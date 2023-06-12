package commands

import (
	"fmt"
	"github.com/eurozulu/argdecoder"
	"github.com/eurozulu/pempal/builder"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/resourceio"
	"io"
	"strings"
)

const defaultResourceFormat = resourceio.FormatPEM

type MakeCommand struct {
	BuildType string   `flag:"type,resource-type,resourcetype,buildtype"`
	Format    string   `flag:"format,fm"`
	flags     []string `flag:"-"`
}

func (cmd MakeCommand) Execute(args []string, out io.Writer) error {
	if len(args) == 0 {
		if err := showTemplateNames(out); err != nil {
			return err
		}
		return fmt.Errorf("must provide one or more template names to build")
	}

	tm, err := config.TemplateStore()
	if err != nil {
		return err
	}

	km, err := config.KeyManager()
	if err != nil {
		return err
	}

	temps, err := tm.ExtendedTemplatesByName(args...)
	if err != nil {
		return err
	}

	var rt model.ResourceType
	if cmd.BuildType != "" {
		rt = model.ParseResourceType(cmd.BuildType)
		if rt == model.Unknown {
			return fmt.Errorf("%s is an unknown build type", cmd.BuildType)
		}
	} else {
		rt = model.DetectResourceType(temps...)
		if rt == model.Unknown {
			return fmt.Errorf("failed to detect build type, ensure a template extends a resource type")
		}
	}

	build, err := builder.NewBuilder(rt, km)
	if err != nil {
		return err
	}

	if err = build.ApplyTemplate(temps...); err != nil {
		return err
	}
	if len(cmd.flags) > 0 {
		if _, err := argdecoder.ApplyArguments(cmd.flags, build); err != nil {
			return err
		}
	}

	// Keep validating until either no errors or user aborts (EOF)
	err = build.Validate()
	for err != nil {
		err = processInvalidTemplate(build, err)
		if err == io.EOF {
			return fmt.Errorf("aborted")
		}
	}

	r, err := build.Build()
	if err != nil {
		return err
	}

	if err = cmd.writeResource(out, r); err != nil {
		return err
	}
	return nil
}

func (cmd *MakeCommand) ApplyFlags(flags map[string]*string) error {
	// reassemble into commandline args
	var args []string
	for k, v := range flags {
		args = append(args, strings.Join([]string{"-", k}, ""))
		if v != nil {
			args = append(args, *v)
		}
	}
	cmd.flags = args
	return nil
}

func (cmd MakeCommand) writeResource(out io.Writer, r ...model.Resource) error {
	formOut, err := cmd.getResourceFormatter()
	if err != nil {
		return err
	}
	data, err := formOut.FormatResources(r...)
	if err != nil {
		return err
	}
	_, err = out.Write(data)
	return err
}

func (cmd MakeCommand) getResourceFormatter() (resourceio.ResourceFormatter, error) {
	var err error
	rf := defaultResourceFormat
	if cmd.Format != "" {
		rf, err = resourceio.ParseResourceFormat(cmd.Format)
		return nil, err
	}
	return resourceio.NewResourceFormatter(rf), nil
}

func showTemplateNames(out io.Writer) error {
	tc := &TemplateCommand{}
	return tc.Execute(nil, out)
}

func processInvalidTemplate(build builder.Builder, err error) error {
	logger.Error("make failed:\n%v", err)
	return io.EOF
}
