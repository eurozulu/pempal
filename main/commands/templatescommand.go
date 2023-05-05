package commands

import (
	"fmt"
	"io"
	"pempal/resourceio"
	"sort"
)

type TemplatesCommand struct {
	Name         string `flag:"name"`
	TemplatePath string `flag:"template-path,path"`
}

func (t TemplatesCommand) Execute(args []string, out io.Writer) error {
	if t.TemplatePath == "" {
		t.TemplatePath = configuration.TemplatePath
	}

	tm, err := resourceio.NewResourceTemplateManager(t.TemplatePath)
	if err != nil {
		return err
	}
	names := tm.Names(args...)
	sort.Strings(names)
	for _, n := range names {
		fmt.Fprintln(out, n)
	}
	return err
}
