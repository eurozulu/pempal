package templates

import (
	"pempal/resources"
)

var resourceTemplates = map[resources.ResourceType]newTemplateFunc{
	resources.Unknown:        newEmptyTemplate,
	resources.Key:            newKeyTemplate,
	resources.Certificate:    newCertificateTemplate,
	resources.Request:        newCSRTemplate,
	resources.RevocationList: newCRLTemplate,
}

type newTemplateFunc func() Template

func newKeyTemplate() Template {
	return &KeyTemplate{}
}
func newCertificateTemplate() Template {
	return &CertificateTemplate{}
}
func newCSRTemplate() Template {
	return &CSRTemplate{}
}
func newCRLTemplate() Template {
	return &CRLTemplate{}
}
func newEmptyTemplate() Template {
	return &EmptyTemplate{}
}
