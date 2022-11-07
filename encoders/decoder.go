package encoders

import (
	"encoding/pem"
	"pempal/templates"
)

type Decoder interface {
	Decode(t templates.Template) (*pem.Block, error)
}
