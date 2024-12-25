package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed included
var standardTemplates embed.FS

var baseTemplates = []Template{
	&CertificateTemplate{},
	&KeyTemplate{},
	&CSRTemplate{},
	&CRLTemplate{},
}

const templateExt = ".yaml"

type TemplateLib interface {
	// GetTemplateNames gets all the known template names.
	GetTemplateNames() []string
	GetTemplates(name string) ([]Template, error)
	BaseTemplate(name string) Template
}

type templateLib struct {
	tempFS        LayerFs
	baseTemplates []Template
}

func (tl templateLib) GetTemplateNames() []string {
	names := tl.BaseTemplateNames()
	unique := make(map[string]bool)
	err := fs.WalkDir(tl.tempFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || path == "." {
			return err
		}
		// ignore 'hidden' files &dirs
		if strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if !strings.EqualFold(filepath.Ext(d.Name()), templateExt) {
			return nil
		}
		name := path[:len(path)-len(filepath.Ext(path))]
		if unique[name] {
			return nil
		}
		unique[name] = true
		names = append(names, name)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return names
}

func (tl templateLib) GetTemplates(name string) ([]Template, error) {
	n, err := tl.ResolveName(name)
	if err != nil {
		return nil, err
	}

	// check if it can open as template file
	fn := strings.Join([]string{n, templateExt}, "")
	var tps []Template
	if tl.existsAsFile(fn) {
		data, err := fs.ReadFile(tl.tempFS, fn)
		if err != nil {
			return nil, err
		}
		tps = []Template{
			&template{
				name: n,
				data: data,
			},
		}
	} else if bt := tl.BaseTemplate(n); bt != nil {
		// name is a known base template
		tps = []Template{bt}
	} else {
		return nil, fs.ErrNotExist
	}
	// iterate up the template name, looking for parent templates
	pn := filepath.Dir(n)
	if pn != "." {
		ptps, err := tl.GetTemplates(pn)
		if err != nil {
			return nil, err
		}
		tps = append(ptps, tps...)
	}
	return tps, nil
}

// ResolveName attempts to match the given name to a full template name.
// If the given name has no preceeding path it will be matched to the template with that name
func (tl templateLib) ResolveName(name string) (string, error) {
	names := tl.GetTemplateNames()
	nameLen := len(strings.Split(name, string(filepath.Separator)))

	var found []string
	for _, path := range names {
		sp := strings.Split(path, string(filepath.Separator))
		if len(sp) < nameLen {
			continue
		}
		// trim off leading path so same size as name path
		if !strings.EqualFold(strings.Join(sp[len(sp)-nameLen:], string(filepath.Separator)), name) {
			continue
		}
		found = append(found, path)
	}
	if len(found) == 0 {
		return "", fmt.Errorf("%s not a known template", name)
	}
	if len(found) > 1 {
		//if i := slices.Index(found, name); i >= 0 {
		//	// more than one, check if an exact match name if found
		//	found = []string{found[i]}
		//} else {
		return "", fmt.Errorf("%s is ambiguous: %s", name, strings.Join(found, ", "))
		//}
	}
	return found[0], nil
}

func (tl templateLib) BaseTemplateNames() []string {
	names := make([]string, len(tl.baseTemplates))
	for i, t := range tl.baseTemplates {
		names[i] = t.Name()
	}
	return names
}

func (tl templateLib) BaseTemplate(name string) Template {
	for _, t := range tl.baseTemplates {
		if t.Name() == name {
			return t
		}
	}
	return nil
}

func (tl templateLib) existsAsFile(path string) bool {
	fi, err := fs.Stat(tl.tempFS, path)
	if err != nil {
		return false
	}
	return !fi.IsDir()
}

func (tl templateLib) existsAsDirectory(name string) bool {
	fi, err := fs.Stat(tl.tempFS, name)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func NewTemplateLib(templatePath ...string) TemplateLib {
	stdlib, err := fs.Sub(standardTemplates, "included")
	if err != nil {
		log.Fatal(err)
	}
	tempFS := LayerFs{stdlib}
	for _, root := range templatePath {
		tempFS = append(tempFS, os.DirFS(root))
	}
	return &templateLib{
		tempFS:        tempFS,
		baseTemplates: baseTemplates,
	}
}
