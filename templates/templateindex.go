package templates

import (
	"context"
	"fmt"
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/resourcefiles"
	"path/filepath"
	"sort"
	"strings"
)

type TemplateIndex struct {
	index map[string]*model.TemplateFile
}

func (idx TemplateIndex) Names() []string {
	names := make([]string, len(idx.index))
	for n := range idx.index {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func (idx TemplateIndex) ByName(name string) (Template, error) {
	n, err := idx.LookupName(name)
	if err != nil {
		return nil, err
	}
	t, ok := idx.index[n]
	if !ok {
		return nil, fmt.Errorf("template %q is not known", name)
	}
	return t, nil
}

func (idx TemplateIndex) LookupName(name string) (string, error) {
	if _, ok := idx.index[name]; ok {
		return name, nil
	}
	if n, err := LookupBaseName(name); err == nil {
		return n, nil
	}
	return "", fmt.Errorf("template '%s' is not known", name)
}

func (tz TemplateIndex) build(root string) {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()

	for _, name := range BaseTemplateNames() {
		tz.index[name] = nil
	}

	tempFiles := resourcefiles.TemplateFiles(root).Find(ctx, nil)
	for file := range tempFiles {
		path := strings.TrimSuffix(file.Path, filepath.Ext(file.Path))
		var name string
		for {
			name = filepath.Join(filepath.Base(path), name)
			path = filepath.Dir(path)
			if _, err := tz.LookupName(name); err != nil {
				break
			}
			if path == "." {
				logging.Warning("duplicate template name %s ", name)
				break
			}
		}
		tz.index[name] = file
	}
}

func NewTemplateIndex(root string) *TemplateIndex {
	idx := &TemplateIndex{
		index: make(map[string]*model.TemplateFile),
	}
	idx.build(root)
	return idx
}
