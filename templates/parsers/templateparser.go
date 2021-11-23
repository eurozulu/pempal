package parsers

import (
	"encoding/pem"
)

const PEM_TYPE = "type"

type TemplateParser interface {
	Parse(b *pem.Block) (map[string]interface{}, error)
}
