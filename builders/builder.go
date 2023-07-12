package builders

import (
	"fmt"
	"github.com/eurozulu/pempal/identity"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/utils"
)

type Builder interface {
	AddTemplate(t ...templates.Template)
	Validate() utils.CompoundErrors
	BuildTemplate() templates.Template
	Build() (resources.Resource, error)
}

func NewBuilder(resourceType resources.ResourceType) (Builder, error) {
	switch resourceType {
	case resources.PrivateKey:
		return &keyBuilder{}, nil
	case resources.Certificate:
		return &certificateBuilder{}, nil

	default:
		return nil, fmt.Errorf("a %s has can not be built", resourceType.String())
	}
}

func NewSigningBuilder(resourceType resources.ResourceType, issuerz identity.Issuers) (Builder, error) {
	switch resourceType {
	case resources.Certificate:
		return &certificateBuilder{knownIssuers: issuerz}, nil

	default:
		return nil, fmt.Errorf("a %s has can not be built", resourceType.String())
	}
}
