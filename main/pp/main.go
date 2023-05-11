package main

import (
	"fmt"
	"github.com/eurozulu/argdecoder"
	"io"
	"os"
	"pempal/logger"
	"pempal/main/commands"
	"pempal/main/help"
	"pempal/utils"
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

	// Check if command parses custom flags
	cmdFlag, ok := cmd.(commands.FlagsParsingCommand)
	if ok {
		args, err = cmdFlag.ApplyFlags(args)
		if err != nil {
			logger.Log(logger.Error, "%v\n", err)
			return
		}
	}

	// Check all flags have been consumed
	_, flags := argdecoder.ParseArgs(args)
	if len(flags) > 0 {
		logger.Log(logger.Error, "unknown flag(s) %v\n", unknownFlagNames(flags))
		return
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

	// Set logging level output
	level := logger.Info
	if commands.CommonFlags.Verbose {
		level = logger.Warning
	}
	if commands.CommonFlags.Debug {
		level = logger.Debug
	}
	logger.DefaultLogger.SetLevel(level)

	if err = cmd.Execute(args, out); err != nil {
		logger.Log(logger.Error, "%v\n", err)
	}
}

func unknownFlagNames(m map[string]*string) []string {
	names := make([]string, len(m))
	var index int
	for k := range m {
		names[index] = k
	}
	return names
}
