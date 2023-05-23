package stdapp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ddollar/coalesce"
	"github.com/ddollar/migrate"
	"github.com/ddollar/stdcli"
	"github.com/pkg/errors"
)

func (a *App) cliApi(ctx *stdcli.Context) error {
	if ctx.Bool("development") {
		return a.watchAndReload(parseExtensions(ctx.String("watch")), "api", "--port", fmt.Sprint(ctx.Int("port")))
	}

	g, err := a.graphQL()
	if err != nil {
		return errors.WithStack(err)
	}

	port := coalesce.Any(ctx.Int("port"), 8000)

	if err := g.server.Listen("https", fmt.Sprintf(":%d", port)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *App) cliCmd(ctx *stdcli.Context) error {
	if ctx.Bool("development") {
		return a.watchAndReload(parseExtensions(ctx.String("watch")), "cmd", ctx.Args...)
	}

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
