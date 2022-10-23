package parsers

import (
	"encoding/pem"
	"fmt"
	"pempal/templates"
)

var PEMTypes = map[string]templates.TemplateParser{
	"CERTIFICATE":         &DERCertificateParser{},
	"CERTIFICATE REQUEST": &CSRCertificateParser{},
}

type PEMParser struct {
}

func (P PEMParser) Parse(by []byte) (templates.Template, error) {
	data := by
	var blk *pem.Block
	for {
		if len(data) == 0 {
			break
		}
		blk, data = pem.Decode(by)
		if blk == nil {
			return nil, fmt.Errorf("no recognised PEM block found in data")
		}
		parser := PEMTypes[blk.Type]
		if parser == nil {
			continue
		}
		return parser.Parse(blk.Bytes)
	}
	return nil, fmt.Errorf("no recognised PEM block found in data")
}
