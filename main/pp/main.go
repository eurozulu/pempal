package main

import (
	"fmt"
	"io"
	"os"
	"pempal/logger"
	"pempal/main/argdecoder"
	"pempal/main/commands"
	"pempal/main/help"
	"pempal/utils"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		logger.Log(logger.Error, "requires at least one command.\n")
		fmt.Fprintln(os.Stdout, help.HelpCommands())
		return
	}
	// use first param as command name
	cmd, err := commands.NewCommand(os.Args[1])
	if err != nil {
		logger.Log(logger.Error, "%v", err)
		return
	}

	// Apply common flags
	args, err := commands.ApplyCommonFlags(os.Args[2:])
	if err != nil {
		logger.Log(logger.Error, "%v", err)
		return
	}

	// apply remaining flags to command
	args, err = argdecoder.ApplyArguments(args, cmd)
	if err != nil {
		logger.Log(logger.Error, "%v\n", err)
		return
	}

	// remove any remaining flags from args
	params, flags := argdecoder.ParseArgs(args)
	if len(flags) > 0 {
		// Remaining flags, which command (or common) did not consume
		// If command supports custom flags, pass them, otherwise error
		cmdFlag, ok := cmd.(commands.CommandWithFlags)
		if !ok {
			logger.Log(logger.Error, "unknown flag %v\n", mapKeys(flags))
			return
		}
		if err = cmdFlag.SetFlags(flags); err != nil {
			logger.Log(logger.Error, "%v", err)
			return
		}
	}

	// establish the output stream
	out := os.Stdout
	if commands.CommonFlags.Out != "" {
		if utils.FileExists(commands.CommonFlags.Out) && !commands.CommonFlags.ForceOut {
			logger.Log(logger.Error, "%s already exists\n", commands.CommonFlags.Out)
			return
		}
		f, err := os.OpenFile(commands.CommonFlags.Out, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Log(logger.Error, "failed to open %s  %v\n", commands.CommonFlags.Out, err)
			return
		}
		out = f
		defer func(out io.WriteCloser) {
			if err := out.Close(); err != nil {
				logger.Log(logger.Error, "failed to close %s  %v\n", commands.CommonFlags.Out, err)
			}
		}(f)
	}

	level := logger.Info
	if commands.CommonFlags.Verbose {
		level = logger.Warning
	}
	if commands.CommonFlags.Debug {
		level = logger.Debug
	}
	logger.DefaultLogger.SetLevel(level)

	if err = cmd.Execute(params, out); err != nil {
		logger.Log(logger.Error, "%v\n", err)
	}
}

func mapKeys(m map[string]*string) []string {
	var index int
	keys := make([]string, len(m))
	for k := range m {
		keys[index] = strings.Join([]string{"-", k}, "")
		index++
	}
	return keys
}
