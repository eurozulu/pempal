package commands

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/pempal/commandline/commonflags"
	"github.com/eurozulu/pempal/resourceio"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
	"io"
	"strings"
)

const (
	CommandMake = "make"
	CommandKeys = "keys"

	CommandTemplate  = "template"
	CommandTemplates = "templates"

	CommandType  = "type"
	CommandFind  = "find"
	CommandPaths = "paths"
)

// Commands maps the command name to the Command instance
var Commands = map[string]Command{
	CommandMake: &MakeCommand{},
	CommandKeys: &keysCommand{},
	CommandType: &TypeCommand{},
	//CommandFind:     &FindCommand{},
	//CommandTemplate: &TemplateCommand{},
	//CommandConfig:   &ConfigCommand{},
}

// CommandAliases maps alternative names for commands, to the actual command name
var CommandAliases = map[string]string{
	"mk": CommandMake,
	"ks": CommandKeys,

	"f":         CommandFind,
	"t":         CommandTemplate,
	"ts":        CommandTemplates,
	"templates": CommandTemplate,
	"cfg":       CommandPaths,
}

// Command executes a single operation using the given arguments and any flags assigned to the Commands public fields.
type Command interface {
	Execute(args []string, out io.Writer) error
}

func argumentsToTemplates(args []string) ([]templates.Template, error) {
	tm, err := resources.NewResourceTemplatesManager(commonflags.ResolvePath(commonflags.CommonFlags.TemplatePath))
	if err != nil {
		return nil, err
	}

	var tps []templates.Template
	for _, arg := range args {
		// Check each arugment is a known template name or file location
		var errs []error
		// first check if it's a template
		if t, err := tm.TemplateByName(arg); err == nil {
			tps = append(tps, t)
			continue
		} else {
			errs = append(errs, err)
		}

		// Not a template, check if known to filepaths
		if loc, err := commonflags.CommonFlags.FindInPath(arg, false); err == nil {
			ts, err := resourceio.ResourceLocationToTemplates(loc)
			if err != nil {
				return nil, err
			}
			tps = append(tps, ts...)
			continue
		} else {
			errs = append(errs, err)
		}

		// neither template or file location valid, report errors
		buf := bytes.NewBuffer(nil)
		for _, err := range errs {
			if buf.Len() > 0 {
				buf.WriteString(" or ")
			}
			buf.WriteString(err.Error())
		}
		return nil, fmt.Errorf("%s", buf.String())
	}
	return tps, nil
}

func NewCommand(name string) (Command, error) {
	n := strings.ToLower(name)
	alias, ok := CommandAliases[n]
	if ok {
		n = alias
	}
	cmd, ok := Commands[n]
	if !ok {
		return nil, fmt.Errorf("%s is an unknown command", name)
	}
	return cmd, nil
}
