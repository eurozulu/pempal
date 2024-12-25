package commands

import (
	"fmt"
	"github.com/eurozulu/pempal/factories"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/validation"
	"io"
)

type MakeCommand struct {
	Output io.Writer
}

func (mc MakeCommand) Exec(args ...string) error {
	tc := &TemplateCommand{}
	template, err := tc.BuildTemplate(args...)
	if err != nil {
		return err
	}
	if err := validation.Validate(template); err != nil {
		return err
	}
	data, err := mc.Make(template)
	if err != nil {
		return err
	}
	if mc.Output != nil {
		// write to given output
		_, err := mc.Output.Write(data)
		return err
	}
	// Output is not set, write to default location for resource
	path, err := factories.SaveResource(data)
	if err != nil {
		return err
	}
	fmt.Printf("new %s created at %q", template.Name(), path)
	return nil
}

func (mc MakeCommand) Make(template templates.Template) ([]byte, error) {
	fac, err := factories.FactoryForType(template.Name())
	if err != nil {
		return nil, err
	}
	return fac.Build(template)
}
