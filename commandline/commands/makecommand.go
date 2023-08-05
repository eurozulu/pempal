package commands

import (
	"fmt"
	"github.com/eurozulu/pempal/builders"
	"github.com/eurozulu/pempal/commandline/commonflags"
	"github.com/eurozulu/pempal/commandline/prompts"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
	"io"
)

type MakeCommand struct {
	flagsTemplate templates.Template
}

func (cmd MakeCommand) Execute(args []string, out io.Writer) error {
	tb := templates.NewTemplateBuilder()

	// check if any named templates or file resources given in args
	argTemps, err := argumentsToTemplates(args)
	if err != nil {
		return err
	}
	tb.Add(argTemps...)
	if cmd.flagsTemplate != nil {
		tb.Add(cmd.flagsTemplate)
	}

	rt, err := cmd.detectResourceType(tb.Build())
	if err != nil {
		return err
	}

	if !commonflags.CommonFlags.Quiet {
		t, err := cmd.confirmBuild(rt, tb.Build())
		if err != nil {
			return err
		}
		if t != nil {
			tb.Add(t)
		}
	} else {
		if err := cmd.validate(rt, tb.Build()); err != nil {
			return err
		}
	}
	return PerformBuild(rt, tb.Build(), out)
}

func (cmd MakeCommand) validate(rt resources.ResourceType, t templates.Template) error {
	build, err := builders.NewBuilder(rt)
	if err != nil {
		return err
	}
	return build.Validate(t)
}

func (cmd MakeCommand) detectResourceType(t templates.Template) (resources.ResourceType, error) {
	rts, err := resources.TemplateTypes(t)
	if err != nil || len(rts) == 0 {
		// no type known.  If quiet, report as error, otherwise request type
		if commonflags.CommonFlags.Quiet {
			return 0, err
		}
		rts = []resources.ResourceType{}
	}
	if len(rts) != 1 {
		if commonflags.CommonFlags.Quiet {
			return 0, fmt.Errorf("template has unknown or ambiguous resource type(s) %v", rts)
		}
		rt, err := cmd.showRequestType(rts)
		if err != nil {
			return 0, err
		}
		rts = []resources.ResourceType{rt}
	}
	return rts[0], nil
}

func (cmd MakeCommand) confirmBuild(rt resources.ResourceType, t templates.Template) (templates.Template, error) {
	confirm, err := prompts.NewConfirmBuild(rt)
	if err != nil {
		return nil, err
	}
	return confirm.Confirm(t)
}

func (cmd MakeCommand) showRequestType(types []resources.ResourceType) (resources.ResourceType, error) {
	var text string
	if len(types) > 0 {
		text = types[0].String()
	}
	return prompts.NewRequestTypeView("resource type", text, types...).Show()
}

func (cmd *MakeCommand) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&cmd.flagsTemplate)
}

func (cmd MakeCommand) MarshalYAML() (interface{}, error) {
	return &cmd.flagsTemplate, nil
}
