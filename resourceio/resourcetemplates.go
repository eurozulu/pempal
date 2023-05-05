package resourceio

import (
	"pempal/model"
	"pempal/templates"
	"strings"
)

func NewResourceTemplateManager(root string) (templates.TemplateManager, error) {
	tm, err := templates.NewTemplateManager(root)
	if err != nil {
		return nil, err
	}
	addDefaultForResourceType(model.PublicKey, tm)
	addDefaultForResourceType(model.PrivateKey, tm)
	addDefaultForResourceType(model.Certificate, tm)
	addDefaultForResourceType(model.CertificateRequest, tm)
	addDefaultForResourceType(model.RevokationList, tm)

	for name, data := range standardTemplates {
		tm.AddDefaultTemplate(name, data)
	}
	return tm, nil
}

func addDefaultForResourceType(r model.ResourceType, tm templates.TemplateManager) {
	p := strings.Join([]string{"#type", r.String()}, " ")
	tm.AddDefaultTemplate(r.String(), []byte(p))
}
