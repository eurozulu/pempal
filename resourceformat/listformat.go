package resourceformat

import (
	"github.com/eurozulu/colout"
	"github.com/eurozulu/pempal/model"
	"io"
	"strings"
)

type ListFormat struct{}

var listColumns = buildColumns()

func (l ListFormat) Format(out io.Writer, p *model.PemFile) error {
	outCols := colout.ColumnWriter{
		Columns:         listColumns,
		ColumnSpacer:    "  ",
		StringDelimiter: "---",
	}
	for _, r := range p.Resources() {
		cols := []string{
			r.ResourceType().String(),
			r.String(),
			p.Path,
		}
		if _, err := outCols.WriteString(strings.Join(cols, "---")); err != nil {
			return err
		}
	}
	return nil
}

func buildColumns() []colout.Column {
	return []colout.Column{
		{
			Name:      "type",
			Alignment: colout.Left,
			Width:     15,
		},
		{
			Name:      "details",
			Alignment: colout.Left,
			Width:     90,
		},
		{
			Name:      "path",
			Alignment: colout.Left,
			Width:     -1,
		},
	}
}
