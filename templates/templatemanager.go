package templates

import (
	"bytes"
	"fmt"
	"github.com/go-yaml/yaml"
	"os"
	"strings"
)

const extendsTag = "#extends"
const templateFileExtension = "yaml"

var ErrTemplateNotFound = fmt.Errorf("template is not known")

type TemplateManager interface {
	Names(includeDefaults bool) []string
	TemplateByName(name string) (Template, error)
	AddDefaultTemplate(name string, data []byte) error
}

type templateManager struct {
	defaults ByteStore
	store    ByteStore
}

func (tm templateManager) Names(includeDefaults bool) []string {
	if !includeDefaults {
		return tm.store.Names()
	}
	return removeDuplicates(append(tm.defaults.Names(), tm.store.Names()...))
}

// TemplateByName retrieves the named template, expanded with any 'extends' templates
func (tm templateManager) TemplateByName(name string) (Template, error) {
	return tm.templatesByNames([]string{name}, map[string]bool{})
}

func (tm templateManager) templatesByNames(names []string, usedNames map[string]bool) (Template, error) {
	tb := NewTemplateBuilder()
	for _, name := range names {
		if usedNames[name] {
			return nil, fmt.Errorf("cyclic extends, %s extends itself", name)
		}
		data, err := tm.TemplateSourceByName(name)
		if err != nil {
			return nil, fmt.Errorf("template '%s' %v", name, err)
		}
		usedNames[name] = true
		extends := parseExtendsNames(data)
		if len(extends) > 0 {
			et, err := tm.templatesByNames(extends, usedNames)
			if err != nil {
				return nil, err
			}
			tb.Add(et)
		}
		t, err := parseTemplateBytes(data)
		if err != nil {
			return nil, err
		}
		tb.Add(t)
	}
	return tb.Build(), nil
}

// TemplateSourceByName retrieves the raw bytes of a template by its name
func (tm templateManager) TemplateSourceByName(name string) ([]byte, error) {
	if !tm.store.Contains(name) {
		return tm.defaults.Read(name)
	}
	return tm.store.Read(name)
}

func (tm *templateManager) AddDefaultTemplate(name string, data []byte) error {
	return tm.defaults.Write(name, data)
}

func parseTemplateBytes(data []byte) (Template, error) {
	m := map[string]string{}
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func parseExtendsNames(data []byte) []string {
	if !bytes.HasPrefix(data, []byte(extendsTag)) {
		return nil
	}
	eol := bytes.IndexByte(data, '\n')
	if eol < 0 {
		eol = len(data)
	}
	return strings.Split(strings.TrimSpace(string(data[len(extendsTag):eol])), ",")
}

func NewTemplateManager(templatePath []string) (TemplateManager, error) {
	if len(templatePath) == 0 {
		templatePath = []string{os.ExpandEnv("$PWD")}
	}
	fs, err := newFileStore(templatePath[0], templateFileExtension)

	if err != nil {
		return nil, err
	}
	return &templateManager{
		defaults: newMemoryStore(),
		store:    fs,
	}, nil
}
