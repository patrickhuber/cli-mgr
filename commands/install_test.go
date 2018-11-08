package commands_test

import (
	"fmt"
	"net/http/httptest"
	"path"

	"github.com/patrickhuber/wrangle/fakes"
	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/ui"

	"github.com/patrickhuber/wrangle/tasks"

	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/commands"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/packages"
)

var _ = Describe("Install", func() {
	var (
		platform     string
		packagesPath string
		fileSystem   filesystem.FsWrapper
		manager      packages.Manager
		loader       config.Loader
	)
	Describe("NewInstall", func() {
		It("returns install command", func() {
			command, err := commands.NewInstall(platform, packagesPath, fileSystem, manager, loader)
			Expect(err).To(BeNil())
			Expect(command).ToNot(BeNil())
		})
	})
	Describe("Execute", func() {
		const packagesRootPosix = "/opt/wrangle/pacakges"
		const packagesRootWindows = "c:" + packagesRootPosix
		var (
			platform         string
			downloadFileName string
			archive          string
			destination      string
			server           *httptest.Server
		)
		BeforeSuite(func() {
			server = fakes.NewHTTPServerWithArchive(
				[]fakes.TestFile{
					{Path: "/test", Data: "this is data"},
					{Path: "/test.exe", Data: "this is data"},
				})
		})
		AfterSuite(func() {
			server.Close()
		})
		AfterEach(func() {
			url := server.URL
			packageVersion := "1.0.0"
			packageName := "test"
			packagesRoot := packagesRootPosix
			if platform == "windows" {
				packagesRoot = packagesRootWindows
			}

			// create command dependencies
			console := ui.NewMemoryConsole()
			fs := filesystem.NewMemMapFs()
			loader := config.NewLoader(fs)

			taskProviders := tasks.NewProviderRegistry()
			taskProviders.Register(tasks.NewExtractProvider(fs, console))
			taskProviders.Register(tasks.NewDownloadProvider(fs, console))
			taskProviders.Register(tasks.NewMoveProvider(fs, console))
			taskProviders.Register(tasks.NewLinkProvider(fs, console))

			manager := packages.NewManager(fs, taskProviders)

			out := filepath.Join("/", downloadFileName)
			url = path.Join(url, downloadFileName)

			// create the package manifest
			packageManifest, err := createPackageManifest(packageName, packageVersion, platform, url, out, archive, destination)
			Expect(err).To(BeNil())

			packagePath := filepath.Join(packagesRoot, packageName, packageVersion)
			packageManifestPath := filepath.Join(packagePath, fmt.Sprintf("%s.%s.yml", packageName, packageVersion))
			err = afero.WriteFile(fs, packageManifestPath, []byte(packageManifest), 0600)
			Expect(err).To(BeNil())

			// create the command and execute it
			command, err := commands.NewInstall(platform, packagesRoot, fs, manager, loader)
			Expect(err).To(BeNil())

			err = command.Execute(packageName, packageVersion)
			Expect(err).To(BeNil())
		})
		When("Windows", func() {
			BeforeEach(func() {
				platform = "windows"
			})
			When("Tar", func() {
				It("installs", func() {
					downloadFileName = "test.tar"
				})
			})
			When("Tgz", func() {
				It("installs", func() {
					downloadFileName = "test.tgz"
				})
			})
			When("Zip", func() {
				It("installs", func() {
					downloadFileName = "test.zip"
				})
			})
			When("Binary", func() {
				It("installs", func() {
					downloadFileName = "test.exe"
				})
			})
		})
		When("Linux", func() {
			BeforeEach(func() {
				platform = "linux"
			})
		})
		When("Darwin", func() {
			BeforeEach(func() {
				platform = "darwin"
			})
		})
	})
})

func createPackageManifest(
	name string,
	version string,
	platform string,
	url string,
	outFile string,
	archive string,
	destination string) (string, error) {
	taskList := []config.Task{
		config.Task{
			Name: "download",
			Type: "download",
			Params: map[string]interface{}{
				"url": url,
				"out": outFile,
			},
		},
	}
	if archive != "" && destination != "" {
		extractTask := config.Task{
			Name: "extract",
			Type: "extract",
			Params: map[string]interface{}{
				"archive":     archive,
				"destination": destination,
			},
		}
		taskList = append(taskList, extractTask)
	}
	packageConfig := &config.Package{
		Name:    name,
		Version: version,
		Platforms: []config.Platform{
			config.Platform{
				Name:  platform,
				Tasks: taskList,
			},
		},
	}

	return config.SerializePackage(packageConfig)
}
