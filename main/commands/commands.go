package commands

import (
	"fmt"
	"io"
	"pempal/config"
	"pempal/keymanager"
	"pempal/logger"
	"pempal/main/argdecoder"
	"pempal/resourceio"
	"pempal/templates"
	"strings"
)

var Commands = map[string]Command{
	"find":      &FindCommand{},
	"make":      &MakeCommand{},
	"config":    &ConfigCommand{},
	"templates": &TemplatesCommand{},
	"template":  &TemplateCommand{},
	"keys":      &KeysCommand{},
}

// CommonFlags are flags which all command can use without the need to declare them in the command class.
var CommonFlags CommonFlagsStruct

// Configuration contains the shared config for all Commands
var Configuration config.Config

// ResourceTemplates is the shared TemplateManager
var ResourceTemplates templates.TemplateManager

var Keys keymanager.KeyManager

// Command executes a single operation using the given arguments and any flags assigned to the Commands public fields.
type Command interface {
	Execute(args []string, out io.Writer) error
}

// CommandWithFlags is a custom Command which processes arbitrary flags.
type CommandWithFlags interface {
	Command
	SetFlags(flags map[string]*string) error
}

// CommonFlagsStruct contains the flags used by all Commands
type CommonFlagsStruct struct {
	Out        string `flag:"out"`
	ForceOut   bool   `flag:"force,f"`
	ConfigPath string `flag:"config"`
	Verbose    bool   `flag:"v"`
	Debug      bool   `flag:"vv"`
	Help       bool   `flag:"help"`
}

func ApplyCommonFlags(args []string) ([]string, error) {
	args, err := argdecoder.ApplyArguments(args, &CommonFlags)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse common flags  %v", err)
	}
	Configuration = config.NewConfig(CommonFlags.ConfigPath)
	ResourceTemplates, err = resourceio.NewResourceTemplateManager(Configuration.TemplatePath)
	if err != nil {
		logger.Log(logger.Error, "Failed to load template manager. %v", err)
	}

	Keys, err = keymanager.NewKeyManager(Configuration.KeyPath, Configuration.CertPath)
	if err != nil {
		logger.Log(logger.Error, "Failed to load key manager. %v", err)
	}
	return args, nil
}

func NewCommand(name string) (Command, error) {
	cmd, ok := Commands[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("%s is an unknown command", name)
	}
	return cmd, nil
}
