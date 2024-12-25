package commands

import (
	"fmt"
	"github.com/eurozulu/pempal/validation"
	"io"
	"os"
)

type ValidateCommand struct {
	Output      io.Writer
	Interactive bool
}

func (vc *ValidateCommand) Exec(args ...string) error {
	tc := &TemplateCommand{}
	ft, err := tc.BuildTemplate(args...)
	if err != nil {
		return err
	}

	if err = validation.Validate(ft); err != nil {
		vc.writeOut("validating %s template failed\n", ft.Name())
		return err
	}
	vc.writeOut("%s template is valid\n", ft.Name())
	return nil
}

func (vc *ValidateCommand) writeOut(format string, a ...any) error {
	o := vc.Output
	if o == nil {
		o = os.Stdout
	}
	_, err := fmt.Fprintf(o, format, a...)
	return err
}
