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
	resolver   any
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

	// c.Command("reload", "reload a command on changes", a.cliReload, stdcli.CommandOptions{
	// 	Flags: []stdcli.Flag{
	// 		stdcli.StringFlag("extensions", "e", "comma separated list of file extensions to watch"),
	// 	},
	// })

	c.Command("web", "start web server", a.cliWeb, stdcli.CommandOptions{
		Flags: []stdcli.Flag{
			stdcli.BoolFlag("development", "d", "run in development mode"),
		},
	})

	return c.Execute(args)
}

// func (a *App) watchAndExit(extensions []string) error {
// 	a.logger.At("watch").Logf("extensions=%q", strings.Join(extensions, ","))

// 	for {
// 		cmd := exec.Command("go", append([]string{"run", "."}, ctx.Args...)...)

// 		cmd.SysProcAttr = &syscall.SysProcAttr{
// 			Setpgid: true,
// 		}

// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr

// 		if err := cmd.Start(); err != nil {
// 			return errors.WithStack(err)
// 		}

// 		if err := a.watchChanges(extensions); err != nil {
// 			return errors.WithStack(err)
// 		}

// 		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
// 			return errors.WithStack(err)
// 		}

// 		if _, err := cmd.Process.Wait(); err != nil {
// 			return errors.WithStack(err)
// 		}
// 	}
// }

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
