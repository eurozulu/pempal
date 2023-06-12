package builder

import (
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/templates"
)

type keyBuilder struct {
	dto model.PrivateKeyDTO
}

func (k keyBuilder) ApplyTemplate(tp ...templates.Template) error {
	//TODO implement me
	panic("implement me")
}

func (k keyBuilder) Validate() error {
	//TODO implement me
	panic("implement me")
}

func (k keyBuilder) Build() (model.Resource, error) {
	//TODO implement me
	panic("implement me")
}
