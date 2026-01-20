package resourceformat

import (
	"github.com/eurozulu/colout"
	"github.com/eurozulu/pempal/model"
	"io"
)

type ListFormat struct{}

var listColumns = buildColumns()

func (l ListFormat) Format(out io.Writer, p *model.PemFile) error {
	outCols := colout.ColumnWriter{
		Columns:      listColumns,
		ColumnSpacer: " ",
		Out:          out,
	}
	for _, r := range p.Resources() {
		cols := []string{
			r.ResourceType().String(),
			r.String(),
			p.Path,
		}
		if _, err := outCols.WriteStrings(cols); err != nil {
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
			Width:     20,
		},
		{
			Name:      "details",
			Alignment: colout.Left,
			Width:     60,
		},
		{
			Name:      "path",
			Alignment: colout.Left,
			Width:     -1,
		},
	}
}
