package stdcli

import (
	"fmt"
	"strings"
)

type InfoWriter interface {
	Add(header, value string)
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

func (i *infoWriter) Add(header, value string) {
	i.rows = append(i.rows, infoRow{header: header, value: value})
}

func (i *infoWriter) Print() error {
	f := i.formatString()

	for _, r := range i.rows {
		value := strings.Replace(r.value, "\n", fmt.Sprintf(fmt.Sprintf("\n%%%ds  ", i.headerWidth()), ""), -1)
		i.ctx.Writef(f, r.header, value) //nolint:errcheck
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
