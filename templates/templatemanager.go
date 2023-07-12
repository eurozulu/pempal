package templates

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"path/filepath"
	"strings"
)

const templateFileExtension = ".yaml"

var ErrTemplateNotFound = fmt.Errorf("template is not known")

type TemplateManager interface {
	TemplateByName(name string) (Template, error)
}

type templateManager struct {
	store BlobStore
}

func (tm templateManager) TemplateByName(name string) (Template, error) {
	data, err := tm.readTemplateBytes(name)
	if err != nil {
		return nil, fmt.Errorf("template '%s' %v", name, err)
	}
	return parseTemplateBytes(data)
}

func (tm templateManager) readTemplateBytes(name string) ([]byte, error) {
	fname := name
	if filepath.Ext(fname) == "" {
		fname = strings.Join([]string{name, templateFileExtension}, "")
	}
	return tm.store.Read(fname)
}

func parseTemplateBytes(data []byte) (Template, error) {
	m := map[string]string{}
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func NewTemplateManager(templatePath []string) (TemplateManager, error) {
	bs, err := newFileBlobStore(templatePath[0])
	if err != nil {
		return nil, err
	}
	return &templateManager{store: bs}, nil
}
