package builders

import (
	"bytes"
	"fmt"
	"pempal/model"
	"pempal/templates"
)

type ResourceBuilder interface {
	ApplyTemplate(tp ...templates.Template) error
	Validate() []error
	Build() (model.PEMResource, error)
}

func NewResourceBuilder(t model.ResourceType) (ResourceBuilder, error) {
	switch t {
	case model.Certificate:
		return &CertificateBuilder{}, nil
	case model.CertificateRequest:
		return &CertificateRequestBuilder{}, nil
	case model.PrivateKey:
		return &PrivateKeyBuilder{}, nil
	case model.RevokationList:
		return &RevokationListBuilder{}, nil
	default:
		return nil, fmt.Errorf("Invalid resource type. Can not build %s", t.String())
	}
}

func collectErrorList(errs []error) string {
	buf := bytes.NewBuffer(nil)
	for i, err := range errs {
		if i > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(err.Error())
	}
	return buf.String()
}
