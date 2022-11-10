package builders

import (
	"encoding/pem"
	"fmt"
	"pempal/pemtypes"
	"pempal/templates"
)

var ErrMissingValue = fmt.Errorf("missing %s")
var ErrInvalidValue = fmt.Errorf("invalid %s")

type Builder interface {
	AddResource(p ...*pem.Block)
	AddTemplate(t ...templates.Template) error
	Validate() []error
	Build() ([]*pem.Block, error)
}

func NewBuilder(pemType pemtypes.PEMType) Builder {

}
