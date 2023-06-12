package builder

import (
	"fmt"
	"github.com/eurozulu/pempal/keys"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/templates"
)

type Builder interface {
	ApplyTemplate(tp ...templates.Template) error
	Validate() error
	Build() (model.Resource, error)
}

func NewBuilder(t model.ResourceType, keys keys.Keys) (Builder, error) {
	switch t {
	case model.Certificate:
		return &certificateBuilder{keys: keys}, nil
	case model.CertificateRequest:
		return &certificateRequestBuilder{keys: keys}, nil
	case model.PrivateKey:
		return &keyBuilder{}, nil

	default:
		return nil, fmt.Errorf("%ss can not be built", t.String())
	}
}
