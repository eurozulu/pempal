//go:generate ../spud gen . -packagename commands -spudpackage commands -f -v

package commands

import (
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/templates"
	"gopkg.in/yaml.v2"
	"strconv"
	"strings"
)

// Verbose, when true, displays additional data about each template
// @Flag(verbose,v)
var Verbose bool

// @Flag(vv)
var VeryVerbose bool

func SetLoggingOutput() {
	if Verbose {
		logging.DefaultLogger.SetLogLevel(logging.LogInfo)
	}
	if VeryVerbose {
		logging.DefaultLogger.SetLogLevel(logging.LogDebug)
	}
}

func ArgFlagsToTemplate(args []string) (templates.Template, []string, error) {
	flags := map[string]interface{}{}
	var argz []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") {
			argz = append(argz, arg)
			continue
		}
		arg = strings.TrimPrefix(arg, "-")
		if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
			flags[arg] = true
			continue
		}
		var v interface{}
		if i+1 < len(args) {
			i++
			v = stringToValue(args[i])
		}
		flags[arg] = v
	}
	var data []byte
	if len(flags) > 0 {
		var err error
		data, err = yaml.Marshal(flags)
		if err != nil {
			return nil, nil, err
		}
	}
	return &model.TemplateFile{
		Path: "",
		Data: data,
	}, argz, nil
}

func stringToValue(s string) interface{} {
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
}
