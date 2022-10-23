package templates

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"pempal/filepath"
	"strings"
)

type TemplateFilePath interface {
	Path() []string
	Find(names ...string) ([]Template, error)
	Names(ctx context.Context) <-chan string
}

type templateFilePath struct {
	fPath filepath.FilePath
}

func (tp templateFilePath) Path() []string {
	return tp.fPath.Path()
}

// Find locates the named templates.
// The result is an ordered slice of Templates, matching the order of the given names.
// if ANY of the named templates are not found, an error is thrown
func (tp templateFilePath) Find(names ...string) ([]Template, error) {
	found := make([]Template, len(names))
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()

	var foundCount int
	for p := range tp.fPath.Find(ctx, templateFileFilter{names: names}) {
		i := indexOfName(p, names)
		if i < 0 || found[i] != nil {
			continue
		}
		var err error
		found[i], err = LoadTemplate(p)
		foundCount++
		if err != nil {
			return nil, err
		}
	}

	// Check all names found or list the ones not found
	if foundCount < len(names) {
		var unfound []string
		for i, f := range found {
			if f == nil {
				unfound = append(unfound, names[i])
			}
		}
		return nil, fmt.Errorf("failed to find %v", unfound)
	}
	return found, nil
}

func (tp templateFilePath) Names(ctx context.Context) <-chan string {
	return tp.fPath.Find(ctx, templateFileFilter{names: nil})
}

// templateFileFilter is a recurive filefilter which limits files to those with extension 'yaml' or 'template'.
// Optionally, it may contain a list of names which it will match against the templates found.
type templateFileFilter struct {
	names []string
}

func (t templateFileFilter) Accept(p string) bool {
	ext := strings.ToLower(path.Ext(p))
	if ext != "yaml" && ext != "template" {
		return false
	}
	// no names == all names
	if len(t.names) == 0 {
		return true
	}
	return indexOfName(filepath.NameFromPath(p), t.names) >= 0
}

func (t templateFileFilter) AcceptDirectory(p string) bool {
	return !strings.HasPrefix(path.Base(p), ".")
}

func indexOfName(name string, names []string) int {
	for i, n := range names {
		if n == name {
			return i
		}
	}
	return -1
}

// LoadTemplate loads the given file as a template.
// The given file may be any file supported by the TemplateFileTypes.
func LoadTemplate(location string) (Template, error) {
	tr := TemplateFileTypes[strings.ToLower(path.Ext(location))]
	if tr == nil {
		return nil, fmt.Errorf("'%s' is not a known template file types", path.Ext(location))
	}
	by, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file '%s'  %v", location, err)
	}
	return tr.Parse(by)
}

// NewTemplatePath creates a new Path using the given, coma delimited list of directories
// Template path is recursive, so it searches subdirectories of the directories given in the path.
func NewTemplatePath(templatepath string) TemplateFilePath {
	return &templateFilePath{fPath: filepath.NewFilePath(templatepath)}
}
