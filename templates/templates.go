package templates

import (
	"context"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"pempal/fileformats"
	"pempal/filepathscanner"
	"pempal/pemresources"
	"sort"
	"strings"
)

const FileTag = "#"

const ENV_TemplatePath = "PEMPAL_TEMPLATEPATH"

var TemplatePath = strings.TrimSpace(os.Getenv(ENV_TemplatePath))

// TemplateNames lists the known names of all the templates, including buuld in ones
func TemplateNames(includeBuiltIn bool) []string {
	paths := templateFileNames()
	names := make([]string, len(paths))
	var i int
	for k := range paths {
		names[i] = k
	}
	sort.Strings(names)
	// add built in names
	if includeBuiltIn {
		names = append(names, templateInBuiltNames()...)
	}
	return names
}

// CompoundTemplate generates a new temapltes, which is the combination of all the given templates.
func CompoundTemplate(base Template, names ...string) (Template, error) {
	temps, err := Templates(names...)
	if err != nil {
		return nil, err
	}
	if base != nil {
		temps = append([]Template{base}, temps...)
	}
	tb := &TemplateBuilder{Templates: temps}
	return tb.Build()
}

// Templates loads the named templates
// names can be # preceeded names, to indicate template, or resource path to existing resources, to make into template.
func Templates(names ...string) ([]Template, error) {
	var temps []Template
	paths := templateFileNames()
	for _, name := range names {
		key := cleanName(name)
		if _, ok := builtInTemplates[key]; ok {
			t, err := loadBuildInTemplate(key)
			if err != nil {
				return nil, err
			}
			temps = append(temps, t)
			continue
		}
		// check if name is a known template name or just a file path.
		p := name
		if !strings.Contains(p, "/") {
			p = paths[key]
		}
		t, err := loadResourceAsTemplate(p)
		if err != nil {
			return nil, fmt.Errorf("problem with built in template!  %w", err)
		}
		temps = append(temps, t)
	}
	return temps, nil
}

func loadBuildInTemplate(name string) (Template, error) {
	tb, ok := builtInTemplates[strings.ToLower(strings.TrimLeft(name, FileTag))]
	if !ok {
		return nil, fmt.Errorf("unknown template %s", name)
	}
	return ParseTemplate(tb.([]byte))
}

func loadResourceAsTemplate(p string) (Template, error) {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	ts := pemresources.PemScanner{
		Recursive: true,
		Verbose:   false,
		Reader:    fileformats.NewFormatReader(),
	}
	tch := ts.Scan(ctx, p)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case blk, ok := <-tch:
		if !ok {
			return nil, os.ErrNotExist
		}
		return BlockToTemplate(blk)
	}
}

// templateInBuiltNames gets a list of all the templates hard coded into the application
func templateInBuiltNames() []string {
	var names []string
	for k := range builtInTemplates {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// templateFileNames gets a list of all the template names found in the templateKeyPath
func templateFileNames() map[string]string {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	ts := filepathscanner.FilePathScanner{
		Recursive: true,
		ExtFilter: map[string]bool{"yaml": true},
	}
	pathCh := ts.Scan(ctx, GetTemplatePath()...)
	names := map[string]string{}
	for p := range pathCh {
		if !strings.Contains(p, FileTag) {
			continue
		}
		name := cleanName(p)
		names[name] = p
	}
	return names
}

func cleanName(p string) string {
	n := path.Base(p)
	// strip the file extension
	e := path.Ext(n)
	if len(e) > 0 {
		n = n[:len(n)-len(e)]
	}
	n = strings.TrimLeft(n, FileTag)
	return n
}

// GetTemplatePath gets the template path, the path(s) to search for private keys
// Template path is the current working directory with any additional (comma demited) paths set in the ENV_TemplatePath
func GetTemplatePath() []string {
	tp := []string{os.ExpandEnv("$PWD")}
	if TemplatePath != "" {
		tp = append(tp, strings.Split(os.ExpandEnv(TemplatePath), ":")...)
	}
	return tp
}

func ParseTemplate(by []byte) (Template, error) {
	reader := fileformats.NewFormatReader("yaml")
	tblk, err := reader.Unmarshal(by)
	if err != nil {
		return nil, err
	}
	if len(tblk) < 1 {
		return nil, fmt.Errorf("failed to parse template")
	}
	return BlockToTemplate(tblk[0])
}

func BlockToTemplate(blk *pem.Block) (Template, error) {
	location := blk.Headers[pemresources.LocationHeaderKey]
	var t Template
	pt := blk.Type
	isTemplate := strings.HasSuffix(pt, fileformats.PEM_TEMPLATE)
	if isTemplate {
		pt = pt[:len(pt)-(len(fileformats.PEM_TEMPLATE)+1)]
	}
	if fileformats.PemTypesCertificate[pt] {
		t = &pemresources.Certificate{PemResource: pemresources.PemResource{Location: location}}
	}
	if fileformats.PemTypesCertificateRequest[pt] {
		t = &pemresources.CertificateRequest{PemResource: pemresources.PemResource{Location: location}}
	}
	if fileformats.PemTypesPrivateKey[pt] {
		t = &pemresources.PrivateKey{PemResource: pemresources.PemResource{Location: location}}
	}
	if fileformats.PemTypesPublicKey[pt] {
		t = &pemresources.PublicKey{PemResource: pemresources.PemResource{Location: location}}
	}
	if t == nil {
		return nil, fmt.Errorf("%s is an unsupported pem type", blk.Type)
	}

	// If template, parse as a yaml;
	var err error
	if isTemplate {
		err = yaml.Unmarshal(blk.Bytes, t)
	} else {
		// otherwise parse as a pem
		err = t.UnmarshalPem(blk)
	}
	return t, err
}
