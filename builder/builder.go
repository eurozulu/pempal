package builder

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/pempal/keys"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/templates"
)

type Builder interface {
	Validate(t templates.Template) []error
	Build(t templates.Template) (model.Resource, error)
}

func NewBuilder(rt model.ResourceType, kez keys.Keys) (Builder, error) {
	switch rt {
	case model.Certificate:
		return &certificateBuilder{keys: kez}, nil
	//case model.CertificateRequest:
	//	return &certificateRequestBuilder{keys: keys}, nil
	case model.PrivateKey:
		return &keyBuilder{}, nil

	default:
		return nil, fmt.Errorf("invalid resource type.  %ss can not be built", rt.String())
	}
}

func CombineErrors(errs []error) error {
	buf := bytes.NewBuffer(nil)
	for i, err := range errs {
		if i > 0 {
			buf.WriteRune('\n')
		}
		buf.WriteString(err.Error())
	}
	return fmt.Errorf(buf.String())
}
