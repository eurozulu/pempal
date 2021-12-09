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
	names := templateFileNames()
	// add built in names
	if includeBuiltIn {
		names = append(names, templateInBuiltNames()...)
	}
	// TODO: // search for duplicate names as directory is not considured. When found, add the parset path to it. e.g. mytemps/#thattemp
	for i, n := range names {
		names[i] = cleanName(n)
	}
	return names
}

func FindTemplate(name string) (Template, error) {
	key := strings.ToLower(strings.TrimLeft(name, FileTag))
	if t, ok := builtInTemplates[key]; ok {
		return ParseTemplate(t.([]byte))
	}
	// not built in, search for file
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	ts := TemplateScanner{
		Recursive: true,
		Verbose:   false,
		Reader:    fileformats.NewFormatReader(),
		FilterFunc: func(p string) bool {
			// get name without extension
			n := path.Base(p)
			e := len(path.Ext(n))
			n = n[:len(n)-e]
			return strings.EqualFold(n, name)
		},
	}
	tch := ts.Find(ctx, templatePath()...)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case t, ok := <-tch:
		if !ok {
			return nil, os.ErrNotExist
		}
		return t, nil
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
func templateFileNames() []string {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	ts := filepathscanner.FilePathScanner{
		Recursive: true,
		ExtFilter: map[string]bool{"yaml": true},
	}
	nameCh := ts.Scan(ctx, templatePath()...)
	var names []string
	for name := range nameCh {
		if !strings.HasPrefix(name, FileTag) {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func cleanName(p string) string {
	n := path.Base(p)
	// strip the file extension
	e := path.Ext(n)
	if len(e) > 0 {
		n = n[:len(n)-len(e)]
	}
	if !strings.HasPrefix(n, FileTag) {
		n = strings.Join([]string{FileTag, n}, "")
	}
	return n
}

// templatePath gets the template path, the path(s) to search for private keys
// Template path is the current working directory with any additional (comma demited) paths set in the ENV_TemplatePath
func templatePath() []string {
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
	return TemplateFromBlock(tblk[0], "")
}

func TemplateFromBlock(blk *pem.Block, location string) (Template, error) {
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
