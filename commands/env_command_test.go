package commands

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/cli-mgr/store"
	"github.com/patrickhuber/cli-mgr/store/memory"

	"github.com/patrickhuber/cli-mgr/ui"
	"github.com/spf13/afero"
)

func TestEnvCommand(t *testing.T) {
	t.Run("CanRunCommand", func(t *testing.T) {
		r := require.New(t)

		// create filesystem
		fileSystem := afero.NewMemMapFs()

		// create config file
		configFileContent := `
---
config-sources:
processes:
- name: echo
  environments:
  - name: lab
    process: echo
    env:
      CLI_MGR_TEST: value`
		configFileContent = strings.Replace(configFileContent, "\t", "  ", -1)
		afero.WriteFile(fileSystem, "/config", []byte(configFileContent), 0644)

		// create store manager
		manager := store.NewManager()
		manager.Register(memory.NewMemoryStoreProvider())

		// create console
		console := ui.NewMemoryConsole()

		// create and run command
		cmd := NewEnvCommand(manager, fileSystem, "linux", console)
		runCommandParams := NewRunCommandParams("/config", "echo", "lab")
		err := cmd.ExecuteCommand(runCommandParams)
		r.Nil(err)

		// verify output
		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)
		r.Equal("export CLI_MGR_TEST=value\n", b.String())
	})
}
