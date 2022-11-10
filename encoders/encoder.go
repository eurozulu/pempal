package encoders

import (
	"encoding/pem"
	"pempal/templates"
)

type Encoder interface {
	Encode(p *pem.Block) (templates.Template, error)
}
