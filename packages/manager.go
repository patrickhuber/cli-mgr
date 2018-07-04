package packages

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/afero"
)

type manager struct {
	fileSystem afero.Fs
}

// Manager defines a manager interface
type Manager interface {
	Download(p Package) error
	Extract(p Package) error
}

// NewManager creates a new package manager
func NewManager(fileSystem afero.Fs) Manager {
	return &manager{fileSystem: fileSystem}
}

func (m *manager) Download(p Package) error {

	if p.Download() == nil {
		return errors.New("package Download() is required")
	}

	// create the file
	file, err := m.fileSystem.Create(p.Download().OutPath())
	if err != nil {
		return err
	}
	defer file.Close()

	// get the file data
	resp, err := http.Get(p.Download().URL())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (m *manager) Extract(p Package) error {

	if p.Download() == nil {
		return errors.New("package Download() is required")
	}
	if p.Extract() == nil {
		return errors.New("package Extract() is required")
	}

	// open the file for reading
	path := p.Download().OutPath()
	file, err := m.fileSystem.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	var reader io.Reader = file

	// based on extension process the file differently
	extension := filepath.Ext(p.Download().OutFile())

	// file is gzipped
	if extension == ".tgz" || extension == ".gz" {
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return err
		}

		if strings.HasSuffix(p.Download().OutFile(), ".tar.gz") {
			extension = ".tar"
		}
	}

	//  the file is a tar archive
	if extension == ".tgz" || extension == ".tar" {
		err = m.extractTar(reader, p.Extract())
		if err != nil {
			return err
		}
		return nil
	}

	// the file is a zip archive
	if extension == ".zip" {
		err = m.extractZip(file, p.Extract())
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unrecoginzed file extension '%s'", extension)
}

func (m *manager) extractTar(reader io.Reader, extract Extract) error {
	// https://gist.github.com/indraniel/1a91458984179ab4cf80
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		name := header.Name
		match, err := isMatch(header.Name, extract.Filter())
		if err != nil {
			return err
		}
		if !match {
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			// create the destination file
			targetFile := extract.OutPath()
			destination, err := m.fileSystem.Create(targetFile)
			if err != nil {
				return err
			}
			defer destination.Close()

			// copy the data to the destination file
			_, err = io.Copy(destination, tarReader)
			if err != nil {
				return err
			}
			m.fileSystem.Chmod(targetFile, 0755)
		default:
			return fmt.Errorf("unable to determine type : '%c' for file '%s' in package", header.Typeflag, name)
		}
	}
	return nil
}

func isMatch(name string, filter string) (bool, error) {
	normalizedName := strings.Replace(name, "\\", "/", -1)

	if normalizedName == filter {
		return true, nil
	}
	return regexp.MatchString(filter, normalizedName)
}

func (m *manager) extractZip(file afero.File, extract Extract) error {
	// http://golang-examples.tumblr.com/post/104726613899/extract-an-uploaded-zip-file

	// get file stat to get file size
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	// create reader
	reader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return err
	}

	for _, zipFile := range reader.File {

		targetFile := extract.OutPath()

		// open destination
		destination, err := m.fileSystem.Create(targetFile)
		if err != nil {
			return err
		}
		defer destination.Close()

		// open source
		source, err := zipFile.Open()
		if err != nil {
			return err
		}
		defer source.Close()

		// copy to the destination
		_, err = io.Copy(destination, source)
		if err != nil {
			return err
		}
		m.fileSystem.Chmod(targetFile, 0755)
	}

	return nil
}
