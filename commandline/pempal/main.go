package main

import (
	"fmt"
	"github.com/eurozulu/argdecoder"
	"github.com/eurozulu/pempal/commandline/commands"
	"github.com/eurozulu/pempal/commandline/help"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/utils"
	"io"
	"os"
)

func main() {
	// Apply common flags
	args, err := commands.ApplyCommonFlags(os.Args[1:])
	if err != nil {
		logger.Error("%v", err)
		return
	}

	// Set logging level output
	setLoggerLevel()

	// use first arg as command, if present
	var cmdArg string
	if len(args) > 0 {
		cmdArg = args[0]
		args = args[1:]
	}

	// if help flag specified, show help on command or general help if no command given
	if commands.CommonFlags.Help {
		fmt.Fprintln(os.Stdout, help.HelpFor(cmdArg))
		return
	}

	if cmdArg == "" {
		fmt.Fprintln(os.Stdout, "requires at least one command. Use -help to list available commands")
		return
	}

	// create Command and apply flags
	cmd, err := commands.NewCommand(cmdArg)
	if err != nil {
		logger.Error("%v", err)
		return
	}

	// apply any flags to command
	args, err = argdecoder.ApplyArguments(args, cmd)
	if err != nil {
		logger.Error("%v\n", err)
		return
	}

	// Ensure all flags have been consumed
	if ags, flags := argdecoder.ParseArgs(args); len(flags) > 0 {
		if flagCmd, ok := cmd.(commands.CustomFlagsCommand); ok {
			flagCmd.ApplyFlags(flags)
		} else {
			logger.Error("unknown flag(s) %v\n", unknownFlagNames(flags))
			return
		}
		args = ags
	}

	// establish the output stream
	out := os.Stdout
	if commands.CommonFlags.Out != "" {
		if utils.FileExists(commands.CommonFlags.Out) && !commands.CommonFlags.ForceOut {
			logger.Error("%s already exists\n", commands.CommonFlags.Out)
			return
		}
		f, err := os.OpenFile(commands.CommonFlags.Out, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			logger.Error("failed to open %s  %v\n", commands.CommonFlags.Out, err)
			return
		}
		out = f
		defer func(out io.WriteCloser) {
			if err := out.Close(); err != nil {
				logger.Error("failed to close %s  %v\n", commands.CommonFlags.Out, err)
			}
		}(f)
	}

	if err = cmd.Execute(args, out); err != nil {
		logger.Error("%v\n", err)
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

func setLoggerLevel() {
	level := logger.LevelInfo
	if commands.CommonFlags.Verbose {
		level = logger.LevelWarning
	}
	if commands.CommonFlags.Debug {
		level = logger.LevelDebug
	}
	logger.DefaultLogger.SetLevel(level)
	logger.DefaultLogger.SetShowTimeStamp(level > logger.LevelInfo)
}
