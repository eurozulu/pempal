package commands

import (
	"fmt"
	"io"
	"pempal/builders"
	"pempal/logger"
	"pempal/model"
	"pempal/templates"
)

type MakeCommand struct {
	Resource_Type string `flag:"resource-type,type"`
}

func (cmd MakeCommand) Execute(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("must provide one or more template names to build")
	}
	temps, err := parseArgumentsToTemplates(args)
	if err != nil {
		return err
	}

	pl := isPlural(args)
	logger.Log(logger.Debug, "using template%s: %v", pl, args)

	rt, err := cmd.resolveResourceType(temps)
	if err != nil {
		return err
	}
	logger.Log(logger.Debug, "Make is building a new %s resource", rt.String())

	builder, err := builders.NewResourceBuilder(rt)
	if err != nil {
		return err
	}

	if err = builder.ApplyTemplate(temps...); err != nil {
		return err
	}
	logger.Log(logger.Debug, "Templates set, building resource...")
	pemRes, err := builder.Build()
	if err != nil {
		return err
	}
	logger.Log(logger.Debug, "building successful of %s resource", pemRes.ResourceType().String())
	pemBytes, err := pemRes.MarshalPEM()
	if err != nil {
		return err
	}
	_, err = out.Write(pemBytes)
	return err
}

func (cmd MakeCommand) resolveResourceType(temps []templates.Template) (model.ResourceType, error) {
	if cmd.Resource_Type != "" {
		rt := model.ParseResourceType(cmd.Resource_Type)
		if rt == model.Unknown {
			return 0, fmt.Errorf("resource type %s is not known", cmd.Resource_Type)
		}
		return rt, nil
	}

	// no forced type, detect type from template
	// search for the first template with a Tag named TAG_TYPE
	rts := struct {
		ResourceType string `yaml:"resource-type"`
	}{}
	for _, t := range temps {
		rts.ResourceType = ""
		if err := t.Apply(&rts); err != nil {
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

func isPlural(args []string) string {
	if len(args) == 1 {
		return ""
	}
	return "s"
}
