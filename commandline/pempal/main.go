package main

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/pempal/commandline/commands"
	"github.com/eurozulu/pempal/commandline/commonflags"
	"github.com/eurozulu/pempal/commandline/help"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/utils"
	"io"
	"os"
)

func main() {
	args, flags := commands.ParseArgs(os.Args[1:])
	if err := flags.ApplyTo(commonflags.CommonFlags); err != nil {
		fmt.Fprintln(os.Stderr, "Invalid arguments  %v", err)
		return
	}

	// Set logging level output
	setLoggerLevel()

	var arg string
	if len(args) > 0 {
		arg = args[0]
		args = args[1:]
	}
	// if help flag specified, show help on command or general help if no command given
	if commonflags.CommonFlags.Help || arg == "" {
		fmt.Fprintln(os.Stdout, help.HelpFor(arg))
		return
	}

	// create Command with the first argument
	cmd, err := commands.NewCommand(arg)
	if err != nil {
		logger.Error("%v", err)
		return
	}

	// apply any flags to command
	if err = flags.ApplyTo(cmd); err != nil {
		logger.Error("error parsing flags for %s  %v\n", arg, err)
		return
	}

	if len(flags) > 0 {
		logger.Error("unknown flag '%s' for command %s\n", badFlagNames(flags), arg)
		return
	}
	// establish the output stream
	out := os.Stdout
	if commonflags.CommonFlags.Output != nil {
		outPath := *commonflags.CommonFlags.Output
		if outPath != "" {
			if utils.FileExists(outPath) && !commonflags.CommonFlags.ForceOut {
				logger.Error("%s already exists\n", outPath)
				return
			}
			f, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				logger.Error("failed to open %s  %v\n", outPath, err)
				return
			}
			out = f
			defer func(out io.WriteCloser) {
				if err := out.Close(); err != nil {
					logger.Error("failed to close %s  %v\n", outPath, err)
				}
			}(f)
		}
	}
	if err = cmd.Execute(args, out); err != nil {
		logger.Error("%v\n", err)
	}
}

func setLoggerLevel() {
	level := logger.LevelInfo
	if commonflags.CommonFlags.Quiet {
		level = logger.LevelError
	}
	if commonflags.CommonFlags.Verbose {
		level = logger.LevelWarning
	}
	if commonflags.CommonFlags.Debug {
		level = logger.LevelDebug
	}
	logger.DefaultLogger.SetLevel(level)
	logger.DefaultLogger.SetShowTimeStamp(level == logger.LevelDebug)
	if commonflags.CommonFlags.Quiet && (commonflags.CommonFlags.Verbose || commonflags.CommonFlags.Debug) {
		logger.Warning("Ignoring -quiet flag as -verbose or -debug is active")
	}
}

func badFlagNames(flags map[string]*string) string {
	buf := bytes.NewBuffer(nil)
	var nonFirst bool
	for k := range flags {
		if nonFirst {
			buf.WriteString(", ")
		} else {
			nonFirst = true
		}
		buf.WriteRune('-')
		buf.WriteString(k)
	}
	return buf.String()
}
