package templates

import (
	"fmt"
	"github.com/eurozulu/pempal/model"
	"gopkg.in/yaml.v2"
)

type Template interface {
	fmt.Stringer
}

// MergeTemplates merges the given templates into a single template.
// If the given templates contains a base template, all others are merged into that base template.
// The given templates must NOT contain more than one base template as base templates can not be merged into one another.
// If none of the given templates are a base template, templates are simply merged into a TemplateFile with no path.
func MergeTemplates(temps []Template) (Template, error) {
	if len(temps) == 0 {
		return nil, nil
	}
	base := findFirstBase(temps)
	if base == nil {
		return nil, fmt.Errorf("no base template found")
	}
	for _, t := range temps {
		if base == t {
			continue
		}
		if IsBaseTemplate(t) {
			return nil, fmt.Errorf("cannot merge multiple base templates")
		}
		if err := yaml.Unmarshal([]byte(t.String()), base); err != nil {
			return nil, err
		}
	}
	data, err := yaml.Marshal(base)
	if err != nil {
		return nil, err
	}
	if base == nil {
		return &model.TemplateFile{
			Path: "",
			Data: data,
		}, nil
	}
	if err := yaml.Unmarshal(data, base); err != nil {
		return nil, err
	}
	return base, nil
}

func findFirstBase(temps []Template) Template {
	for _, t := range temps {
		if IsBaseTemplate(t) {
			return t
		}
	}
	return nil
}

func TemplatesOfResources(res []model.PemResource) ([]Template, error) {
	var found []Template
	for _, r := range res {
		t, err := TemplateOfResource(r)
		if err != nil {
			return nil, err
		}
		found = append(found, t)
	}
	return found, nil
}

func TemplateOfResource(res model.PemResource) (Template, error) {
	switch r := res.(type) {
	case *model.Certificate:
		return NewCertificateTemplate(r), nil
	case *model.PrivateKey:
		return NewPrivateKeyTemplate(r), nil
	case *model.CertificateRequest:
		return NewCertificateRequestTemplate(r), nil
	case *model.RevocationList:
		return NewRevocationListTemplate(r), nil
	default:
		return nil, fmt.Errorf("failed to create template as unknown resource type: %T", r.ResourceType())
	}
}
