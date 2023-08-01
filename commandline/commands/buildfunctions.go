package commands

import (
	"fmt"
	"github.com/eurozulu/pempal/builders"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
	"io"
)

func PerformBuild(rt resources.ResourceType, t templates.Template, out io.Writer) error {
	builder, err := builders.NewBuilder(rt)
	if err != nil {
		return err
	}
	res, err := builder.Build(t)
	if err != nil {
		return err
	}

	// flip new keyresource into a PEM string
	dto, err := resources.NewResourceDTO(res)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, dto.String())
	return err
}
