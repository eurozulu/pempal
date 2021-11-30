package commands

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"pempal/templates"
)

// TemplatesCommand lists all the available templates
type TemplatesCommand struct {
	ShowPath         bool
	ShowAllTemplates bool
}

func (cmd *TemplatesCommand) Description() string {
	lines := bytes.NewBufferString("returns a list of all the available templates found in the current directory and $")
	lines.WriteString(templates.ENV_TemplatePath)
	lines.WriteString("\nTemplates are yaml files with a '#' as the first character of the name.\n")
	lines.WriteString("They contain the properties to be applied to new resources when they're being generated.\n")
	lines.WriteString("Users should place their templates in the current directory or in any directory listed in the comma delimited $")
	lines.WriteString(templates.ENV_TemplatePath)
	lines.WriteString(" environment variable\n")
	lines.WriteRune('\n')
	lines.WriteString("There are a number of built in templates to help create common resource types.\n")
	lines.WriteString("By default, templates command will only show user defined templates.  To show all the built in temapltes available, use the -b flag\n")
	return lines.String()
}

func (cmd *TemplatesCommand) Flags(f *flag.FlagSet) {
	f.BoolVar(&cmd.ShowPath, "tp", false, fmt.Sprintf("display the current $%s", templates.ENV_TemplatePath))
	f.BoolVar(&cmd.ShowAllTemplates, "b", false, "include all the built in template names at the end of the list")
}

func (cmd TemplatesCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	if cmd.ShowPath {
		tp := templates.TemplatePath
		if tp == "" {
			tp = "not set"
		}
		_, err := fmt.Fprintf(out, "$%s: %s\n", templates.ENV_TemplatePath, tp)
		if err != nil {
			return err
		}
	}
	for _, p := range templates.TemplateNames(cmd.ShowAllTemplates) {
		_, err := fmt.Fprintln(out, p)
		// TODO: Read each template for a "Description" tag
		if err != nil {
			return err
		}
	}
	return nil
}
