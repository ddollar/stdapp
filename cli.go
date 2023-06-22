package stdapp

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ddollar/coalesce"
	"github.com/ddollar/migrate"
	"github.com/ddollar/stdapi"
	"github.com/ddollar/stdcli"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
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

	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *App) cliCron(ctx *stdcli.Context) error {
	if ctx.Bool("development") {
		return a.watchAndReload(parseExtensions(ctx.String("watch")), "cron")
	}

	cc, err := NewCron(ctx)
	if err != nil {
		return err
	}

	if err := cc.Run(); err != nil {
		return err
	}

	return nil
}

func (a *App) cliDeployment(ctx *stdcli.Context) error {
	data, err := exec.Command("git", "remote", "get-url", "deploy").CombinedOutput()
	if err != nil {
		return fmt.Errorf("no deploy remote found")
	}

	u, err := url.Parse(strings.TrimSpace(string(data)))
	if err != nil {
		return err
	}

	cmd := exec.Command("ssh", "-t", u.Host, fmt.Sprintf(`cd %s && bash -l -c "sa %s"`, strings.TrimPrefix(u.Path, "/"), strings.Join(ctx.Args, " ")))

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *App) cliInit(ctx *stdcli.Context) error {
	name := ctx.Arg(0)

	if err := initApp(name); err != nil {
		return err
	}

	return nil
}

func (a *App) cliMigrate(ctx *stdcli.Context) error {
	if a.opts.Compose {
		return a.run("api", "go", "run", ".", "migrate")
	}

	if err := migrate.Run(a.opts.Database, a.opts.Migrations); err != nil {
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
	return a.run("postgres", "psql", a.opts.Database)
}

func (a *App) cliPgExport(ctx *stdcli.Context) error {
	return a.run("postgres", "pg_dump", "--clean", "--no-acl", "--no-owner", a.opts.Database)
}

func (a *App) cliPgImport(ctx *stdcli.Context) error {
	return a.run("postgres", "psql", a.opts.Database)
}

func (a *App) cliPgReset(ctx *stdcli.Context) error {
	return a.run("postgres", "psql", a.opts.Database, "-c", "drop schema public cascade; create schema public;")
}

func (a *App) cliWeb(ctx *stdcli.Context) error {
	if ctx.Bool("development") {
		return a.webDevelopment()
	}

	s, err := a.spa()
	if err != nil {
		return err
	}

	port := coalesce.Any(ctx.Int("port"), 8000)

	if err := s.server.Listen("https", fmt.Sprintf(":%d", port)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *App) webDevelopment() error {
	eg := new(errgroup.Group)

	eg.Go(a.webDevelopmentProxy)
	eg.Go(a.webDevelopmentVite)

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (a *App) webDevelopmentProxy() error {
	s := stdapi.New(a.opts.Name, a.opts.Name)

	u, err := url.Parse("http://localhost:8001")
	if err != nil {
		return err
	}

	rp := httputil.NewSingleHostReverseProxy(u)

	s.Router.PathPrefix(a.opts.Prefix).Handler(a.WithMiddleware(rp))

	if err := s.Listen("https", ":8000"); err != nil {
		return err
	}

	return nil
}

func (a *App) webDevelopmentVite() error {
	if err := os.Chdir("web"); err != nil {
		return errors.WithStack(err)
	}

	cmd := exec.Command("npx", "vite", "--host")

	cmd.Env = append(os.Environ(),
		"PORT=8001",
		fmt.Sprintf("VITE_PREFIX=%s", coalesce.Any(a.opts.Prefix, "/")),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
