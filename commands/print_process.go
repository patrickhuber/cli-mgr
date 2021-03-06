package commands

import (
	"strings"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/processes"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// CreatePrintProcessCommand Creates a Print Command from the cli context
func CreatePrintProcessCommand(
	app *cli.App,
	printService processes.PrintService,
	fs filesystem.FileSystem) cli.Command {

	command := cli.Command{
		Name:      "process",
		Aliases:   []string{"ps"},
		Usage:     "prints the process as it would be executed",
		ArgsUsage: "<process name> [arguments]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "format, f",
				Usage: "Print for with the given format (bash|powershell)",
			},
		},
		Action: func(context *cli.Context) error {
			processName := context.Args().First()
			if strings.TrimSpace(processName) == "" {
				return errors.New("process name argument is required")
			}

			format := context.String("format")

			configFile := context.GlobalString("config")

			configProvider := config.NewFsProvider(fs, configFile)

			cfg, err := configProvider.Get()
			if err != nil {
				return err
			}

			params := &processes.PrintParams{
				Config:      cfg,
				ProcessName: processName,
				Format:      format,
				Include: processes.PrintParamsInclude{
					ProcessAndArgs: true,
				},
			}

			return printService.Print(params)
		},
	}

	setCommandCustomHelpTemplateWithGlobalOptions(app, &command)
	return command
}
