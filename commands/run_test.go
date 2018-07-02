package commands

import (
	"strings"
	"testing"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/processes"
	"github.com/patrickhuber/cli-mgr/store"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestRunCommand(t *testing.T) {

	t.Run("CanRunGoVersionProcess", func(t *testing.T) {
		r := require.New(t)

		// write out the config file
		configFileData := `
environments:
- name: lab
  processes:
  - name: go
    path: go
    args:
    - version
`
		configFileData = strings.Replace(configFileData, "\t", "  ", -1)
		fileSystem := afero.NewMemMapFs()
		err := afero.WriteFile(fileSystem, "/config", []byte(configFileData), 0644)
		r.Nil(err)

		// create run command
		configStoreManager := store.NewManager()
		runCommand := NewRun(configStoreManager, fileSystem, processes.NewOsFactory())

		// load the config
		loader := config.NewLoader(fileSystem)
		cfg, err := loader.Load("/config")
		r.Nil(err)

		// run the run command
		err = runCommand.Execute(
			NewRunCommandParams(cfg, "go", "lab"))
		r.Nil(err)
	})
}