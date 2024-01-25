package stdcli

import (
	"fmt"
	"strings"
)

type TableWriter interface {
	Append(row ...string)
	Print() error
}

type tableWriter struct {
	ctx     Context
	columns []string
	rows    [][]string
}

var _ TableWriter = &tableWriter{}

func (t *tableWriter) Append(row ...string) {
	t.rows = append(t.rows, row)
}

func (t *tableWriter) Print() error {
	f := t.formatString()

	t.ctx.Writef(fmt.Sprintf("<h1>%s</h1>\n", f), t.columns) //nolint:errcheck

	for _, r := range t.rows {
		t.ctx.Writef(fmt.Sprintf("<value>%s</value>\n", f), r) //nolint:errcheck
	}

	return nil
}

func (t *tableWriter) formatString() string {
	f := []string{}

	ws := t.widths()

	for _, w := range ws {
		f = append(f, fmt.Sprintf("%%-%ds", w))
	}

	return strings.Join(f, "  ")
}

func (t *tableWriter) widths() []int {
	w := make([]int, len(t.columns))

	for i, c := range t.columns {
		w[i] = len(stripTag(c))

		for _, r := range t.rows {
			if len(r) > i {
				if lri := len(stripTag(r[i])); lri > w[i] {
					w[i] = lri
				}
			}
		}
	}

	w[len(w)-1] = 0

	return w
}
