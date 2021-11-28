package cmd

import (
	"context"
	"flag"
	"gopkg.in/yaml.v3"
	"io"
	"pempal/templates"
)

// TemplateCommand displays and builds templates
type TemplateCommand struct {
	BuildTemplate bool
}

func (cmd TemplateCommand) Description() string {
	return "displays listed resources or templates and, optionally, merges them into a single template"
}

func (cmd TemplateCommand) Flags(f *flag.FlagSet) {
	f.BoolVar(&cmd.BuildTemplate, "b", false, "builds a single template, merging all the given locations and templates into one.")
}

func (cmd TemplateCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	tb := templates.NewTemplateBuilder()
	if err := tb.Add(args...); err != nil {
		return err
	}

	var temps []templates.Template
	if cmd.BuildTemplate {
		t, err := tb.Build()
		if err != nil {
			return err
		}
		temps = append(temps, t)
	} else {
		temps = tb.Templates()
	}
	yout := yaml.NewEncoder(out)
	for _, t := range temps {
		if err := yout.Encode(&t); err != nil {
			return err
		}
		_, err := out.Write([]byte("\n---\n"))
		if err != nil {
			return err
		}
	}
	return nil
}
