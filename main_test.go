package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/settings"

	"github.com/patrickhuber/wrangle/packages"
	"github.com/patrickhuber/wrangle/tasks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/processes"
	"github.com/patrickhuber/wrangle/store"
	file "github.com/patrickhuber/wrangle/store/file"
	"github.com/patrickhuber/wrangle/ui"
)

var _ = Describe("Main", func() {

	AfterEach(func() {
		os.Unsetenv(global.ConfigFileKey)
		os.Unsetenv(global.PackagePathKey)
	})

	Describe("Print Env", func() {
		It("can cascade configuration", func() {
			runPrintTest("env", "export WRANGLE_TEST=value\n")
		})
	})
	Describe("Print", func() {
		It("can cascade configurations", func() {
			runPrintTest("process", "export WRANGLE_TEST=value\necho\n")
		})
	})

})

func printDirectory(fileSystem filesystem.FileSystem, path string) {
	fmt.Println(path)
	files, err := fileSystem.ReadDir(path)
	Expect(err).To(BeNil())
	for _, f := range files {
		fmt.Println(f.Name())
	}
}

func runPrintTest(command string, expected string) {
	// create dependencies
	platform := "linux"

	fileSystem := filesystem.NewOs()
	storeManager := store.NewManager()
	storeManager.Register(file.NewFileStoreProvider(fileSystem, nil))
	processFactory := processes.NewOsFactory() // change to fake process factory?

	paths := &settings.Paths{
		Root:     "/opt/wrangle",
		Bin:      "/opt/wrangle/bin",
		Packages: "/opt/wrangle/packages",
	}
	feedService := feed.NewFsService(fileSystem, paths.Packages)
	interfaceReader := packages.NewYamlInterfaceReader()
	contextProvider := packages.NewFsContextProvider(fileSystem, paths)

	console := ui.NewMemoryConsole()
	taskProviders := tasks.NewProviderRegistry()
	taskProviders.Register(tasks.NewDownloadProvider(fileSystem, console))
	taskProviders.Register(tasks.NewExtractProvider(fileSystem, console))
	taskProviders.Register(tasks.NewLinkProvider(fileSystem, console))
	taskProviders.Register(tasks.NewMoveProvider(fileSystem, console))

	packageService := packages.NewService(feedService, interfaceReader, contextProvider, taskProviders)

	// create config file
	configFileContent := `
---
stores:
- name: store1
  type: file
  params:
    path: /store1
- name: store2
  type: file 
  params:
    path: /store2
processes:
- name: echo
  path: echo
  stores: 
  - store1
  - store2
  env:
    WRANGLE_TEST: ((key))`

	// create files
	err := fileSystem.Write("/config", []byte(configFileContent), 0644)
	Expect(err).To(BeNil())

	err = fileSystem.Write("/store1", []byte("key: ((key1))"), 0644)
	Expect(err).To(BeNil())

	err = fileSystem.Write("/store2", []byte("key1: value"), 0644)
	Expect(err).To(BeNil())

	// create cli
	app, err := createApplication(
		"/wd",
		storeManager,
		fileSystem,
		processFactory,
		console,
		platform,
		collections.NewDictionary(),
		packageService)
	Expect(err).To(BeNil())
	Expect(app).ToNot(BeNil())

	// run command
	args := []string{
		"wrangle",
		"-c", "/config",
		"print",
		command,
		"-f", "posix",
		"echo"}
	err = app.Run(args)
	Expect(err).To(BeNil())

	// get the output, validate the chaining works
	buffer, ok := console.Out().(*bytes.Buffer)
	Expect(ok).To(BeTrue())
	Expect(buffer).ToNot(BeNil())
	Expect(buffer.String()).To(Equal(expected))
}
