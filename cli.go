package stdapp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ddollar/migrate"
	"github.com/ddollar/stdcli"
	"github.com/pkg/errors"
)

func (a *App) cliApi(ctx *stdcli.Context) error {
	g, err := a.graphQL()
	if err != nil {
		return errors.WithStack(err)
	}

	if err := g.server.Listen("https", ":8000"); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *App) cliCmd(ctx *stdcli.Context) error {
	cmd := exec.Command("go", append([]string{"run", fmt.Sprintf("./cmd/%s", ctx.Args[0])}, ctx.Args[1:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (a *App) cliMigrate(ctx *stdcli.Context) error {
	if a.compose {
		return a.run("api", "go", "run", ".", "migrate")
	}

	if err := migrate.Run(a.database, a.migrations); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *App) cliMigration(ctx *stdcli.Context) error {
	name := ctx.Arg(0)

	file := filepath.Join("db", "migrate", fmt.Sprintf("%s_%s.sql", time.Now().Format("20060102150405"), name))

	fd, err := os.Create(file)
	if err != nil {
		return errors.WithStack(err)
	}
	defer fd.Close()

	if err := ctx.Writef("%s\n", file); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *App) cliPgConsole(ctx *stdcli.Context) error {
	return a.run("postgres", "psql", a.database)
}

func (a *App) cliPgExport(ctx *stdcli.Context) error {
	return a.run("postgres", "pg_dump", "--clean", "--no-acl", "--no-owner", a.database)
}

func (a *App) cliPgImport(ctx *stdcli.Context) error {
	return a.run("postgres", "psql", a.database)
}

func (a *App) cliPgReset(ctx *stdcli.Context) error {
	return a.run("postgres", "psql", a.database, "-c", "drop schema public cascade; create schema public;")
}

func (a *App) cliReload(ctx *stdcli.Context) error {
	extensions := strings.Split(ctx.String("extensions"), ",")

	a.logger.At("reload").Logf("extensions=%q", strings.Join(extensions, ","))

	for {
		cmd := exec.Command("go", append([]string{"run", "."}, ctx.Args...)...)

		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			return errors.WithStack(err)
		}

		if err := a.watchChanges(extensions); err != nil {
			return errors.WithStack(err)
		}

		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
			return errors.WithStack(err)
		}

		if _, err := cmd.Process.Wait(); err != nil {
			return errors.WithStack(err)
		}
	}
}

func (a *App) cliWeb(ctx *stdcli.Context) error {
	if ctx.Bool("development") {
		return a.cliWebDevelopment()
	}

	return nil
}

func (a *App) cliWebDevelopment() error {
	if err := os.Chdir("web"); err != nil {
		return errors.WithStack(err)
	}

	cmd := exec.Command("npx", "vite", "--host")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
