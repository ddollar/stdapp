package stdapp

import (
	"context"
	"fmt"
	"io/fs"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ddollar/coalesce"
	"github.com/ddollar/errors"
	"github.com/ddollar/migrate"
	"github.com/ddollar/stdapi"
	"github.com/ddollar/stdcli"
	"golang.org/x/sync/errgroup"
)

func (a *App) cliApi(ctx stdcli.Context) error {
	if ctx.Flags().Bool("development") {
		return a.watchAndReload(parseExtensions(ctx.Flags().String("watch")), "api", "--port", fmt.Sprint(ctx.Flags().Int("port")))
	}

	api, err := a.api()
	if err != nil {
		return errors.Wrap(err)
	}

	port := coalesce.Any(ctx.Flags().Int("port"), 8000)

	if err := api.server.Listen("https", fmt.Sprintf(":%d", port)); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (a *App) cliCmd(ctx stdcli.Context) error {
	args := ctx.Args()

	if ctx.Flags().Bool("development") {
		return a.watchAndReload(parseExtensions(ctx.Flags().String("watch")), "cmd", append([]string{"go", "run", fmt.Sprintf("./cmd/%s", args[0])}, args[1:]...)...)
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (a *App) cliCron(ctx stdcli.Context) error {
	if ctx.Flags().Bool("development") {
		return a.watchAndReload(parseExtensions(ctx.Flags().String("watch")), "cron")
	}

	cc, err := NewCron(ctx)
	if err != nil {
		return errors.Wrap(err)
	}

	if err := cc.Run(); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (a *App) cliDeployment(ctx stdcli.Context) error {
	data, err := exec.Command("git", "remote", "get-url", "deploy").CombinedOutput()
	if err != nil {
		return errors.Errorf("no deploy remote found")
	}

	u, err := url.Parse(strings.TrimSpace(string(data)))
	if err != nil {
		return errors.Wrap(err)
	}

	cmd := exec.Command("ssh", "-t", u.Host, fmt.Sprintf(`cd %s && bash -l -c "sa %s"`, strings.TrimPrefix(u.Path, "/"), strings.Join(ctx.Args(), " ")))

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (a *App) cliInit(ctx stdcli.Context) error {
	name := ctx.Arg(0)

	if err := initApp(name); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (a *App) cliMigrate(ctx stdcli.Context) error {
	args := []string{}
	dry := false

	if ctx.Flags().Bool("dry") {
		args = append(args, "--dry")
		dry = true
	}

	if a.opts.Compose {
		return a.run("api", "go", append([]string{"run", ".", "migrate"}, args...)...)
	}

	u, err := url.Parse(a.opts.Database)
	if err != nil {
		return errors.Wrap(err)
	}

	mopts := migrate.Options{
		DryRun: dry,
	}

	for _, domain := range append([]string{""}, a.domains()...) {
		q := u.Query()
		q.Set("search_path", domain)
		u.RawQuery = q.Encode()

		mfs, err := fs.Sub(a.opts.Migrations, filepath.Join("db", "migrate", domain))
		if err != nil {
			return errors.Wrap(err)
		}

		if err := migrate.Run(context.Background(), u.String(), mfs, mopts); err != nil {
			return errors.Wrap(err)
		}
	}

	return nil
}

func (a *App) cliMigration(ctx stdcli.Context) error {
	name := ctx.Arg(0)

	dir := coalesce.Any(ctx.Flags().String("dir"), filepath.Join("db", "migrate"))
	file := filepath.Join(dir, fmt.Sprintf("%s_%s.sql", time.Now().Format("20060102150405"), name))

	fd, err := os.Create(file)
	if err != nil {
		return errors.Wrap(err)
	}
	defer fd.Close()

	ctx.Writef("%s\n", file)

	return nil
}

func (a *App) cliPgConsole(ctx stdcli.Context) error {
	schema := coalesce.Any(ctx.Flags().String("schema"), "public")

	env := map[string]string{
		"PGOPTIONS": fmt.Sprintf("--search_path=%s", schema),
	}

	return a.runEnv("postgres", env, "psql", a.opts.Database)
}

func (a *App) cliPgExport(ctx stdcli.Context) error {
	return a.run("postgres", "pg_dump", "--clean", "--no-acl", "--no-owner", a.opts.Database)
}

func (a *App) cliPgImport(ctx stdcli.Context) error {
	return a.run("postgres", "psql", a.opts.Database)
}

func (a *App) cliPgReset(ctx stdcli.Context) error {
	return a.run("postgres", "psql", a.opts.Database, "-c", "drop schema public cascade; create schema public;")
}

func (a *App) cliSleep(ctx stdcli.Context) error {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	return nil
}

func (a *App) cliWeb(ctx stdcli.Context) error {
	if ctx.Flags().Bool("development") {
		return a.webDevelopment()
	}

	s, err := a.spa()
	if err != nil {
		return errors.Wrap(err)
	}

	port := coalesce.Any(ctx.Flags().Int("port"), 8000)

	if err := s.server.Listen("https", fmt.Sprintf(":%d", port)); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (a *App) webDevelopment() error {
	eg := new(errgroup.Group)

	eg.Go(a.webDevelopmentProxy)
	eg.Go(a.webDevelopmentVite)

	if err := eg.Wait(); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (a *App) webDevelopmentProxy() error {
	s := stdapi.New(a.opts.Name, a.opts.Name)

	u, err := url.Parse("http://localhost:8001")
	if err != nil {
		return errors.Wrap(err)
	}

	rp := httputil.NewSingleHostReverseProxy(u)

	s.Router.PathPrefix(a.opts.Prefix).Handler(a.WithMiddleware(rp))

	if err := s.Listen("https", ":8000"); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (a *App) webDevelopmentVite() error {
	if err := os.Chdir("web"); err != nil {
		return errors.Wrap(err)
	}

	cmd := exec.Command("yarn", "run", "dev")

	cmd.Env = append(os.Environ(),
		"PORT=8001",
		fmt.Sprintf("VITE_PREFIX=%s", coalesce.Any(a.opts.Prefix, "/")),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err)
	}

	return nil
}
