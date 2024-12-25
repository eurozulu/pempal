package commands

import "github.com/eurozulu/pempal/logging"

var CommonFlags CommonFlagsStruct

type CommonFlagsStruct struct {
	Quiet       bool `flag:"quiet,q"`
	Help        bool `flag:"help,h"`
	Verbose     bool `flag:"verbose,v"`
	VeryVerbose bool `flag:"very-verbose,vv"`
}

func SetLogging() {
	level := logging.DefaultLogger.LogLevel()
	if CommonFlags.Verbose {
		level = logging.LogInfo
	}
	if CommonFlags.VeryVerbose {
		level = logging.LogDebug
	}
	logging.DefaultLogger.SetLogLevel(level)
	logging.Info("logging", "set logging to level %s", level)
}
