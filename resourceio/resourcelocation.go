package resourceio

import (
	"github.com/eurozulu/pempal/model"
)

type ResourceLocation interface {
	Location() string
	Resources(resourceType ...model.ResourceType) []model.Resource
}

type resourceLocation struct {
	location  string
	resources []model.Resource
}

func (rl resourceLocation) Location() string {
	return rl.location
}

func (rl resourceLocation) Resources(resourceType ...model.ResourceType) []model.Resource {
	if len(resourceType) == 0 {
		return rl.resources
	}
	var found []model.Resource
	for _, r := range rl.resources {
		if !isFilteredResourceType(r.ResourceType(), resourceType) {
			continue
		}
		found = append(found, r)
	}
	return found
}

func isFilteredResourceType(t model.ResourceType, ts []model.ResourceType) bool {
	for _, rt := range ts {
		if rt == t {
			return true
		}
	}
	return false
}

func NewResourceLocation(location string, resources []model.Resource) ResourceLocation {
	return &resourceLocation{
		location:  location,
		resources: resources,
	}
}
