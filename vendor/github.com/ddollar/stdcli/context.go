package stdcli

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/term"
)

type Context interface {
	context.Context
	io.ReadWriter

	Arg(i int) string
	Args() []string
	Execute(cmd string, args ...string) ([]byte, error)
	Flags() Flags
	InfoWriter() InfoWriter
	IsTerminal() bool
	ReadSecret() (string, error)
	Run(cmd string, args ...string) error
	Terminal(cmd string, args ...string) error
	Version() string
	Writef(format string, args ...any)
}

type defaultContext struct {
	context.Context

	args   []string
	engine *Engine
	flags  Flags
}

var _ Context = &defaultContext{}

func (c *defaultContext) Arg(i int) string {
	if i < len(c.args) {
		return c.args[i]
	}

	return ""
}

func (c *defaultContext) Args() []string {
	return []string(c.args)
}

func (c *defaultContext) Engine() *Engine {
	return c.engine
}

func (c *defaultContext) Execute(cmd string, args ...string) ([]byte, error) {
	if c.engine.Executor == nil {
		return nil, errors.Errorf("no executor")
	}

	data, err := c.engine.Executor.Execute(c, cmd, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return data, nil
}

func (c *defaultContext) Flags() Flags {
	return c.flags
}

func (c *defaultContext) InfoWriter() InfoWriter {
	return &infoWriter{ctx: c}
}

func (c *defaultContext) IsTerminal() bool {
	return c.engine.Reader.IsTerminal()
}

func (c *defaultContext) Read(data []byte) (int, error) {
	n, err := c.engine.Reader.Read(data)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return n, nil
}

func (c *defaultContext) ReadSecret() (string, error) {
	data, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(data), nil
}

func (c *defaultContext) Run(cmd string, args ...string) error {
	if c.engine.Executor == nil {
		return errors.Errorf("no executor")
	}

	if err := c.engine.Executor.Run(c, c, cmd, args...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *defaultContext) Terminal(cmd string, args ...string) error {
	if c.engine.Executor == nil {
		return errors.Errorf("no executor")
	}

	if err := c.engine.Executor.Terminal(c, cmd, args...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *defaultContext) Version() string {
	return c.engine.Version
}

func (c *defaultContext) Write(data []byte) (int, error) {
	return c.engine.Writer.Write(data)
}

func (c *defaultContext) Writef(format string, args ...any) {
	c.engine.Writer.Write([]byte(fmt.Sprintf(format, args...))) //nolint:errcheck
}

func (c *defaultContext) TableWriter(columns ...string) TableWriter {
	return &tableWriter{ctx: c, columns: columns}
}
