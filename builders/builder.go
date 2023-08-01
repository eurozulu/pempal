package builders

import (
	"fmt"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/utils"
)

type Builder interface {
	Validate(t templates.Template) utils.CompoundErrors
	Build(t templates.Template) (resources.Resource, error)
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
