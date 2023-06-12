package commands

import (
	"fmt"
	"github.com/eurozulu/argdecoder"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/utils"
	"path/filepath"
	"strings"
)

// CommonFlags are flags which all command can use without the need to declare them in the command class.
var CommonFlags CommonFlagsStruct

// CommonFlagsStruct contains the flags used by all Commands
type CommonFlagsStruct struct {
	Out        string `flag:"out"`
	ForceOut   bool   `flag:"force,f"`
	ConfigPath string `flag:"config,cfg"`
	Verbose    bool   `flag:"v"`
	Debug      bool   `flag:"vv"`
	Quiet      bool   `flag:"quiet,q"`
	Help       bool   `flag:"help"`
}

func ApplyCommonFlags(args []string) ([]string, error) {
	newArgs, err := argdecoder.ApplyArguments(args, &CommonFlags)
	if err != nil {
		return nil, fmt.Errorf("Failed to read common flags  %v", err)
	}
	if CommonFlags.ConfigPath != "" && !strings.HasSuffix(CommonFlags.ConfigPath, ".config") && !utils.FileExists(CommonFlags.ConfigPath) {
		CommonFlags.ConfigPath = filepath.Join(CommonFlags.ConfigPath, ".config")
	}
	config.ConfigPath = CommonFlags.ConfigPath
	return newArgs, nil
}
