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
	return tm, nil
}

func addDefaultForResourceType(r model.ResourceType, tm templates.TemplateManager) {
	p := strings.Join([]string{"resource-type: ", r.String()}, "")
	tm.AddDefaultTemplate(r.String(), []byte(p))
}

func MergeTemplatesInto(dst interface{}, templatesroot string, names ...string) error {
	tm, err := NewResourceTemplateManager(templatesroot)
	if err != nil {
		return err
	}
	tps, err := tm.TemplatesByName(names...)
	if err != nil {
		return err
	}
	for _, t := range tps {
		if err = t.Apply(dst); err != nil {
			return err
		}
	}
	return nil
}
