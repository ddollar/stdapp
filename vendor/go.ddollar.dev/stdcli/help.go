package stdcli

import (
	"fmt"
	"sort"
)

func help(e *Engine) HandlerFunc {
	return func(ctx Context) error {
		cs := []Command{}

		for _, cmd := range e.Commands {
			if cmd.Invisible {
				continue
			}

			cs = append(cs, cmd)
		}

		sort.Slice(cs, func(i, j int) bool { return cs[i].FullCommand() < cs[j].FullCommand() })

		l := 7

		for _, cmd := range cs {
			c := cmd.FullCommand()

			if len(c) > l {
				l = len(c)
			}
		}

		for _, cmd := range cs {
			e.Writer.Writef(fmt.Sprintf(fmt.Sprintf("<h1>%%-%ds</h1>  <value>%%s</value>\n", l), cmd.FullCommand(), cmd.Description)) // nolint:errcheck
		}

		return nil
	}
}

func helpCommand(e *Engine, cmd *Command) {
	e.Writer.Writef("<h2>USAGE</h2>\n  <value>%s</value> <info>%s</info>\n\n", cmd.FullCommand(), cmd.Usage) //nolint:errcheck
	e.Writer.Writef("<h2>DESCRIPTION</h2>\n  <value>%s</value>\n\n", cmd.Description)                        //nolint:errcheck
	e.Writer.Writef("<h2>OPTIONS</h2>\n")                                                                    //nolint:errcheck

	ll := 0
	ls := 0

	fs := []Flag(cmd.Flags)

	sort.Slice(fs, func(i, j int) bool { return fs[i].Name < fs[j].Name })

	for _, f := range fs {
		l := f.UsageLong()
		s := f.UsageShort()

		if len(l) > ll {
			ll = len(l)
		}

		if len(s) > ls {
			ls = len(s)
		}
	}

	for _, f := range fs {
		e.Writer.Writef(fmt.Sprintf(fmt.Sprintf("  %%-%ds  %%-%ds\n", ll, ls), f.UsageLong(), f.UsageShort())) //nolint:errcheck
	}
}
