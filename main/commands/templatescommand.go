package commands

import (
	"fmt"
	"io"
	"pempal/templates"
	"sort"
)

type TemplatesCommand struct {
	Name   string `flag:"name"`
	Add    string `flag:"add"`
	Remove string `flag:"remove"`
}

func (cmd TemplatesCommand) Execute(args []string, out io.Writer) error {
	if cmd.Remove != "" {
		return cmd.removeTemplates(args)
	}

	store := ResourceTemplates.(templates.TemplateStore)
	names := store.Names(args...)
	sort.Strings(names)
	for _, n := range names {
		fmt.Fprintln(out, n)
	}
	return nil
}

func (cmd TemplatesCommand) addTemplate(name string, data []byte) error {
	t, err := ResourceTemplates.ParseTemplate(data)
	if err != nil {
		return err
	}
	return ResourceTemplates.(templates.TemplateStore).SaveTemplate(name, t)
}

func (cmd TemplatesCommand) removeTemplates(names []string) error {
	for _, name := range names {
		if err := ResourceTemplates.(templates.TemplateStore).DeleteTemplate(name); err != nil {
			return err
		}
	}
	return nil
}
