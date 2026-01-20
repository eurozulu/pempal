package resourceformat

import (
	"fmt"
	"github.com/eurozulu/pempal/model"
	"io"
	"strings"
)

const DefaultFormat = "list"

var views = map[string]ResourceFormat{
	"list": &ListFormat{},
	"pem":  &PemFormat{},
	"der":  &DERFormat{},
	"yaml": &YamlFormat{},
	"json": &JsonFormat{},
}

type ResourceFormat interface {
	Format(out io.Writer, p *model.PemFile) error
}

func NewResourceFormat(format string) (ResourceFormat, error) {
	if format == "" {
		format = DefaultFormat
	}
	format = strings.ToLower(format)
	if rv, ok := views[format]; ok {
		return rv, nil
	}
	for k, v := range views {
		if strings.HasPrefix(k, format) {
			return v, nil
		}
	}
	return nil, fmt.Errorf("%q is not a known resource view format", format)
}

func FormatResources(format string, out io.Writer, p *model.PemFile) error {
	vFormat, err := NewResourceFormat(format)
	if err != nil {
		return err
	}
	return vFormat.Format(out, p)
}
