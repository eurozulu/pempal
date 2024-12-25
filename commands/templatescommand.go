package commands

import (
	"fmt"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/templates"
	"io"
	"os"
)

type TemplatesCommand struct {
	Output io.Writer
}

func (t TemplatesCommand) Exec(args ...string) error {
	tlib := templates.NewTemplateLib(config.Config.Templates...)
	names := tlib.GetTemplateNames()

	if t.Output == nil {
		t.Output = os.Stdout
	}
	for _, name := range names {
		fmt.Fprintln(t.Output, name)
	}
	return nil
}
