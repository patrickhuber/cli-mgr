package commands_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/urfave/cli"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/commands"
	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/packages"
	"github.com/patrickhuber/wrangle/settings"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/ui"
)

var _ = Describe("Install", func() {
	It("can run with environment variables", func() {
		// rewrite this test to use new package management features
		console := ui.NewMemoryConsole()
		variables := collections.NewDictionary()
		fs := filesystem.NewMemory()

		paths := &settings.Paths{
			Root:     "/opt/wrangle",
			Bin:      "/opt/wrangle/bin",
			Packages: "/opt/wrangle/packages",
		}
		feedService := feed.NewFsService(fs, paths.Packages)

		contextProvider := packages.NewFsContextProvider(fs, paths)

		taskProviders := tasks.NewProviderRegistry()
		taskProviders.Register(tasks.NewDownloadProvider(fs, console))
		taskProviders.Register(tasks.NewExtractProvider(fs, console))
		taskProviders.Register(tasks.NewLinkProvider(fs, console))
		taskProviders.Register(tasks.NewMoveProvider(fs, console))

		interfaceReader := packages.NewYamlInterfaceReader()

		service := packages.NewService(feedService, interfaceReader, contextProvider, taskProviders)

		variables.Set(global.PackagePathKey, paths.Packages)
		os.Setenv(global.PackagePathKey, paths.Packages)

		// setup the test server
		message := "this is a message"

		// start the local http server
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(message))
		}))

		defer server.Close()

		content := `
name: test
version: 1.0.0
targets:
- platform: linux
  tasks:
  - download:      
      url: %s
      out: test.html
`
		content = fmt.Sprintf(content, server.URL)

		err := fs.Mkdir(paths.Packages+"/test/1.0.0", 0666)
		Expect(err).To(BeNil())

		err = fs.Write(paths.Packages+"/test/1.0.0/test.1.0.0.yml", []byte(content), 0666)
		Expect(err).To(BeNil())

		app := cli.NewApp()
		app.Flags = []cli.Flag{
			cli.StringFlag{
				Name:   "config, c",
				Usage:  "Load configuration from `FILE`",
				EnvVar: global.ConfigFileKey,
				Value:  "/config",
			},
		}
		app.Commands = []cli.Command{
			*commands.CreateInstallCommand(service, "linux"),
		}

		err = app.Run([]string{
			"wrangle",
			"install",
			"test",
			"-v", "1.0.0",
			"-r", "/wrangle",
		})
		Expect(err).To(BeNil())

		ok, err := fs.Exists(paths.Packages + "/test/1.0.0/test.html")
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

})
