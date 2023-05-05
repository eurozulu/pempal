package builders

import (
	"bytes"
	"pempal/model"
	"pempal/templates"
)

type ResourceBuilder interface {
	ApplyTemplate(tp ...templates.Template) error
	Validate() []error
	Build() (model.PEMResource, error)
}

func NewResourceBuilder(t model.ResourceType) ResourceBuilder {
	switch t {
	case model.Certificate:
		return &CertificateBuilder{}
	case model.CertificateRequest:
		return &CertificateRequestBuilder{}

	}
}

func collectErrorList(errs []error, delimit string) string {
	buf := bytes.NewBuffer(nil)
	for i, err := range errs {
		if i > 0 {
			buf.WriteString(delimit)
		}
		buf.WriteString(err.Error())
	}
	return buf.String()
}
