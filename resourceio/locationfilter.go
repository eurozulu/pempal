package resourceio

import (
	"fmt"
	"pempal/model"
	"regexp"
	"strings"
)

type LocationFilter interface {
	Filter(loc ResourceLocation) []model.PEMResource
}

type ResourceTypeLocationFilter struct {
	resourceTypes []model.ResourceType
}

func (t ResourceTypeLocationFilter) Filter(loc ResourceLocation) []model.PEMResource {
	var filtered []model.PEMResource
	for _, r := range loc.Resources {
		if !model.ContainsType(r.ResourceType(), t.resourceTypes) {
			continue
		}
		filtered = append(filtered, r)
	}
	return filtered
}

type FileNameLocationFilter struct {
	exp *regexp.Regexp
}

func (n FileNameLocationFilter) Filter(loc *ResourceLocation) []model.PEMResource {
	if n.exp.MatchString(loc.Path) {
		return loc.Resources
	}
	return nil
}

func NewTypeResourceLocationFilter(resourceTypes ...model.ResourceType) *ResourceTypeLocationFilter {
	return &ResourceTypeLocationFilter{resourceTypes: resourceTypes}
}

func NewFileNameResourceLocationFilter(name string) (*FileNameLocationFilter, error) {
	expr := strings.Replace(name, ".", "\\.", -1)
	expr = strings.Replace(expr, "*", ".", -1)
	exp, err := regexp.Compile(expr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse name into expression  %v", err)
	}
	return &FileNameLocationFilter{exp: exp}, nil
}
