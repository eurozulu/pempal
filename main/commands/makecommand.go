package commands

import (
	"fmt"
	"io"
	"pempal/builders"
	"pempal/logger"
	"pempal/model"
)

const TAG_TYPE = "type"

type MakeCommand struct {
	Resource_Type string `flag:"resource-type,type"`
}

func (cmd MakeCommand) Execute(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("must provide one or more template names to build")
	}
	rt, err := cmd.establishResourceType(args)
	if err != nil {
		return err
	}
	if rt == model.Unknown {
		return fmt.Errorf("no resource type was detected. Ensure either one template has a '#%s' tag or the '-type' command line flag specifies the type", TAG_TYPE)
	}
	logger.Log(logger.Debug, "Make is building a new %s resource", rt.String())
	temps, err := ResourceTemplates.TemplatesByName(args...)
	if err != nil {
		return err
	}
	pl := isPlural(args)
	logger.Log(logger.Debug, "using template%s: %v", pl, args)

	builder := builders.NewResourceBuilder(rt)
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

func (cmd MakeCommand) establishResourceType(names []string) (model.ResourceType, error) {
	if cmd.Resource_Type == "" {
		// no forced type, detect type from template
		rt, err := findResourceTypeInTemplates(names)
		if err != nil {
			return model.Unknown, fmt.Errorf("failed to detect the resource type %v", err)
		}
		cmd.Resource_Type = rt.String()
	}
	rt := model.ParseResourceType(cmd.Resource_Type)
	return rt, nil
}

func findResourceTypeInTemplates(names []string) (model.ResourceType, error) {
	// Collect named templates and search for the first one with a Tag named TAG_TYPE
	temps, err := ResourceTemplates.TemplatesByName(names...)
	if err != nil {
		return 0, err
	}
	for _, t := range temps {
		tags := t.Tags().TagsByName(TAG_TYPE)
		if len(tags) == 0 {
			continue
		}
		rt := model.ParseResourceType(tags[0].Value)
		if rt == model.Unknown {
			return model.Unknown, fmt.Errorf("resource type %s is unknown", tags[0].Value)
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
