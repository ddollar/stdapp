package stdcli

import (
	"fmt"
	"strings"
)

type ColumnWriter interface {
	Append(items ...any)
	Print() error
}

type columnWriter struct {
	ctx  Context
	rows [][]any
}

var _ ColumnWriter = &columnWriter{}

func (c *columnWriter) Append(items ...any) {
	c.rows = append(c.rows, items)
}

func (c *columnWriter) Print() error {
	if len(c.rows) == 0 {
		return nil
	}

	widths := c.widths()

	for _, row := range c.rows {
		parts := []string{}
		for i, item := range row {
			itemStr := fmt.Sprintf("%v", item)
			strippedLen := len(stripTags(item))

			if i < len(widths) && widths[i] > 0 {
				// Calculate how many spaces we need to add for padding
				padding := widths[i] - strippedLen
				if padding > 0 {
					itemStr = itemStr + strings.Repeat(" ", padding)
				}
			}
			parts = append(parts, itemStr)
		}
		c.ctx.Writef("<value>%s</value>\n", strings.Join(parts, "  ")) //nolint:errcheck
	}

	return nil
}

func (c *columnWriter) widths() []int {
	if len(c.rows) == 0 {
		return []int{}
	}

	// Find the maximum number of columns
	maxCols := 0
	for _, row := range c.rows {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	widths := make([]int, maxCols)

	// Calculate max width for each column
	for _, row := range c.rows {
		for i, item := range row {
			if length := len(stripTags(item)); length > widths[i] {
				widths[i] = length
			}
		}
	}

	// Last column doesn't need padding
	if len(widths) > 0 {
		widths[len(widths)-1] = 0
	}

	return widths
}
