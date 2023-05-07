package commands

import (
	"fmt"
	"io"
	"pempal/config"
	"pempal/main/argparser"
	"pempal/resourceio"
	"pempal/templates"
	"strings"
)

var commands = map[string]Command{
	"find":      &FindCommand{},
	"make":      &MakeCommand{},
	"config":    &ConfigCommand{},
	"templates": &TemplatesCommand{},
	"template":  &TemplateCommand{},
}

// CommonFlags are flags which all command can use without the need to declare them in the command class.
var CommonFlags CommonFlagsStruct

// Configuration contains the shared config for all commands
var Configuration config.Config

// ResourceTemplates is the shared TemplateManager
var ResourceTemplates templates.TemplateManager

// Command executes a single operation using the given arguments and any flags assigned to the commands public fields.
type Command interface {
	Execute(args []string, out io.Writer) error
}

// CommandWithFlags is a custom Command which processes arbitrary flags.
type CommandWithFlags interface {
	Command
	SetFlags(flags map[string]*string) error
}

// CommonFlagsStruct contains the flags used by all commands
type CommonFlagsStruct struct {
	Out        string `flag:"out"`
	ForceOut   bool   `flag:"force,f"`
	ConfigPath string `flag:"config"`
	Verbose    bool   `flag:"v"`
	Debug      bool   `flag:"vv"`
}

func ApplyCommonFlags(args []string) ([]string, error) {
	args, err := argparser.ApplyArguments(args, &CommonFlags)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse common flags  %v", err)
	}
	Configuration = config.NewConfig(CommonFlags.ConfigPath)
	ResourceTemplates, err = resourceio.NewResourceTemplateManager(Configuration.TemplatePath)
	if err != nil {
		return nil, err
	}
	return args, nil
}

func NewCommand(name string) (Command, error) {
	cmd, ok := commands[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("%s is an unknown command", name)
	}
	return cmd, nil
}
