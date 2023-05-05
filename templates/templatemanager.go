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

	Names(s ...string) []string

	MergeTemplatesInto(dst interface{}, names ...string) error

	// AddTemplate adds a new template to the store under the given name.
	// returns error if the name already exists.
	AddTemplate(name string, t Template) error

	// RemoveTemplate removes a named template from the store.
	//  returns error if the name is not known
	RemoveTemplate(name string) error

	AddDefaultTemplate(name string, data []byte)
}

type templateManager struct {
	store          BlobStore
	defaults       map[string][]byte
	typedTemplates map[string][]byte
}

func (tm templateManager) ParseTemplate(data []byte) (Template, error) {
	tags, _ := parseTags(data)
	extendTemplates, err := tm.TemplatesByName(tagValues(tags.TagsByName(TAG_EXTENDS))...)
	if err != nil {
		return nil, err
	}
	imports, err := tm.namedTemplatesMap(tagValues(tags.TagsByName(TAG_IMPORTS)))
	if err != nil {
		return nil, err
	}
	return newYamlTemplate(data, extendTemplates, imports)
}

func (tm templateManager) TemplatesByName(names ...string) ([]Template, error) {
	temps := make([]Template, len(names))
	for i, name := range names {
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
		if hasQuery && !constainsString(n, s) {
			continue
		}
		names = append(names, n)
	}
	for k := range tm.defaults {
		if constainsString(k, names) || (hasQuery && !constainsString(k, s)) {
			continue
		}
		names = append(names, k)
	}
	return names
}

func (tm templateManager) MergeTemplatesInto(dst interface{}, names ...string) error {
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

func (tm templateManager) AddTemplate(name string, t Template) error {
	return tm.store.Write(name, t.Raw())
}

func (tm templateManager) RemoveTemplate(name string) error {
	return tm.store.Delete(name)
}

func (tm *templateManager) AddDefaultTemplate(name string, data []byte) {
	tm.defaults[strings.ToLower(name)] = data
}

func (tm templateManager) namedTemplatesMap(names []string) (map[string]interface{}, error) {
	imports := map[string]interface{}{}
	importTemplates, err := tm.TemplatesByName(names...)
	if err != nil {
		return nil, err
	}
	for i, name := range names {
		// place map of each template into the root map
		m := map[string]interface{}{}
		if err = importTemplates[i].Apply(&m); err != nil {
			return nil, err
		}
		imports[name] = m
	}
	return imports, nil
}

func constainsString(s string, ss []string) bool {
	for _, sz := range ss {
		if strings.Contains(s, sz) {
			return true
		}
	}
	return false
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
