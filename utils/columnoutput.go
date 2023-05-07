package utils

import (
	"bytes"
	"io"
	"strings"
)

const DefaultColumnWidth = 20

type ColumnOutput struct {
	// ColumnWidths lists the custom widths of induvidual columns.
	// All columns are set to DefaultColumnWidth unless an non zero value is in this slice
	ColumnWidths []int

	// Delimiter to split the data being written
	Delimiter string
	// ColumnDelimiter defines how columns are seperated
	ColumnDelimiter string

	out io.Writer
}

func (t ColumnOutput) Write(p []byte) (n int, err error) {
	if strings.TrimSpace(string(p)) == "" {
		return t.out.Write(p)
	}
	return t.WriteSlice(strings.Split(string(p), t.Delimiter))
}

func (t ColumnOutput) WriteSlice(ss []string) (n int, err error) {
	buf := bytes.NewBuffer(nil)
	last := len(ss) - 1
	for i, s := range ss {
		if i > 0 {
			buf.WriteString(t.ColumnDelimiter)
		}
		w := DefaultColumnWidth
		if i < len(t.ColumnWidths) && t.ColumnWidths[i] > 0 {
			w = t.ColumnWidths[i]
		}
		if i < last {
			s = fixSizeString(s, w)
		}
		buf.WriteString(s)

	}
	return t.out.Write(buf.Bytes())
}

func fixSizeString(s string, width int) string {
	if len(s) > width {
		s = s[:width]

	} else if width > len(s) {
		pad := width - len(s)
		s = strings.Join([]string{s, strings.Repeat(" ", pad)}, "")
	}
	return s
}

func NewColumnOutput(out io.Writer, widths ...int) *ColumnOutput {
	return &ColumnOutput{
		ColumnWidths:    widths,
		Delimiter:       ",",
		ColumnDelimiter: "\t",
		out:             out,
	}
}
