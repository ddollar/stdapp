package stdcli

import (
	"encoding/json"
	"fmt"
	"strings"
)

type TableWriter interface {
	Append(row ...any)
	Print() error
}

type tableWriter struct {
	ctx     Context
	columns []any
	rows    [][]any
}

var _ TableWriter = &tableWriter{}

func (t *tableWriter) Append(row ...any) {
	t.rows = append(t.rows, row)
}

func (t *tableWriter) Print() error {
	switch t.ctx.Flags().String("output") {
	case "json":
		return t.printJSON()
	default:
		return t.printText()
	}
}

func (t *tableWriter) printJSON() error {
	lccs := []string{}

	for _, c := range t.columns {
		lccs = append(lccs, strings.ToLower(c.(string)))
	}

	vs := []map[string]any{}

	for _, r := range t.rows {
		v := map[string]any{}

		for i := range t.columns {
			v[lccs[i]] = r[i]
		}

		vs = append(vs, v)
	}

	data, err := json.MarshalIndent(vs, "", "  ")
	if err != nil {
		return err //nowrap
	}

	if _, err := t.ctx.Write(data); err != nil {
		return err //nowrap
	}

	return nil
}

func (t *tableWriter) printText() error {
	cw := t.ctx.Columns()

	cs := make([]any, len(t.columns))

	for i, c := range t.columns {
		cs[i] = fmt.Sprintf("<h1>%s</h1>", c)
	}

	cw.Append(cs...)

	for _, r := range t.rows {
		cw.Append(r...)
	}

	return cw.Print()
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
