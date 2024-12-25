package main

import (
	"fmt"
	"github.com/eurozulu/argflags"
	"github.com/eurozulu/pempal/commands"
	"os"
)

func main() {
	var err error
	args := argflags.ArgFlags(os.Args[1:])
	// Apply any common flags
	args, err = args.ApplyTo(&commands.CommonFlags)
	if len(args) == 0 || commands.CommonFlags.Help {
		// when help flag given or no command given, force the help command to run
		args = append([]string{"help"}, args...)
	}
	commands.SetLogging()
	if replace, ok := commands.CommandAliases[args[0]]; ok {
		args = append(replace, args[1:]...)
	}
	cmd := commands.Commands[args[0]]
	if cmd == nil {
		exitError(fmt.Errorf("unknown command: %s", args[0]))
	}

	// apply the command specific flags
	args, err = args[1:].ApplyTo(cmd)
	if err != nil {
		exitError(err)
	}

	if err = cmd.Exec(args...); err != nil {
		exitError(err)
	}
}

func exitError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
