package commands

import (
	"fmt"
	"github.com/eurozulu/pempal/resourceio"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/utils"
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
	tm, err := resources.NewResourceTemplatesManager(ResolvePath(CommonFlags.TemplatePath))
	if err != nil {
		return nil, err
	}

	var tps []templates.Template
	for _, arg := range args {
		// Check each arugment is valid filepath or template name
		var ts []templates.Template
		var err error
		if utils.FileExists(arg) {
			ts, err = resourceio.LoadTemplatesFromFile(arg)
		} else {
			// not a file, check if it's a template
			if t, er := tm.TemplateByName(arg); er != nil {
				err = er
			} else {
				ts = []templates.Template{t}
			}

		}
		if err != nil {
			return nil, err
		}
		tps = append(tps, ts...)
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
