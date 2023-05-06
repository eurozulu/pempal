package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type TemplateManager interface {
	// ParseTemplate parses the given raw bytes into a Template, processing any #tags at the beginning of the bytes.
	// Tags:
	// '#extends' tags should contain a single named template which to merge this template into.
	// Multiple #extends tags may be given, each are applied in the order they appear in the bytes.
	// (i.e. last one overrides any clashing property names)
	//
	// '#imports' tags can specify a single named template to make available to {{ marco }} expressions.
	// multiple imports can be given, each referenced by the template name.
	// e.g. #imports mytemplate
	// Will allow macros to reference a property in that template with 'mytemplate.thatproperty'
	// An alias name of imports can be given, should the template name clash with a property name.
	// The alias name should preced the template name, delimited with a space.
	// e.g. #imports othername mytemplate
	ParseTemplate(data []byte) (Template, error)

	// TemplateByName retrieves a named template or an error if the given name is not known
	TemplatesByName(name ...string) ([]Template, error)

	AddTemplate(name string, data []byte)
}

type templateManager struct {
	store          BlobStore
	defaults       map[string][]byte
	typedTemplates map[string][]byte
}

func (tm templateManager) ParseTemplate(data []byte) (Template, error) {
	tags, parsed := parseTags(data)

	if containsGoTemplates(data) {
		importData, err := tm.buildImportData(tags.TagsByName(TAG_IMPORTS))
		if err != nil {
			return nil, err
		}
		parsed, err = executeGoTemplate(data, importData)
		if err != nil {
			return nil, fmt.Errorf("failed to execute templates  %v", err)
		}
	}

	extTags := tags.TagsByName(TAG_EXTENDS)
	var extends []Template
	if len(extTags) > 0 {
		ext, err := tm.TemplatesByName(tagValues(extTags)...)
		if err != nil {
			return nil, err
		}
		extends = ext
	}
	return newYamlTemplate(data, tags, parsed, extends)
}

func (tm templateManager) TemplatesByName(names ...string) ([]Template, error) {
	temps := make([]Template, len(names))
	for i, name := range names {
		by, err := tm.readTemplate(name)
		if err != nil {
			return nil, err
		}
		t, err := tm.ParseTemplate(by)
		if err != nil {
			return nil, err
		}
		temps[i] = t
	}
	return temps, nil
}

func (tm templateManager) Names(s ...string) []string {
	var names []string
	hasQuery := len(s) > 0
	for _, n := range tm.store.Names() {
		if strings.EqualFold(filepath.Ext(n), ".yaml") {
			n = n[:len(n)-len(filepath.Ext(n))]
		}
		if hasQuery && !constainsString(n, s, false) {
			continue
		}
		names = append(names, n)
	}
	for k := range tm.defaults {
		if constainsString(k, names, true) || (hasQuery && !constainsString(k, s, false)) {
			continue
		}
		names = append(names, k)
	}
	return names
}

func (tm templateManager) SaveTemplate(name string, t Template) error {
	return tm.store.Write(name, t.Raw())
}

func (tm templateManager) DeleteTemplate(name string) error {
	return tm.store.Delete(name)
}

func (tm *templateManager) AddTemplate(name string, data []byte) {
	tm.defaults[strings.ToLower(name)] = data
}

func (tm templateManager) readTemplate(name string) ([]byte, error) {
	fname := name
	if filepath.Ext(fname) == "" {
		fname = strings.Join([]string{name, "yaml"}, ".")
	}
	by, err := tm.store.Read(fname)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to open template '%s'  %v", name, err)
		}
		// not exits as 'real' template, check if it's a default
		if len(tm.defaults) > 0 {
			b, ok := tm.defaults[strings.ToLower(name)]
			if !ok {
				return nil, fmt.Errorf("template %s is not known", name)
			}
			by = b
		}
	}
	return by, nil
}

func (tm templateManager) buildImportData(tags Tags) (map[string]interface{}, error) {
	mapped := map[string]interface{}{}
	for _, tag := range tags {
		name, key := tag.ParseAsImport()
		imported, err := tm.TemplatesByName(name)
		if err != nil {
			return nil, err
		}
		m := map[string]interface{}{}
		if err = imported[0].Apply(&m); err != nil {
			return nil, err
		}
		mapped[key] = m
	}
	return mapped, nil
}

func constainsString(s string, ss []string, exact bool) bool {
	for _, sz := range ss {
		if exact {
			if sz == s {
				return true
			}
			continue
		}
		if strings.Contains(s, sz) {
			return true
		}
	}
	return false
}

func ApplyTemplatesTo(dst interface{}, templates []Template) error {
	for _, t := range templates {
		if err := t.Apply(dst); err != nil {
			return err
		}
	}
	return nil
}

func NewTemplateManager(rootpath string) (TemplateManager, error) {
	store, err := NewFileBlobStore(rootpath, "yaml", "template")
	if err != nil {
		return nil, err
	}
	return &templateManager{
		store:    store,
		defaults: map[string][]byte{},
	}, nil
}
