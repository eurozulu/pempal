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
// When one template name is given, that template is shown in its formatted state. (After any marco and extends have been proccessed)
// When more than one name given the right most template is said to extend the template name it follows.
// The extended templates (right) is 'applied' to a copy of the the template being extended (left), adding any keys not already present and
// overwritting keys which are present.
// Multiple templates may be 'chained' in extenntions, each adding or changing the final values within the template.
// The final result is a single template containing all of the keys from all of the templates.
// The extended of the named templates can be turned off using the '-list' flag.  When set, each named template is shown in its formatted form.
// Templates are shown in the formatted format, after any tags have been processed.  To show the template in its unformetted state,
// i.e. how it is stored, prior to tag processing, use the -raw flag.
type TemplateCommand struct {
	Unformatted bool `flag:"unformatted,noformat,raw"`
	NoMerge     bool `flag:"no-merge,nomerge,list"`

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
	if ns, err := cmd.addNewTemplates(assignments); err != nil {
		return err
	} else if !CommonFlags.Quiet {
		names = append(names, ns...)
	}

	if cmd.Unformatted {
		return cmd.writeRawTemplates(out, names)
	}
	if len(args) == 0 {
		// if no name given, list all names
		return cmd.writeTemplateNames(out)
	}

	return cmd.writeTemplates(out, names)
}

func (cmd TemplateCommand) writeTemplates(out io.Writer, names []string) error {
	temps, err := cmd.templateStore.ExtendedTemplatesByName(names...)
	if err != nil {
		return err
	}
	if !cmd.NoMerge {
		t, err := builder.MergeTemplates(temps)
		if err != nil {
			return err
		}
		temps = []templates.Template{t}
	}
	buf := bytes.NewBuffer(nil)
	for i, t := range temps {
		if i > 0 {
			buf.WriteString("\n---\n")
		}
		if _, err := buf.Write(t.Bytes()); err != nil {
			return err
		}
	}
	_, err = out.Write(buf.Bytes())
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
		buf.WriteString(t.String())
	}
	_, err := out.Write(buf.Bytes())
	return err
}

func (cmd TemplateCommand) writeTemplateNames(out io.Writer) error {
	names := cmd.templateStore.Names()
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
	return nil
}

func (cmd TemplateCommand) addNewTemplates(names []string) ([]string, error) {
	for i, name := range names {
		ss := strings.SplitN(name, "=", 2)
		name = ss[0]
		var temp templates.Template
		if len(ss) > 1 {
			t, err := templates.ParseInlineTemplate(ss[1])
			if err != nil {
				return nil, err
			}
			temp = t
		}
		if CommonFlags.ForceOut && cmd.templateStore.Exists(name) {
			if err := cmd.templateStore.DeleteTemplate(name); err != nil {
				return nil, err
			}
		}
		if err := cmd.templateStore.SaveTemplate(name, temp); err != nil {
			return nil, fmt.Errorf("Failed to save template '%s'  %v", name, err)
		}
		if !CommonFlags.Quiet {
			logger.Info("created template %s\n", name)
		}
		names[i] = name
	}
	return names, nil
}
