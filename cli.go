package stdapp

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/ddollar/stdcli"
	"github.com/pkg/errors"
)

func (a *App) cliApi(ctx *stdcli.Context) error {
	g, err := a.graphQL()
	if err != nil {
		return err
	}

	if err := g.server.Listen("https", ":8000"); err != nil {
		return errors.WithStack(err)
	}

	return nil
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
			return err
		}

		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
			return errors.WithStack(err)
		}

		if _, err := cmd.Process.Wait(); err != nil {
			return errors.WithStack(err)
		}
	}
}
