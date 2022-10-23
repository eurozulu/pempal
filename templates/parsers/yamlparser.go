package parsers

import (
	"pempal/templates"
)

type YAMLParser struct {
}

func (Y YAMLParser) Parse(by []byte) (templates.Template, error) {
	return templates.NewTemplate(by), nil
}
