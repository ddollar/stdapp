package stdapp

import (
	"io/fs"

	"github.com/ddollar/logger"
	"github.com/ddollar/stdcli"
)

var version = "dev"

type App struct {
	database   Database
	logger     *logger.Logger
	migrations fs.FS
	name       string
	resolver   any
	schema     string
}

func (a *App) Run(args []string) int {
	// fmt.Printf("args: %+v\n", args)

	c := stdcli.New(a.name, version)

	c.Command("api", "run the api server", a.cliApi, stdcli.CommandOptions{})

	c.Command("migration", "create a migration", a.cliMigration, stdcli.CommandOptions{
		Usage:    "<name>",
		Validate: stdcli.Args(1),
	})

	c.Command("reload", "reload a command on changes", a.cliReload, stdcli.CommandOptions{
		Flags: []stdcli.Flag{
			stdcli.StringFlag("extensions", "e", "comma separated list of file extensions to watch"),
		},
	})

	return c.Execute(args)
}
