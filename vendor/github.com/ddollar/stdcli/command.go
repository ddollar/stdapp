package stdcli

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type Command struct {
	Command     []string
	Description string
	Flags       []Flag
	Invisible   bool
	Handler     HandlerFunc
	Usage       string
	Validate    Validator

	engine *Engine
}

type CommandOptions struct {
	Flags     []Flag
	Invisible bool
	Usage     string
	Validate  Validator
}

type HandlerFunc func(Context) error

func (c *Command) ExecuteContext(ctx context.Context, args []string) error {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.Usage = func() { helpCommand(c.engine, c) }

	flags := []*Flag{}

	for _, f := range c.Flags {
		g := f
		flags = append(flags, &g)
		flag := fs.VarPF(&g, f.Name, f.Short, f.Description)
		if f.Type() == "bool" {
			flag.NoOptDefVal = "true"
		}
	}

	if err := fs.Parse(args); err != nil {
		if strings.HasPrefix(err.Error(), "unknown shorthand flag") {
			parts := strings.Split(err.Error(), " ")
			return errors.Errorf("unknown flag: %s", parts[len(parts)-1])
		}
		if err == pflag.ErrHelp {
			return nil
		}
		return errors.WithStack(err)
	}

	cc := &defaultContext{
		Context: ctx,
		args:    fs.Args(),
		flags:   flags,
		engine:  c.engine,
	}

	if c.Validate != nil {
		if err := c.Validate(cc); err != nil {
			return err
		}
	}

	if err := c.Handler(cc); err != nil {
		return err
	}

	return nil
}

func (c *Command) FullCommand() string {
	return filepath.Base(os.Args[0]) + " " + strings.Join(c.Command, " ")
}

func (c *Command) Match(args []string) ([]string, bool) {
	if len(args) < len(c.Command) {
		return args, false
	}

	for i := range c.Command {
		if args[i] != c.Command[i] {
			return args, false
		}
	}

	return args[len(c.Command):], true
}

type CommandDefinition struct {
	Command     string
	Description string
	Handler     HandlerFunc
	Options     CommandOptions
}

type CommandDefinitions []CommandDefinition

func (cs CommandDefinitions) Apply(e *Engine) {
	for _, c := range cs {
		e.Command(c.Command, c.Description, c.Handler, c.Options)
	}
}

func (cs *CommandDefinitions) Register(command, description string, fn HandlerFunc, opts CommandOptions) {
	*cs = append(*cs, CommandDefinition{
		Command:     command,
		Description: description,
		Handler:     fn,
		Options:     opts,
	})
}
