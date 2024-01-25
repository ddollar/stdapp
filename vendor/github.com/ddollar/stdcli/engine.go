package stdcli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
)

type Engine struct {
	Commands []Command
	Executor Executor
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
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	defer stop(c, cancel)
	go wait(ctx, c, cancel)

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
	switch t := err.(type) {
	case nil:
		return 0
	case ExitCoder:
		return t.Code()
	default:
		e.Writer.Errorf("%s", t) // nolint:errcheck
		return 1
	}
}

func stop(c chan<- os.Signal, cancel func()) {
	signal.Stop(c)
	cancel()
}

func wait(ctx context.Context, c <-chan os.Signal, cancel func()) {
	select {
	case <-c:
		cancel()
	case <-ctx.Done():
	}
}
