package stdcli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"go.ddollar.dev/errors"
)

type Engine struct {
	Commands []Command
	Executor Executor
	Flags    []Flag
	Name     string
	Reader   *Reader
	Settings string
	Version  string
	Writer   *Writer
}

func (e *Engine) Command(command, description string, fn HandlerFunc, opts CommandOptions) {
	e.Commands = append(e.Commands, Command{
		Command:     strings.Split(command, " "),
		Description: description,
		Handler:     fn,
		Flags:       opts.Flags,
		Invisible:   opts.Invisible,
		Usage:       opts.Usage,
		Validate:    opts.Validate,
		engine:      e,
	})
}

func (e *Engine) Execute(args []string) int {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	return e.ExecuteContext(ctx, args)
}

func (e *Engine) ExecuteContext(ctx context.Context, args []string) int {
	if len(args) > 0 && (args[0] == "-v" || args[0] == "--version") {
		fmt.Println(e.Version) // nolint:forbidigo
		return 0
	}

	var m *Command
	var cargs []string

	for _, c := range e.Commands {
		d := c
		if a, ok := d.Match(args); ok {
			if m == nil || len(m.Command) < len(c.Command) {
				m = &d
				cargs = a
			}
		}
	}

	if m == nil {
		m = &(e.Commands[0])
	}

	err := m.ExecuteContext(ctx, cargs)
	switch t := errors.Cause(err).(type) {
	case nil:
		return 0
	case ExitCoder:
		return t.ExitCode()
	default:
		e.Writer.Error(err) //nolint:errcheck
		return 1
	}
}
