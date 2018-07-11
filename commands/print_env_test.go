package commands

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/store/file"

	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"
)

func TestPrintEnvCommand(t *testing.T) {
	t.Run("CanRunCommand", func(t *testing.T) {
		r := require.New(t)

		// create filesystem
		fileSystem := afero.NewMemMapFs()

		// create config file
		configFileContent := `
---
config-sources:
environments:
- name: lab
  processes:
  - name: echo
    path: echo
    env:
      WRANGLE_TEST: value`
		configFileContent = strings.Replace(configFileContent, "\t", "  ", -1)
		afero.WriteFile(fileSystem, "/config", []byte(configFileContent), 0644)

		// create store manager
		manager := store.NewManager()
		manager.Register(file.NewFileStoreProvider(fileSystem))

		// create console
		console := ui.NewMemoryConsole()

		// load the config
		loader := config.NewLoader(fileSystem)
		cfg, err := loader.Load("/config")
		r.Nil(err)

		// create and run command
		cmd := NewPrintEnv(manager, fileSystem, "linux", console)
		runCommandParams := NewProcessParams(cfg, "lab", "echo")
		err = cmd.Execute(runCommandParams)
		r.Nil(err)

		// verify output
		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)
		r.Equal("export WRANGLE_TEST=value\n", b.String())
	})

	t.Run("CanRunReplaceVariable", func(t *testing.T) {
		r := require.New(t)

		// create filesystem
		fileSystem := afero.NewMemMapFs()

		// create config file
		configFileContent := `
---
config-sources:
- name: store1
  type: file
  params: 
    path: /store1
environments:
- name: lab
  processes:
  - name: echo
    path: echo
    configurations: [ store1 ]
    env:
      WRANGLE_TEST: ((/key))`
		configFileContent = strings.Replace(configFileContent, "\t", "  ", -1)
		afero.WriteFile(fileSystem, "/config", []byte(configFileContent), 0644)
		afero.WriteFile(fileSystem, "/store1", []byte("key: value"), 0644)

		// create store manager
		manager := store.NewManager()
		manager.Register(file.NewFileStoreProvider(fileSystem))

		// create console
		console := ui.NewMemoryConsole()

		// load the config
		loader := config.NewLoader(fileSystem)
		cfg, err := loader.Load("/config")
		r.Nil(err)

		// create and run command
		cmd := NewPrintEnv(manager, fileSystem, "linux", console)
		runCommandParams := NewProcessParams(cfg, "lab", "echo")
		err = cmd.Execute(runCommandParams)
		r.Nil(err)

		// verify output
		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)
		r.Equal("export WRANGLE_TEST=value\n", b.String())
	})

	t.Run("CanChainReplaceVariables", func(t *testing.T) {

	})
}