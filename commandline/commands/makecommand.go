package commands

import (
	"fmt"
	"github.com/eurozulu/argdecoder"
	"github.com/eurozulu/pempal/builder"
	"github.com/eurozulu/pempal/commandline/formselect"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/resourceio"
	"github.com/eurozulu/pempal/templates"
	"io"
	"strings"
)

const defaultResourceFormat = resourceio.FormatPEM

type MakeCommand struct {
	Format string `flag:"format,fm"`
}

func (cmd MakeCommand) Execute(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("must provide one or more template names to build")
	}

	km, err := config.KeyManager()
	if err != nil {
		return err
	}

	// Collect the templates to build
	tm, err := config.TemplateManager()
	if err != nil {
		return err
	}

	// collect the named templates from the arguments
	names, flags := argdecoder.ParseArgs(args)
	temps, err := tm.ExtendedTemplatesByName(names...)
	if err != nil {
		return err
	}
	// Build a 'flag' template and add to the list
	if len(flags) > 0 {
		t, err := flagTemplate(flags)
		if err != nil {
			return err
		}
		temps = append(temps, t)
	}

	// Create builder with templates.  Should be able to establish build type from these.
	tb := builder.TemplateBuilder(temps)
	// Build the final template, based on the established type
	t, err := tb.MergeTemplates()

	build, err := builder.NewBuilder(tb.ResourceType(), km)
	if err != nil {
		return err
	}

	errs := build.Validate(t)
	if CommonFlags.Quiet && len(errs) > 0 {
		return builder.CombineErrors(errs)
	}
	// request new values to correct errors
	for len(errs) > 0 {
		ct, err := correctionTemplate(t)
		if err != nil {
			return err
		}
		t, err = builder.TemplateBuilder([]templates.Template{t, ct}).MergeTemplates()
		if err != nil {
			return err
		}
		errs = build.Validate(t)
	}

	r, err := build.Build(t)
	if err != nil {
		return err
	}

	if err = cmd.writeResource(out, r); err != nil {
		return err
	}
	return nil
}

func (cmd MakeCommand) writeResource(out io.Writer, r ...model.Resource) error {
	formOut, err := cmd.getResourceFormatter()
	if err != nil {
		return err
	}
	data, err := formOut.FormatResources(r...)
	if err != nil {
		return err
	}
	_, err = out.Write(data)
	return err
}

func (cmd MakeCommand) getResourceFormatter() (resourceio.ResourceFormatter, error) {
	var err error
	rf := defaultResourceFormat
	if cmd.Format != "" {
		rf, err = resourceio.ParseResourceFormat(cmd.Format)
		return nil, err
	}
	return resourceio.NewResourceFormatter(rf), nil
}

func parseErrorNames(errs []error) []string {
	names := make([]string, len(errs))
	for i, err := range errs {
		names[i] = strings.SplitN(err.Error(), " ", 2)[0]
	}
	return names
}

func correctionTemplate(t templates.Template) (templates.Template, error) {

	f := formselect.NewForm(5, 10, lines)
	f.SelectLine(-1)
	return io.EOF
}
