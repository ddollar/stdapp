package stdcli

import (
	"sort"
)

func writeFlags(ctx Context, e *Engine, title string, flags []Flag) {
	if len(flags) == 0 {
		return
	}

	e.Writer.Writef("<h2>%s</h2>\n", title) //nolint:errcheck

	// Sort flags alphabetically by name
	sorted := make([]Flag, len(flags))
	copy(sorted, flags)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Name < sorted[j].Name })

	cw := ctx.Columns()

	for _, f := range sorted {
		cw.Append("", f.Usage(), f.Description)
	}

	cw.Print()

	e.Writer.Writef("\n") //nolint:errcheck
}

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
			e.Writer.Writef("<h1>%-*s</h1>  <value>%s</value>\n", l, cmd.FullCommand(), cmd.Description) // nolint:errcheck
		}

		return nil
	}
}

func helpCommand(ctx Context, e *Engine, cmd *Command) {
	e.Writer.Writef("<h2>USAGE</h2>\n  <value>%s</value> <info>%s</info>\n\n", cmd.FullCommand(), cmd.Usage) //nolint:errcheck
	e.Writer.Writef("<h2>DESCRIPTION</h2>\n  <value>%s</value>\n\n", cmd.Description)                        //nolint:errcheck

	writeFlags(ctx, e, "OPTIONS", cmd.Flags)
	writeFlags(ctx, e, "GLOBAL OPTIONS", e.Flags)
}
