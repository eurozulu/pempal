package commands

import (
	"context"
	"fmt"
	"io"
	"pempal/model"
	"pempal/resourceio"
	"strings"
)

const defaultOutputFormat = resourceio.LIST

type FindCommand struct {
	ResourceType string `flag:"resource-type,type"`
	Query        string `flag:"query, q"`
	Format       string `flag:"format,f"`
	Recursive    bool   `flag:"r,recursive"`
}

func (fc FindCommand) Execute(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("find requires at least one path to a file or directory.")
	}

	var filters []resourceio.LocationFilter

	typeFilter, err := fc.parseResourceType()
	if err != nil {
		return err
	}
	if typeFilter != nil {
		filters = append(filters, typeFilter)
	}

	queryFilter, err := fc.parseQuery()
	if err != nil {
		return err
	}
	if queryFilter != nil {
		filters = append(filters, queryFilter)
	}
	format, err := fc.parseFormat()
	if err != nil {
		return err
	}
	printOut := resourceio.NewResourcePrinter(out, format)

	rs := resourceio.NewResourceScanner(fc.Recursive, filters...)
	var count int
	for loc := range rs.Scan(context.Background(), args...) {
		if err := printOut.Write(*loc); err != nil {
			return err
		}
		count++
	}
	if CommonFlags.Out == "" {
		fmt.Fprintf(out, "found %d resources\n", count)
	}
	return nil
}

func (fc FindCommand) parseResourceType() (resourceio.LocationFilter, error) {
	if fc.ResourceType == "" {
		return nil, nil
	}
	var types []model.ResourceType
	for _, r := range strings.Split(fc.ResourceType, ",") {
		rt := model.ParseResourceType(r)
		if rt == model.Unknown {
			return nil, fmt.Errorf("%s is not a known resource type", r)
		}
		types = append(types, rt)
	}
	return resourceio.NewTypeResourceLocationFilter(types...), nil
}

func (fc FindCommand) parseQuery() (resourceio.LocationFilter, error) {
	if fc.Query == "" {
		return nil, nil
	}
	return nil, fmt.Errorf("query not yet implemented")
}

func (fc FindCommand) parseFormat() (resourceio.ResourceFormat, error) {
	if fc.Format == "" {
		return defaultOutputFormat, nil
	}
	return resourceio.ParseResourceFormat(fc.Format)
}
