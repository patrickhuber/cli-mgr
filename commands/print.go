package commands

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/renderers"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"
)

type print struct {
	fileSystem      afero.Fs
	console         ui.Console
	manager         store.Manager
	rendererFactory renderers.Factory
}

// PrintParams defines parameters for the print command
type PrintParams struct {
	Configuration   *config.Config
	EnvironmentName string
	ProcessName     string
	Shell           string
}

// Print represents an environment command
type Print interface {
	Execute(params *PrintParams) error
}

// NewPrint creates a new environment command
func NewPrint(
	manager store.Manager,
	fileSystem afero.Fs,
	console ui.Console,
	rendererFactory renderers.Factory) Print {
	return &print{
		manager:         manager,
		fileSystem:      fileSystem,
		console:         console,
		rendererFactory: rendererFactory}
}

func (cmd *print) Execute(params *PrintParams) error {

	processName := params.ProcessName
	environmentName := params.EnvironmentName

	if processName == "" {
		return errors.New("process name is required for the run command")
	}

	if environmentName == "" {
		return errors.New("environment name is required for the run command")
	}

	cfg := params.Configuration
	if cfg == nil {
		return errors.New("unable to load configuration")
	}

	processTemplate, err := store.NewProcessTemplate(cfg, cmd.manager)
	if err != nil {
		return err
	}

	process, err := processTemplate.Evaluate(environmentName, processName)
	if err != nil {
		return err
	}

	renderer, err := cmd.rendererFactory.Create(params.Shell)
	if err != nil {
		return err
	}

	fmt.Fprint(
		cmd.console.Out(),
		renderer.RenderProcess(
			process.Path,
			process.Args,
			process.Vars))
	return nil
}
