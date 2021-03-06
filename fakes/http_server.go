package fakes

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/patrickhuber/wrangle/archiver"
	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/filesystem"
)

// NewHTTPServer creates a new http server
func NewHTTPServer() *httptest.Server {
	return NewHTTPServerWithArchive([]TestFile{{"/data", "this is data"}})
}

// NewHTTPServerWithArchive creates a new http server with
func NewHTTPServerWithArchive(testFiles []TestFile) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		fs := filesystem.NewMemory()

		var filePaths = make([]string, 0)

		// write out the data files
		for _, testFile := range testFiles {
			filePaths = append(filePaths, testFile.Path)
			err := fs.Write(testFile.Path, []byte(testFile.Data), 0666)
			if err != nil {
				rw.WriteHeader(400)
				rw.Write([]byte("error creating executable"))
				return
			}
		}

		// create the archiver for the file extension if not found assume binary
		var a archiver.Archiver
		if strings.HasSuffix(path, ".tgz") || strings.HasSuffix(path, ".tar.gz") {
			a = archiver.NewTargz(fs)
		} else if strings.HasSuffix(path, ".tar") {
			a = archiver.NewTar(fs)
		} else if strings.HasSuffix(path, ".zip") {
			a = archiver.NewZip(fs)
		} else {
			ok, err := fs.Exists(path)
			if err != nil {
				rw.WriteHeader(400)
				rw.Write([]byte(fmt.Sprintf("error finding file %s: %s", path, err.Error())))
				return
			}
			if !ok {
				rw.WriteHeader(404)
				rw.Write([]byte(fmt.Sprintf("file %s not found", path)))
				return
			}
		}

		_, fileName := filepath.Split(path)
		loadPath := path
		if a != nil {
			err := a.Archive(fileName, filePaths)
			if err != nil {
				rw.WriteHeader(400)
				rw.Write([]byte(fmt.Sprintf("error creating archive: %s", err.Error())))
				return
			}
			loadPath = fileName
		}

		buf, err := fs.Read(loadPath)
		if err != nil {
			rw.WriteHeader(400)
			rw.Write([]byte(fmt.Sprintf("error reading file %s: %s", fileName, err.Error())))
			return
		}
		rw.WriteHeader(200)
		rw.Write(buf)
		return
	}))
}
