package commands

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/pempal/builder"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/templates"
	"io"
	"sort"
	"strings"
)

// TemplateCommand, when used with no parameters, lists the names of all the available templates.
type TemplateCommand struct {
	Apply         bool `flag:"apply"`
	All           bool `flag:"all,a"`
	templateStore templates.TemplateStore
}

func (cmd TemplateCommand) Execute(args []string, out io.Writer) error {
	// setup the template store
	if ts, err := config.TemplateStore(); err == nil {
		cmd.templateStore = ts
	} else {
		return fmt.Errorf("template store unavailable, %v", err)
	}

	// parse argument into simple names and assignments (names followed by '=')
	assignments, names := parseArgsForAssignments(args)
	if err := cmd.addNewTemplates(assignments); err != nil {
		return err
	}

	if len(args) == 0 {
		// if no name given, list all names
		return cmd.writeTemplateNames(out)
	}
	if cmd.Apply {
		return cmd.applyTemplates(out, names)
	}
	return cmd.writeRawTemplates(out, names)
}

func (cmd TemplateCommand) applyTemplates(out io.Writer, names []string) error {
	temps, err := cmd.templateStore.ExtendedTemplatesByName(names...)
	if err != nil {
		return err
	}

	t, err := builder.MergeTemplates(temps)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, string(t.Bytes()))
	return err
}

func (cmd TemplateCommand) writeRawTemplates(out io.Writer, names []string) error {
	buf := bytes.NewBuffer(nil)
	for i, n := range names {
		if i > 0 {
			buf.WriteString("\n---\n")
		}
		t, err := cmd.templateStore.TemplateByName(n)
		if err != nil {
			return err
		}
		if !CommonFlags.Quiet {
			buf.WriteString("Template ")
			buf.WriteString(n)
			buf.WriteString(":-\n")
		}
		buf.WriteString(t.String())
	}
	_, err := out.Write(buf.Bytes())
	return err
}

func (cmd TemplateCommand) writeTemplateNames(out io.Writer) error {
	var names []string
	if cmd.All {
		names = cmd.templateStore.AllNames()
	} else {
		names = cmd.templateStore.Names()
	}
	if !CommonFlags.Quiet && len(names) == 0 {
		logger.Info("No templates found")
		return nil
	}
	sort.Strings(names)
	for _, n := range names {
		if _, err := fmt.Fprintln(out, n); err != nil {
			return err
		}
	}
	if !CommonFlags.Quiet {
		logger.Info("%d templates found", len(names))
	}
	return nil
}

func (cmd TemplateCommand) addNewTemplates(names []string) error {
	for _, name := range names {
		ss := strings.SplitN(name, "=", 2)
		name = ss[0]
		var temp templates.Template
		if len(ss) > 1 {
			t, err := templates.ParseInlineTemplate(ss[1])
			if err != nil {
				return err
			}
			temp = t
		}
		if CommonFlags.ForceOut && cmd.templateStore.Exists(name) {
			if err := cmd.templateStore.DeleteTemplate(name); err != nil {
				return err
			}
		}
		if err := cmd.templateStore.SaveTemplate(name, temp); err != nil {
			return fmt.Errorf("Failed to save template '%s'  %v", name, err)
		}
		if !CommonFlags.Quiet {
			logger.Info("created template %s\n", name)
		}
	}
	return nil
}
