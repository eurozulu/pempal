package commands

import (
	"fmt"
	"io"
	"pempal/formats"
	"pempal/templates"
)

type makeCommand struct {
	Quiet bool `flag:"quiet"`
	quiet bool `flag:"q"`
}

func (m makeCommand) Main(args Arguments, out io.Writer) error {
	var builder formats.Builder

	temps, err := templateManager.FindTemplates(args.Parameters()...)
	if err != nil {
		return err
	}
	if len(temps) > 0 {
		if err = builder.AddTemplate(temps...); err != nil {
			return err
		}
	}
	// loop until we have all the required argument or aborted
	argTemplate := argsToTemplate(args)
	for {
		if len(argTemplate) > 0 {
			if err = builder.AddTemplate(temps...); err != nil {
				return err
			}
		}
		missing := builder.Validate()
		if len(missing) == 0 {
			// all required has values, continue
			break
		}
		// If non interactive mode, throw error for missing values
		if m.Quiet || m.quiet {
			return fmt.Errorf("required values not found.  %v", missing)
		}

		// request missing values from user
		for _, name := range missing {
			k := name.Error()[len("missing "):] // strip leading 'missing '
			v, err := requestMissing(k)
			if err != nil {
				return err
			}
			argTemplate[k] = v
		}
	}
}

func requestMissing(name string) (string, error) {
	prompt := prompertyPrompts[name]
	if prompt == nil {
		prompt = stringPrompt{}
	}
	prompt.Request(fmt.Sprintf("enter %s:", name))
}

func argsToTemplate(args Arguments) templates.EmptyTemplate {
	names := args.FlagNames()
	if len(names) == 0 {
		return nil
	}
	t := templates.EmptyTemplate{}
	for _, k := range names {
		t[k] = args.FlagValue(k)
	}
	return t
}
