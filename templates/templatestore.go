package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const templateFileExtension = ".yaml"

type TemplateStore interface {

	// TemplateByName retrieves a named template or an error if the given name is not known
	TemplateByName(name string) (Template, error)

	// ExtendedTemplatesByName retrieves a named template and all the templates it extends, or an error if the given name or any extended name is not known.
	// Providing more than one name is the same as extending the second named template with the first.  i.e. a chain of templates can
	// be given OR be fixed by a templates tag.
	// If the given named template extends another template, the extended template is loaded first.  If the extended template also extends
	// another, the process is repeated until a template not extending another is reached.  The unextending template will be the first template returned.
	// the template extending it, next and so on until the given named template, whihc is the last returned.
	ExtendedTemplatesByName(name ...string) ([]Template, error)

	// Names lists the names of all the available templates.
	Names(s ...string) []string

	Exists(name string) bool

	// SaveTemplate adds a new template to the store under the given name.
	// returns error if the name already exists.
	SaveTemplate(name string, t Template) error

	// DeleteTemplate removes a named template from the store.
	//  returns error if the name is not known
	// note: default tempaltes can NOT be deleted.
	DeleteTemplate(name string) error
}

type templateStore struct {
	store    BlobStore
	defaults map[string][]byte
}

func (ts templateStore) TemplateByName(name string) (Template, error) {
	data, err := ts.readTemplateBytes(name)
	if err != nil {
		return nil, fmt.Errorf("template '%s' error %v", name, err)
	}
	return NewTemplate(data)
}

func (ts templateStore) ExtendedTemplatesByName(names ...string) ([]Template, error) {
	var temps []Template
	for _, name := range names {
		ts, err := ts.extendedTemplates(name, uniqueNames{})
		if err != nil {
			return nil, err
		}
		temps = append(temps, ts...)
	}
	return temps, nil
} //TODO: combine extended an d imports as only the tag name differs.

func (ts templateStore) ImportedTemplate(t Template) ([]Template, error) {
	names := strings.Split(t.Tags().TagByName(TAG_IMPORTS).Value(), ",")
	var temps []Template
	for _, name := range names {
		ts, err := ts.extendedTemplates(name, uniqueNames{})
		if err != nil {
			return nil, err
		}
		temps = append(temps, ts...)
	}
	return temps, nil
}

func (ts templateStore) extendedTemplates(name string, used uniqueNames) ([]Template, error) {
	if used[name] {
		return nil, fmt.Errorf("extending template '%s' is a circular reference.", name)
	}
	used[name] = true

	var temps []Template
	tp, err := ts.TemplateByName(name)
	if err != nil {
		return nil, err
	}

	if tp.Tags().ContainsTag(TAG_EXTENDS) {
		// recusrive call with the extended names
		names := strings.Split(tp.Tags().TagByName(TAG_EXTENDS).Value(), ",")
		for _, exn := range names {
			// recursive call using the extended name
			ts, err := ts.extendedTemplates(exn, used)
			if err != nil {
				return nil, err
			}
			temps = append(temps, ts...)
		}
	}
	temps = append(temps, tp)
	return temps, nil
}

func (ts templateStore) Names(s ...string) []string {
	var names []string
	hasQuery := len(s) > 0
	for _, n := range ts.store.Names() {
		if strings.EqualFold(filepath.Ext(n), templateFileExtension) {
			n = n[:len(n)-len(templateFileExtension)]
		}
		if hasQuery && !compareString(n, s, false) {
			continue
		}
		names = append(names, n)
	}
	for k := range ts.defaults {
		if compareString(k, names, true) || (hasQuery && !compareString(k, s, false)) {
			continue
		}
		names = append(names, k)
	}
	return names
}

func (ts templateStore) SaveTemplate(name string, t Template) error {
	if filepath.Ext(name) == "" {
		name = strings.Join([]string{name, templateFileExtension}, "")
	}
	return ts.store.Write(name, []byte(t.String()))
}

func (ts templateStore) DeleteTemplate(name string) error {
	if filepath.Ext(name) == "" {
		name = strings.Join([]string{name, templateFileExtension}, "")
	}
	return ts.store.Delete(name)
}

func (ts templateStore) Exists(name string) bool {
	fname := name
	if filepath.Ext(fname) == "" {
		fname = strings.Join([]string{name, templateFileExtension}, "")
	}
	return ts.store.Contains(fname)
}

func (ts templateStore) readTemplateBytes(name string) ([]byte, error) {
	fname := name
	if filepath.Ext(fname) == "" {
		fname = strings.Join([]string{name, templateFileExtension}, "")
	}
	by, err := ts.store.Read(fname)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to open template '%s'  %v", name, err)
		}
		// not exits as 'real' template, check if it's a default
		b, ok := ts.defaults[strings.ToLower(name)]
		if !ok {
			return nil, fmt.Errorf("template %s is not known", name)
		}
		by = b
	}
	return by, nil
}

func compareString(s string, ss []string, exact bool) bool {
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

// NewTemplateStore creates a new Template store using the given root path as
// the location of where to search and save template files.
// An optional map of default templates may be provided.  When given defaults will be served when that name can not be found in the store.
func NewTemplateStore(rootpath string, defaultTemplates map[string][]byte) (TemplateStore, error) {
	store, err := NewFileBlobStore(rootpath, "yaml", "template")
	if err != nil {
		return nil, err
	}
	if defaultTemplates == nil {
		defaultTemplates = map[string][]byte{}
	}
	return &templateStore{
		store:    store,
		defaults: defaultTemplates,
	}, nil
}
