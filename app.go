package stdapp

import (
	"fmt"
	"io/fs"
	"os"
	"time"

	"go.ddollar.dev/coalesce"
	"go.ddollar.dev/errors"
	"go.ddollar.dev/logger"
	"go.ddollar.dev/stdcli"
)

var version = "dev"

var (
	flagDevelopment = stdcli.BoolFlag("development", "d", "run in development mode")
	flagWatch       = stdcli.StringFlag("watch", "w", "comma separated list of file extensions to watch in development mode")
)

type App struct {
	logger *logger.Logger
	opts   Options
}

type Options struct {
	Compose      bool
	Database     string
	Domains      []string
	Middleware   []Middleware
	Migrations   fs.FS
	Name         string
	Prefix       string
	Resolver     ResolverFunc
	Router       RouterFunc
	Web          fs.FS
	WriteTimeout time.Duration
}

func (a *App) Run(args []string) int {
	c := stdcli.New(a.opts.Name, version)

	c.Command("api", "run the api server", a.cliApi, stdcli.CommandOptions{
		Flags: []stdcli.Flag{
			flagDevelopment,
			flagWatch,
			stdcli.IntFlag("port", "p", "port to listen on"),
		},
	})

	c.Command("cmd", "run a command", a.cliCmd, stdcli.CommandOptions{
		Flags: []stdcli.Flag{
			flagDevelopment,
			flagWatch,
		},
		Validate: stdcli.ArgsMin(1),
	})

	c.Command("cron", "run cron daemon", a.cliCron, stdcli.CommandOptions{
		Flags: []stdcli.Flag{
			flagDevelopment,
			flagWatch,
		},
	})

	c.Command("deployment", "run a command on the deploy target", a.cliDeployment, stdcli.CommandOptions{})

	c.Command("init", "initialize a new project", a.cliInit, stdcli.CommandOptions{
		Validate: stdcli.Args(1),
	})

	c.Command("migrate", "run migrations", a.cliMigrate, stdcli.CommandOptions{
		Flags: []stdcli.Flag{
			stdcli.BoolFlag("dry", "", "dry run"),
		},
	})

	c.Command("migration", "create a migration", a.cliMigration, stdcli.CommandOptions{
		Flags: []stdcli.Flag{
			stdcli.StringFlag("dir", "d", "dir in which to create migration"),
		},
		Usage:    "<name>",
		Validate: stdcli.Args(1),
	})

	c.Command("pg console", "run database console", a.cliPgConsole, stdcli.CommandOptions{
		Flags: []stdcli.Flag{
			stdcli.StringFlag("schema", "s", "database schema to run migrations in"),
		},
	})

	c.Command("pg import", "import contents", a.cliPgImport, stdcli.CommandOptions{})

	c.Command("pg export", "export contents", a.cliPgExport, stdcli.CommandOptions{})

	c.Command("pg reset", "reset databaser", a.cliPgReset, stdcli.CommandOptions{})

	c.Command("sleep", "sleep forever", a.cliSleep, stdcli.CommandOptions{})

	c.Command("web", "start web server", a.cliWeb, stdcli.CommandOptions{
		Flags: []stdcli.Flag{
			flagDevelopment,
			stdcli.IntFlag("port", "p", "port to listen on"),
		},
	})

	return c.Execute(args)
}

func (a *App) domains() []string {
	return coalesce.Any(a.opts.Domains, []string{"public"})
}

func (a *App) run(container, command string, args ...string) error {
	return a.runEnv(container, nil, command, args...)
}

func (a *App) runEnv(container string, env map[string]string, command string, args ...string) error {
	r := RunnerLocal

	if a.opts.Compose {
		tty, err := isTTY(os.Stdin)
		if err != nil {
			return errors.Wrap(err)
		}

		r = RunnerCompose(container, tty, env)
	}

	cmd := r(command, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = os.Environ()

	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func isTTY(f *os.File) (bool, error) {
	stat, err := f.Stat()
	if err != nil {
		return false, errors.Wrap(err)
	}

	return (stat.Mode() & os.ModeCharDevice) != 0, nil
}
