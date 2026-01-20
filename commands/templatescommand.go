package commands

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/repositories"
	"github.com/eurozulu/pempal/templates"
	"io"
	"strings"
)

var templateRepo = repositories.Templates(config.TemplatePath())

// TemplatesCommand manages the templates, listing, displaying and merging them.
// @Command(templates, "template,temps,temp")
type TemplatesCommand struct {

	// Merge when set will merge multiple templates into one.
	// Only applies when outputting templates. i.e. ignored is listing or expanding names
	// @Flag(merge, m)
	Merge bool

	// Extends when set, displays the given template name(s) preceeded with any template names they extend.
	// @Flag(extends, e)
	Extends bool
}

// ListNames shows all the template names, beginning with the base templates.
// @Action(list)
func (c *TemplatesCommand) ListNames() string {
	var lines []string
	for _, n := range templateRepo.TemplateNames() {
		lines = append(lines, n)
	}
	return strings.Join(lines, "\n")
}

// ShowTemplate displays the templates of the given names.
// Each name should be a known template.  When given with no flags,
// each template is output, onbe after another.
// When -merge is used, the single merged template is output.
// If no names are given, all of the template names are listed.
// @Action
func (c *TemplatesCommand) ShowTemplate(name ...string) (string, error) {
	argFlags, names, err := ArgFlagsToTemplate(name)
	if err != nil {
		return "", err
	}
	if len(names) == 0 {
		return "", fmt.Errorf("specifiy one or more template names:\n%s", c.ListNames())
	}

	if c.Extends {
		s, err := c.ExpandNames(names...)
		if err != nil {
			return "", err
		}
		return strings.Join(s, " "), nil
	}
	if c.Merge {
		names, err = templateRepo.ExpandNames(names...)
		if err != nil {
			return "", err
		}
	}
	temps, err := templateRepo.ByName(names...)
	if err != nil {
		return "", err
	}

	if argFlags.String() != "" {
		temps = append(temps, argFlags)
		names = append(names, "")
	}

	if c.Merge {
		t, err := templates.MergeTemplates(temps)
		if err != nil {
			return "", err
		}
		temps = []templates.Template{t}
		names = []string{strings.Join(insertArrowsToNames(names), " ")}
	}

	buf := bytes.NewBuffer(nil)
	for i, temp := range temps {
		if i > 0 {
			buf.WriteString("\n----\n")
		}
		c.writeTemplate(names[i], temp, buf)
	}
	return buf.String(), nil
}

// ExpandNames displays the full chain of templates which are applied for the given names.
// If any of the given names point to a template with an 'extends' clause, the names of those extends
// are preceeded the name.  Aliases for base names are resolved to full names.
func (c *TemplatesCommand) ExpandNames(name ...string) ([]string, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("specifiy one or more template names:\n%s", c.ListNames())
	}
	names, err := templateRepo.ExpandNames(name...)
	if err != nil {
		return nil, err
	}
	return insertArrowsToNames(names), nil
}

func insertArrowsToNames(names []string) []string {
	for i := 0; i < len(names)-1; i++ {
		names[i] = names[i] + " <-"
	}
	return names
}

func (c *TemplatesCommand) writeTemplate(name string, t templates.Template, buf io.StringWriter) {
	if Verbose {
		buf.WriteString("name: ")
		buf.WriteString(name)
		buf.WriteString("\n")

		if tf, ok := t.(*model.TemplateFile); ok {
			buf.WriteString("path: ")
			buf.WriteString(tf.Path)
			buf.WriteString("\n")
		}
		if tt := templates.TypeOfBaseTemplate(t); tt != "" {
			buf.WriteString("type: ")
			buf.WriteString(tt)
			buf.WriteString("\n")
		}
	}
	s := t.String()
	buf.WriteString(s)
	if !strings.HasSuffix(s, "\n") {
		buf.WriteString("\n")
	}
}
