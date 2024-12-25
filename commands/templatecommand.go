package commands

import (
	"fmt"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/utils"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strings"
)

type TemplateCommand struct {
	Output io.Writer
}

func (cmd TemplateCommand) Exec(args ...string) error {
	ft, err := cmd.BuildTemplate(args...)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(ft)
	if err != nil {
		return err
	}
	out := cmd.Output
	if out == nil {
		out = os.Stdout
		fmt.Fprintf(out, "%q builds a %s template:\n", strings.Join(args, ", "), ft.Name())
	}
	_, err = out.Write(data)
	return err
}

func (cmd TemplateCommand) BuildTemplate(args ...string) (templates.Template, error) {
	builder := templates.NewTemplateBuilder(config.Config.Templates...)

	var flags map[string]interface{}
	args, flags = readFlagArgs(args)

	// add the named templates from the (non flag) arguments
	if err := builder.AddTemplate(args...); err != nil {
		return nil, err
	}
	by, err := cmd.flagsAsTemplate(flags)
	if err != nil {
		return nil, err
	}
	if len(by) > 0 {
		if err = builder.AddNewTemplate("command-line-flags", by); err != nil {
			return nil, err
		}
	}
	return builder.Build()
}

func readFlagArgs(args []string) ([]string, map[string]interface{}) {
	flags := map[string]interface{}{}
	var remain []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") {
			remain = append(remain, arg)
			continue
		}
		var value interface{}
		if i+1 < len(args) {
			value = utils.StringToType(strings.TrimSpace(args[i+1]))
			i++
		}
		flags[strings.TrimLeft(arg, "-")] = value
	}
	return remain, flags
}

func (vc TemplateCommand) flagsAsTemplate(flags map[string]interface{}) ([]byte, error) {
	if len(flags) == 0 {
		return nil, nil
	}
	return yaml.Marshal(&flags)
}
