package main

import (
	"github.com/pempal/pemio"
	"github.com/pempal/templates"
)

type undefinedTemplate map[string]interface{}

type TemplateCommand struct {
}

func (tc TemplateCommand) Template(templates ...string) {

}

// NewTemplateFiles loads all the PEM blocks found in the slice of file paths.
func NewTemplateFiles(filepaths []string) ([]templates.Template, error) {
	bls, err := pemio.ReadPEMsFiles(filepaths)
	if err != nil {
		return nil, err
	}
	ts, err := templates.NewTemplates(bls)
	if err != nil {
		return nil, err
	}
	return ts, nil
}
