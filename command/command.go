package command

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
)

// commandsMap is the available commands, mapped to the string command(s)
var commandsMap = map[string]Command{
	"find": &findCommand{},
	"f":    &findCommand{},
	//"template":  &templateCommand{},
	//"t":         &templateCommand{},
	//"templates": &templatesCommand{},
}

// Command represents a single operation
type Command interface {
	Run(ctx context.Context, args Arguments, out io.Writer) error
}

// ApplyArguments applies any values in the given arguments whic are flags, and match fields in the command struct tagged with the same name.
// Other than bool flags, all other fields require a value following the flag argument.  This value must be able to be parsed into the field type it matched.
// If the field is a Bool value, the Flag is not expected to contain a value. Its presents indacates true.
// all arguments matched to fields are 'consumed' and the remaining arguments returned.
func ApplyArguments(cmd interface{}, args []string) (Arguments, error) {
	flagTags, err := newFlagFields(cmd)
	if err != nil {
		// log fatal is should never pass non struct to this]
		log.Fatalln(err)
	}
	a := &arguments{args: args}
	for _, fl := range flagTags.Names() {
		if !a.ContainsFlag(fl) {
			continue
		}
		var v string
		if flagTags.IsBool(fl) {
			// If bool flag no value is expected
			a.removeFlag(fl)
		} else {
			// non bool must have value
			if !a.HasValue(fl) {
				return nil, fmt.Errorf("flag %s is missing a value", fl)
			}
			v = a.FlagValue(fl)
			a.removeFlagAndValue(fl)
		}
		if err = flagTags.SetValue(fl, v); err != nil {
			return nil, err
		}
	}
	return a, nil
}

// NewCommand will parse the given argsuments into a new Command.
// The first argument is expected to be a known command.  This is mapped to the Command instance
// Any remaining arguments are applied to the commands flags.  Each flag used by the command is removed from the arguments.
// Remaining arguments, not consumed by the Command flags are returned with the new Command
func NewCommand(args ...string) (Command, Arguments, error) {
	if len(args) == 0 {
		return nil, nil, fmt.Errorf("no command given")
	}
	cmd := commandsMap[strings.ToLower(args[0])]
	if cmd == nil {
		return nil, nil, fmt.Errorf("%s is an unknown command", args[0])
	}
	arguments, err := ApplyArguments(cmd, args[1:])
	if err != nil {
		return nil, nil, err
	}
	return cmd, arguments, nil
}
