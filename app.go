package stdapp

import (
	"io/fs"
	"os"

	"github.com/ddollar/logger"
	"github.com/ddollar/stdcli"
	"github.com/pkg/errors"
)

var version = "dev"

type App struct {
	compose    bool
	database   string
	logger     *logger.Logger
	migrations fs.FS
	name       string
	resolver   ResolverFunc
	schema     string
}

func (a *App) Run(args []string) int {
	// fmt.Printf("args: %+v\n", args)

	c := stdcli.New(a.name, version)

	c.Command("api", "run the api server", a.cliApi, stdcli.CommandOptions{
		Flags: []stdcli.Flag{
			stdcli.BoolFlag("development", "d", "run in development mode"),
			stdcli.IntFlag("port", "p", "port to listen on"),
			stdcli.StringFlag("watch", "w", "comma separated list of file extensions to watch in development mode"),
		},
	})

	c.Command("cmd", "run a command", a.cliCmd, stdcli.CommandOptions{
		Flags: []stdcli.Flag{
			stdcli.BoolFlag("development", "d", "run in development mode"),
			stdcli.StringFlag("watch", "w", "comma separated list of file extensions to watch in development mode"),
		},
		Validate: stdcli.ArgsMin(1),
	})

	c.Command("migrate", "run migrations", a.cliMigrate, stdcli.CommandOptions{})

	c.Command("migration", "create a migration", a.cliMigration, stdcli.CommandOptions{
		Usage:    "<name>",
		Validate: stdcli.Args(1),
	})

	c.Command("pg console", "run database console", a.cliPgConsole, stdcli.CommandOptions{})

	c.Command("pg import", "import contents", a.cliPgImport, stdcli.CommandOptions{})

	c.Command("pg export", "export contents", a.cliPgExport, stdcli.CommandOptions{})

	c.Command("pg reset", "reset databaser", a.cliPgReset, stdcli.CommandOptions{})

	c.Command("web", "start web server", a.cliWeb, stdcli.CommandOptions{
		Flags: []stdcli.Flag{
			stdcli.BoolFlag("development", "d", "run in development mode (vite)"),
		},
	})

	return c.Execute(args)
}

func (a *App) run(container, command string, args ...string) error {
	r := RunnerLocal

	if a.compose {
		tty, err := isTTY(os.Stdin)
		if err != nil {
			return errors.WithStack(err)
		}

		r = RunnerCompose(container, tty)
	}

	cmd := r(command, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func isTTY(f *os.File) (bool, error) {
	stat, err := f.Stat()
	if err != nil {
		return false, errors.WithStack(err)
	}

	return (stat.Mode() & os.ModeCharDevice) != 0, nil
}
