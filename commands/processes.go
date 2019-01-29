package commands

import (
	"github.com/patrickhuber/wrangle/services"
	"github.com/urfave/cli"
)

// CreateListProcessesCommand creates cli command for listing processes from the cli context
func CreateListProcessesCommand(
	app *cli.App,
	processesService services.ProcessesService) *cli.Command {

	command := &cli.Command{
		Name:  "processes",
		Usage: "prints the list of processes for the given environment in the config file",
		Action: func(context *cli.Context) error {
			configFile := context.GlobalString("config")
			return processesService.List(configFile)
		},
	}
	
	setCommandCustomHelpTemplateWithGlobalOptions(app, command)	
	return command
}
