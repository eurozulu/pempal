package finder

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"pempal/templates"
	"strings"
)

var templateFileExtensions = map[string]bool{"": true, "yaml": true, "template": true}

type templateParser struct {
	temps []templates.Template
}

func (pp templateParser) Parse(path string, data []byte) (Location, error) {
	t := templates.NewTemplate(0)
	if err := yaml.NewDecoder(bytes.NewBuffer(data)).Decode(&t); err != nil {
		return nil, fmt.Errorf("failed to read template '%s'  %v", path, err)
	}
	return &templateLocation{
		path:  path,
		temps: []templates.Template{t},
	}, nil
}

func (pp templateParser) MatchPath(path string) bool {
	return templateFileExtensions[strings.ToLower(strings.TrimLeft(filepath.Ext(path), "."))]
}

func (pp templateParser) FilterLocation(rl Location) Location {
	tl, ok := rl.(*templateLocation)
	if !ok {
		return nil
	}
	return tl
}

func newTemplateParser() *templateParser {
	return &templateParser{}
}
