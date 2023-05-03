package templates

import (
	"fmt"
	"os"
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

	// AddTemplate adds a new template to the store under the given name.
	// returns error if the name already exists.
	AddTemplate(name string, t Template) error

	// RemoveTemplate removes a named template from the store.
	//  returns error if the name is not known
	RemoveTemplate(name string) error
}

type templateManager struct {
	store    BlobStore
	defaults map[string][]byte
}

func (tm templateManager) ParseTemplate(data []byte) (Template, error) {
	tags, raw := parseTags(data)
	extendTemplates, err := tm.TemplatesByName(tagValues(tags.TagsByName(TAG_EXTENDS))...)
	if err != nil {
		return nil, err
	}
	imports, err := tm.namedTemplatesMap(tagValues(tags.TagsByName(TAG_IMPORTS)))
	if err != nil {
		return nil, err
	}
	return newYamlTemplate(tags, raw, extendTemplates, imports)
}

func (tm templateManager) TemplatesByName(names ...string) ([]Template, error) {
	temps := make([]Template, len(names))
	for i, name := range names {
		by, err := tm.store.Read(name)
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to open template '%s'  %v", name, err)
			}
			if len(tm.defaults) > 0 {
				b, ok := tm.defaults[name]
				if !ok {
					return nil, os.ErrNotExist
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

func (tm templateManager) AddTemplate(name string, t Template) error {
	return tm.store.Write(name, t.Raw())
}

func (tm templateManager) RemoveTemplate(name string) error {
	return tm.store.Delete(name)
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

func NewTemplateManager(rootpath string, defaults map[string][]byte) (TemplateManager, error) {
	store, err := NewFileBlobStore(rootpath, "yaml", "template")
	if err != nil {
		return nil, err
	}
	return &templateManager{
		store:    store,
		defaults: defaults,
	}, nil
}
