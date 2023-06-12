package commands

import (
	"fmt"
	"io"
	"strings"
)

const (
	CommandMake     = "make"
	CommandFind     = "find"
	CommandTemplate = "template"
	CommandConfig   = "config"
	CommandKeys     = "keys"
)

// Commands maps the command name to the Command instance
var Commands = map[string]Command{
	CommandMake:     &MakeCommand{},
	CommandFind:     &FindCommand{},
	CommandTemplate: &TemplateCommand{},
	CommandConfig:   &ConfigCommand{},
	CommandKeys:     &keysCommand{},
}

// CommandAliases maps alternative names for commands, to the actual command name
var CommandAliases = map[string]string{
	"mk":        CommandMake,
	"fd":        CommandFind,
	"f":         CommandFind,
	"tp":        CommandTemplate,
	"temp":      CommandTemplate,
	"templates": CommandTemplate,
	"cf":        CommandConfig,
	"cfg":       CommandConfig,
	"k":         CommandKeys,
	"key":       CommandKeys,
}

// Command executes a single operation using the given arguments and any flags assigned to the Commands public fields.
type Command interface {
	Execute(args []string, out io.Writer) error
}

// CustomFlagsCommand is a Command which processes its own arbitrary flags.
type CustomFlagsCommand interface {
	Command
	ApplyFlags(flags map[string]*string) error
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
