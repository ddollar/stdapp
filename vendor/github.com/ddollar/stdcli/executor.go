package stdcli

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

type Executor interface {
	Execute(ctx context.Context, cmd string, args ...string) ([]byte, error)
	Run(ctx context.Context, w io.Writer, cmd string, args ...string) error
	Terminal(ctx context.Context, cmd string, args ...string) error
}

type defaultExecutor struct{}

func (e *defaultExecutor) Execute(ctx context.Context, cmd string, args ...string) ([]byte, error) {
	data, err := exec.CommandContext(ctx, cmd, args...).CombinedOutput()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return data, nil
}

func (e *defaultExecutor) Run(ctx context.Context, w io.Writer, cmd string, args ...string) error {
	c := exec.CommandContext(ctx, cmd, args...)

	c.Stdout = w
	c.Stderr = w

	if err := c.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (e *defaultExecutor) Terminal(ctx context.Context, cmd string, args ...string) error {
	c := exec.CommandContext(ctx, cmd, args...)

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
