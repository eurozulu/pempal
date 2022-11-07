package templates

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"pempal/filepaths"
	"strings"
)

const templatePath = "templates"

var TypedTemplateNames = [...]string{"CERTIFICATE", "KEY", "CSR", "CRL", "NAME"}

type newTypeTemplateFunc func() Template

var typedTemplateFuncs = map[string]newTypeTemplateFunc{
	"CERTIFICATE": NewCertificateTemplate,
	"KEY":         NewKeyTemplate,
	"CSR":         NewCSRTemplate,
	"CRL":         NewCRLTemplate,
	"NAME":        NewNameTemplate,
}

type TemplateManager interface {
	TemplateNames() ([]string, error)
	PropertyNames(t Template) []string
	TemplateByName(name string) (Template, error)
	TemplateByPath(path string) (Template, error)
	MergeTemplates(dst Template, ts ...Template) error
	FindTemplates(names ...string) ([]Template, error)
}

type templateManager struct{}

func (tm templateManager) TemplateNames() ([]string, error) {
	names, err := filepaths.ListFiles(templatePath)
	if err != nil {
		return nil, err
	}
	return append(TypedTemplateNames[:], names...), nil
}

func (tm templateManager) PropertyNames(t Template) []string {
	et, ok := t.(EmptyTemplate)
	if !ok {
		et = EmptyTemplate{}
		tm.MergeTemplates(et, t)
	}
	var names []string
	for k := range et {
		names = append(names, k)
	}
	return names
}

func (tm templateManager) TemplateByPath(path string) (Template, error) {
	by, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	nt := &template{}
	if err := yaml.Unmarshal(by, nt); err != nil {
		return nil, err
	}
	return nt, nil
}

func (tm templateManager) TemplateByName(name string) (Template, error) {
	fn := name
	if filepath.Ext(fn) == "" {
		fn = fmt.Sprintf("%s.yaml", fn)
	}
	// attempt to load a file based one first. (maybe overriding the typed templateManager)
	ft, _ := tm.TemplateByPath(filepath.Join(filepaths.HomePath(), templatePath, fn))

	// Now check if its a type template and apply overriding one, if present
	if i := stringIndex(strings.ToUpper(name), TypedTemplateNames[:]); i >= 0 {
		newTemp := typedTemplateFuncs[TypedTemplateNames[i]]()
		if ft != nil {
			tm.MergeTemplates(newTemp, ft)
		}
		ft = newTemp
	}
	return ft, nil
}

func (tm templateManager) FindTemplates(names ...string) ([]Template, error) {
	var temps []Template
	for _, name := range names {
		nt, err := tm.TemplateByName(name)
		if err != nil || nt == nil {
			// not a named template, try filepath
			nt, err = tm.TemplateByPath(name)
			if err != nil {
				return nil, err
			}
		}
		temps = append(temps, nt)
	}
	return temps, nil
}

func (tm templateManager) MergeTemplates(dst Template, ts ...Template) error {
	for _, t := range ts {
		by, err := yaml.Marshal(t)
		if err != nil {
			return err
		}
		if err = yaml.Unmarshal(by, dst); err != nil {
			return err
		}
	}
	return nil
}

func stringIndex(s string, ss []string) int {
	for i, sz := range ss {
		if s == sz {
			return i
		}
	}
	return -1
}

func NewCertificateTemplate() Template {
	return &CertificateTemplate{}
}

func NewTemplateManager() TemplateManager {
	return &templateManager{}
}
