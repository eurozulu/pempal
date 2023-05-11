package commands

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/argdecoder"
	"github.com/go-yaml/yaml"
	"io"
	"os"
	"pempal/builders"
	"pempal/logger"
	"pempal/model"
	"pempal/templates"
)

type MakeCommand struct {
	Resource_Type string `flag:"resource-type,type"`
	NoPrompt      bool   `flag:"no-prompt,q"`

	resource_Type model.ResourceType
	flagsTemplate templates.Template
}

func (cmd *MakeCommand) ApplyFlags(args []string) ([]string, error) {
	var remain []string
	var err error
	remain, cmd.flagsTemplate, err = cmd.buildArgumentTemplate(args)
	return remain, err
}

func (cmd MakeCommand) Execute(args []string, out io.Writer) error {
	if cmd.resource_Type == model.Unknown {
		rt, err := cmd.resolveResourceType(args)
		if err != nil {
			return err
		}
		cmd.resource_Type = rt
	}
	logger.Log(logger.Debug, "Make is building a new %s resource", cmd.resource_Type.String())

	if len(args) == 0 {
		if err := showTemplateNames(os.Stdout); err != nil {
			return err
		}
		return fmt.Errorf("must provide one or more template names to build")
	}
	builder, err := builders.NewResourceBuilder(cmd.resource_Type)
	if err != nil {
		return err
	}
	pl := isPlural(args)
	logger.Log(logger.Debug, "using template%s: %v", pl, args)

	temps, err := loadNamedTemplates(args)
	if err != nil {
		return err
	}
	if cmd.flagsTemplate != nil {
		temps = append(temps, cmd.flagsTemplate)
	}
	logger.Log(logger.Debug, "Build templates loaded, applying to resource template...")
	if err = builder.ApplyTemplate(temps...); err != nil {
		return err
	}

	logger.Log(logger.Debug, "Resource template loaded, validating...")
	errs := builder.Validate()
	for len(errs) > 0 {
		buf := bytes.NewBuffer(nil)
		for _, err := range errs {
			buf.WriteString(err.Error())
			buf.WriteRune('\n')
		}
		return fmt.Errorf("%s", buf.String())
	}

	logger.Log(logger.Debug, "Resource template valid, starting build of %s...", cmd.resource_Type.String())
	pemRes, err := builder.Build()
	if err != nil {
		return err
	}
	logger.Log(logger.Debug, "%s  build successful", pemRes.ResourceType().String())
	pemBytes, err := pemRes.MarshalPEM()
	if err != nil {
		return err
	}
	_, err = out.Write(pemBytes)
	return err
}

func (cmd MakeCommand) resolveResourceType(args []string) (model.ResourceType, error) {
	if cmd.Resource_Type != "" {
		rt := model.ParseResourceType(cmd.Resource_Type)
		if rt == model.Unknown {
			return 0, fmt.Errorf("resource type %s is not known", cmd.Resource_Type)
		}
		return rt, nil
	}

	// no forced type, detect type from template
	// search for the first template with 'resource-type' property
	names, _ := argdecoder.ParseArgs(args)
	rts := struct {
		ResourceType string `yaml:"resource-type"`
	}{}
	for _, name := range names {
		t, err := ResourceTemplates.TemplatesByName(name)
		if err != nil {
			return 0, err
		}
		if err := t[0].Apply(&rts); err != nil {
			return 0, err
		}
		if rts.ResourceType == "" {
			continue
		}
		rt := model.ParseResourceType(rts.ResourceType)
		if rt == model.Unknown {
			return model.Unknown, fmt.Errorf("resource type %s is unknown", rts.ResourceType)
		}
		return rt, nil
	}
	return model.Unknown, nil
}

func (cmd *MakeCommand) buildArgumentTemplate(args []string) ([]string, templates.Template, error) {
	// establish resource type and apply argument flags to the corresponding DTO
	var err error
	cmd.resource_Type, err = cmd.resolveResourceType(args)
	if err != nil {
		return args, nil, err
	}

	// Use two DTO's and compare diffs to filter out zero value properties
	dto := model.DTOForResourceType(cmd.resource_Type)
	if dto == nil {
		return args, nil, fmt.Errorf("failed to identify resource type %s", cmd.resource_Type.String())
	}
	dtoEmpty := model.DTOForResourceType(cmd.resource_Type)

	remain, err := argdecoder.ApplyArguments(args, dto)
	if err != nil {
		return args, nil, err
	}
	t, err := dtoDiffAsTemplate(dtoEmpty, dto)
	return remain, t, err
}

func showTemplateNames(out io.Writer) error {
	tc := &TemplatesCommand{}
	return tc.Execute(nil, out)
}

func isPlural(args []string) string {
	if len(args) == 1 {
		return ""
	}
	return "s"
}

func dtoDiffAsTemplate(dtoBase, dto model.DTO) (templates.Template, error) {
	m, err := model.DTOToMap(dto)
	mBase, err := model.DTOToMap(dtoBase)
	diff := map[string]interface{}{}
	for k, v := range m {
		if mBase[k] == v {
			continue
		}
		diff[k] = v
	}
	if len(diff) == 0 {
		return nil, nil
	}

	data, err := yaml.Marshal(&diff)
	if err != nil {
		return nil, err
	}
	return ResourceTemplates.ParseTemplate(data)
}
