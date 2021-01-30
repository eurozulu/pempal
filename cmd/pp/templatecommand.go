package main

import (
	"github.com/eurozulu/pempal/templates"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type TemplateCommand struct {
	TemplatePath []string

	Insensitive bool
}

// Template will output the named template.
// Searches the path set on $PP_TEMPLATES for files with the given name (without extension)
// Also includes the three inbuilt templates: server, ca and user.
func (tc TemplateCommand) Template(name string) error {

	t := templates.DefaultTemplate(name)
	if t != nil {
		return yaml.NewEncoder(os.Stdout).Encode(t)
	}
	if len(tc.TemplatePath) == 0 {
		tc.TemplatePath = strings.Split(os.Getenv("PP_TEMPLATES"), ",")
	}

	return nil
}


