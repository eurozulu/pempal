package commands

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"pempal/templates"
)

var templateManager = templates.NewTemplateManager()

type templateCommand struct {
	ShowNames bool `flag:"list"`
}

func (m templateCommand) Main(args Arguments, out io.Writer) error {
	if m.ShowNames {
		l := len(args.Parameters())
		if l > 0 {
			fmt.Printf("ignoring %d parameters as -list flag set\n", l)
		}
		return writeNames(out)
	}

	temps, err := templateManager.FindTemplates(args.Parameters()...)
	if err != nil {
		return err
	}
	if len(temps) == 0 {
		return fmt.Errorf("No template found.  Name at least one valid template name")
	}
	root := temps[0]
	if len(temps) > 1 {
		templateManager.MergeTemplates(root, temps[1:]...)
	}
	return yaml.NewEncoder(out).Encode(root)
}

func writeNames(out io.Writer) error {
	names, err := templateManager.TemplateNames()
	if err != nil {
		return err
	}
	buf := bufio.NewWriter(out)
	for _, n := range names {
		buf.WriteString(n)
		buf.WriteByte('\n')
	}
	return buf.Flush()
}
