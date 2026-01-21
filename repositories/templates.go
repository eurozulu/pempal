package repositories

import (
	"context"
	"fmt"
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/resourcefiles"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/tools"
	"path/filepath"
	"strings"
)

type Templates string

type templateFilter resourcefiles.TemplateFileFilter

func (tz Templates) ByName(name ...string) ([]templates.Template, error) {
	var found []templates.Template
	for _, n := range name {
		temps := tz.FindByName(n)
		if len(temps) == 0 {
			return nil, fmt.Errorf("template %s not found", n)
		}
		if len(temps) > 1 {
			return nil, fmt.Errorf("template %s matches multiple templates", name)
		}
		found = append(found, temps...)
	}
	return found, nil
}

func (tz Templates) FindByName(name string) []templates.Template {
	if n, err := templates.LookupBaseName(name); err == nil {
		t, _ := templates.NewBaseTemplate(n)
		return []templates.Template{t}
	}
	return tz.FindAll(func(file *model.TemplateFile) bool {
		path := strings.TrimSuffix(file.Path, filepath.Ext(file.Path))
		return strings.HasSuffix(path, name)
	})
}

func (tfs Templates) ExpandedByName(names ...string) ([]templates.Template, error) {
	if len(names) == 0 {
		return nil, fmt.Errorf("no template names specified")
	}
	namez, err := tfs.ExpandNames(names...)
	if err != nil {
		return nil, err
	}
	logging.Debug("expanded template names: %v to %v", names, namez)
	return tfs.ByName(namez...)
}

func (tz Templates) TemplateNames() []string {
	names := templates.BaseTemplateNames()
	names = tools.AppendUnique(names, templates.NewTemplateIndex(string(tz)).Names()...)
	return names
}

// ExpandNames uses the given names to expand the list of names to include any #expends names.
// each name, is not q base name, has any #extends read. If extends contains names, these names preceed the extended name,
// All names are resolved to the template name.
func (tz Templates) ExpandNames(names ...string) ([]string, error) {
	idx := templates.NewTemplateIndex(string(tz))
	return tz.expandNames(names, idx)
}

func (tz Templates) expandNames(names []string, idx *templates.TemplateIndex) ([]string, error) {
	var found []string
	for _, name := range names {
		n, err := idx.LookupName(name)
		if err != nil {
			return nil, err
		}
		t, err := idx.ByName(n)
		if err != nil {
			return nil, err
		}
		namez := []string{n}
		tf := t.(*model.TemplateFile)
		// check not a base template
		if tf != nil {
			ext := tools.NotIn(tf.Extends(), found)
			if len(ext) > 0 {
				ext, err = tz.expandNames(ext, idx)
				if err != nil {
					return nil, err
				}
				namez = tools.AppendUnique(ext, namez...)
			}
		}
		found = tools.AppendUnique(found, namez...)
	}
	return found, nil
}

func (tz Templates) FindAll(filter templateFilter) []templates.Template {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	var found []templates.Template
	tempFiles := tz.Find(ctx, filter)
	for file := range tempFiles {
		found = append(found, file)
	}
	return found
}

func (tz Templates) Find(ctx context.Context, filter templateFilter) <-chan templates.Template {
	found := make(chan templates.Template)
	go func() {
		defer close(found)
		tempFiles := resourcefiles.TemplateFiles(tz).Find(ctx, resourcefiles.TemplateFileFilter(filter))
		for file := range tempFiles {
			if filter != nil && !filter(file) {
				continue
			}
			select {
			case <-ctx.Done():
				return
			case found <- file:
			}
		}
	}()
	return found
}
