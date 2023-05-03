package commands

import (
	"fmt"
	"io"
	"strings"
)

var commands = map[string]Command{
	"find": &FindCommand{},
}

// CommonFlags are flags which all command can use without the need to declare them in the command class.
var CommonFlags CommonFlagArgument

// Command executes a single operation using the given arguments and any flags assigned to the commands public fields.
type Command interface {
	Execute(args []string, out io.Writer) error
}

// CommandWithFlags is a custom Command which processes arbitary flags.
type CommandWithFlags interface {
	Command
	SetFlags(flags map[string]*string) error
}

type CommonFlagArgument struct {
	Out      string `flag:"out"`
	ForceOut bool   `flag:"force,f"`
	Verbose  bool   `flag:"v"`
}

func NewCommand(name string) (Command, error) {
	cmd, ok := commands[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("%s is an unknown command", name)
	}
	return cmd, nil
}
