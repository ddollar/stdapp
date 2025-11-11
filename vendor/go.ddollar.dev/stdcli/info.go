package stdcli

import (
	"encoding/json"
	"fmt"
	"strings"
)

type InfoWriter interface {
	Add(header string, value any)
	Print() error
}

type infoWriter struct {
	ctx  Context
	rows []infoRow
}

type infoRow struct {
	header string
	value  string
}

func (i *infoWriter) Add(header string, value any) {
	i.rows = append(i.rows, infoRow{header: header, value: fmt.Sprintf("%v", value)})
}

func (i *infoWriter) Print() error {
	switch i.ctx.Flags().String("output") {
	case "json":
		return i.printJSON()
	default:
		return i.printText()
	}
}

func (i *infoWriter) printJSON() error {
	v := map[string]any{}

	for _, r := range i.rows {
		v[strings.ToLower(r.header)] = r.value
	}

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err //nowrap
	}

	if _, err := i.ctx.Write(data); err != nil {
		return err //nowrap
	}

	return nil
}

func (i *infoWriter) printText() error {
	f := i.formatString()

	for _, r := range i.rows {
		value := strings.Replace(r.value, "\n", fmt.Sprintf("\n%*s  ", i.headerWidth(), ""), -1)
		i.ctx.Writef(f, strings.ToUpper(r.header), value) //nolint:errcheck
	}

	return nil
}

func (i *infoWriter) formatString() string {
	return fmt.Sprintf("<h1>%%-%ds</h1>  <value>%%s</value>\n", i.headerWidth())
}

func (i *infoWriter) headerWidth() int {
	w := 0

	for _, r := range i.rows {
		if len(r.header) > w {
			w = len(r.header)
		}
	}

	return w
}
