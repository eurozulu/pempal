package utils

import (
	"bytes"
	"io"
	"strings"
)

type ColumnWriter struct {
	Columns      []Column
	ColumnSpacer string
	Out          io.Writer
}

type Column struct {
	Name      string
	Title     string
	Alignment ColumnAlignment
	Width     int
}

type ColumnAlignment byte

const (
	Left ColumnAlignment = iota
	Right
	Centre
)

var DefaultColumn = Column{
	Name:      "",
	Alignment: Left,
	Width:     20,
}

func (cw ColumnWriter) Write(p []byte) (n int, err error) {
	return cw.Out.Write(p)
}

func (cw ColumnWriter) WriteString(s string) (n int, err error) {
	return cw.Out.Write([]byte(s))
}

func (cw ColumnWriter) WriteStrings(ss []string) (n int, err error) {
	buf := bytes.NewBuffer(nil)
	for i, sz := range ss {
		if i > 0 && cw.ColumnSpacer != "" {
			buf.WriteString(cw.ColumnSpacer)
		}
		col := cw.ColumnAtIndex(i)
		buf.WriteString(col.FormatString(sz))
	}
	buf.WriteRune('\n')
	return cw.Out.Write(buf.Bytes())
}

func (cw ColumnWriter) ColumnNames() []string {
	names := make([]string, len(cw.Columns))
	for i, c := range cw.Columns {
		title := c.Title
		if title == "" {
			title = c.Name
		}
		names[i] = CapitaliseString(title)
	}
	return names
}

func (cw ColumnWriter) ColumnAtIndex(index int) Column {
	if index >= 0 && index < len(cw.Columns) {
		return cw.Columns[index]
	}
	return DefaultColumn
}

func (cw ColumnWriter) IndexOfColumn(name string) int {
	for i, c := range cw.Columns {
		if strings.EqualFold(c.Name, name) {
			return i
		}
	}
	return -1
}

func (col Column) FormatString(s string) string {
	if col.Width < 0 {
		return s
	}
	if col.Width == 0 {
		return ""
	}

	s = strings.TrimSpace(s)
	padSize := col.Width - len(s)
	if padSize < 0 {
		// trim if too long
		return strings.Join([]string{s[:col.Width-3], "..."}, "")
	}
	if padSize == 0 {
		return s
	}
	pad := strings.Repeat(" ", padSize)
	switch col.Alignment {
	case Left:
		return strings.Join([]string{s, pad}, "")
	case Centre:
		psize := padSize / 2
		lpad := strings.Repeat(" ", psize)
		if psize > 1 {
			pad = pad[psize:]
		}
		return strings.Join([]string{lpad, s, pad}, "")
	case Right:
		return strings.Join([]string{pad, s}, "")
	}
	return s
}
