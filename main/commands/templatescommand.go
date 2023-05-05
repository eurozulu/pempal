package commands

import (
	"fmt"
	"io"
	"sort"
)

type TemplatesCommand struct {
	Name string `flag:"name"`
}

func (t TemplatesCommand) Execute(args []string, out io.Writer) error {
	names := ResourceTemplates.Names(args...)
	sort.Strings(names)
	for _, n := range names {
		fmt.Fprintln(out, n)
	}
	return nil
}
