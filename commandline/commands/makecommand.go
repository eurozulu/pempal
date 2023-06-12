package commands

import (
	"fmt"
	"github.com/eurozulu/pempal/builder"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/resourceio"
	"github.com/eurozulu/pempal/templates"
	"github.com/go-yaml/yaml"
	"io"
)

const defaultResourceFormat = resourceio.FormatPEM

type MakeCommand struct {
	BuildType string             `flag:"type,resource-type,resourcetype,buildtype"`
	Format    string             `flag:"format,fm"`
	flags     templates.Template `flag:"-"`
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
	} else {
		rt = model.DetectResourceType(temps...)
	}
	if rt == model.Unknown {
		return fmt.Errorf("unknown resource type")
	}

	build, err := builder.NewBuilder(rt, km)
	if err != nil {
		return err
	}
	if err = build.ApplyTemplate(temps...); err != nil {
		return err
	}
	if cmd.flags != nil {
		if err = build.ApplyTemplate(cmd.flags); err != nil {
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
	m := model.FlatMap(flags).Expand()
	data, err := yaml.Marshal(&m)
	if err != nil {
		return err
	}
	t, err := templates.NewTemplate(data)
	if err != nil {
		return err
	}
	cmd.flags = t
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
	return io.EOF
}
