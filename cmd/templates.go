package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"pempal/templates"
)

// TemplatesCommand lists all the available templates
type TemplatesCommand struct {
	showTemplatepath     bool
	showBuildInTemplates bool
}

func (cmd *TemplatesCommand) Description() string {
	return fmt.Sprintf("returns a list of all the available templates found in the current directory and $%s", templates.ENV_TemplatePath)
}

func (cmd *TemplatesCommand) Flags(f *flag.FlagSet) {
	f.BoolVar(&cmd.showTemplatepath, "kp", false, "display the current keypath")
	f.BoolVar(&cmd.showBuildInTemplates, "b", false, "include all the built in template names at the end of the list")
}

func (cmd TemplatesCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	if cmd.showTemplatepath {
		_, err := fmt.Fprintf(out, "$%s: %s\n", templates.ENV_TemplatePath, templates.TemplatePath)
		if err != nil {
			return err
		}
	}
	for _, p := range templates.TemplateNames(cmd.showBuildInTemplates) {
		_, err := fmt.Fprintln(out, p)
		// TODO: Read each template for a "Description" tag
		if err != nil {
			return err
		}
	}
	return nil
}
