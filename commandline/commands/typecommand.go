package commands

import (
	"fmt"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
	"io"
	"strings"
)

// TypeCommand displays the resource type of a given template
type TypeCommand struct {
}

func (t TypeCommand) Execute(args []string, out io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("Must provide a template name or filepath to a template")
	}
	temps, err := argumentsToTemplates(args)
	if err != nil {
		return err
	}
	// Merge into single template
	mt := templates.NewTemplateBuilder(temps...).Build()
	rts, err := resources.TemplateTypes(mt)
	if err != nil {
		return err
	}
	names := strings.Join(args, ", ")
	if len(rts) == 1 {
		fmt.Fprintf(out, "template '%s' is of type %s", names, rts[0].String())
	} else {
		fmt.Fprintf(out, "template '%s' supports the following types:\n", names)
		for _, rt := range rts {
			fmt.Fprintf(out, "    %s\n", rt.String())
		}
	}
	return nil
}
