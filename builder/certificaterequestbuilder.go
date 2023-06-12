package builder

import (
	"github.com/eurozulu/pempal/keys"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/templates"
)

type certificateRequestBuilder struct {
	keys keys.Keys
}

func (c certificateRequestBuilder) ApplyTemplate(tp ...templates.Template) error {
	//TODO implement me
	panic("implement me")
}

func (c certificateRequestBuilder) Validate() error {
	//TODO implement me
	panic("implement me")
}

func (c certificateRequestBuilder) Build() (model.Resource, error) {
	//TODO implement me
	panic("implement me")
}
