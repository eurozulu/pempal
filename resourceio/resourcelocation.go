package resourceio

import (
	"github.com/eurozulu/pempal/resources"
)

type ResourceLocation interface {
	Location() string
	Resources(resourceType ...resources.ResourceType) []resources.Resource
}

type resourceLocation struct {
	location  string
	resources []resources.Resource
}

func (rl resourceLocation) Location() string {
	return rl.location
}

func (rl resourceLocation) Resources(resourceType ...resources.ResourceType) []resources.Resource {
	if len(resourceType) == 0 {
		return rl.resources
	}
	var found []resources.Resource
	for _, r := range rl.resources {
		if !isFilteredResourceType(r.ResourceType(), resourceType) {
			continue
		}
		found = append(found, r)
	}
	return found
}

func isFilteredResourceType(t resources.ResourceType, ts []resources.ResourceType) bool {
	for _, rt := range ts {
		if rt == t {
			return true
		}
	}
	return false
}

func NewResourceLocation(location string, resources []resources.Resource) ResourceLocation {
	return &resourceLocation{
		location:  location,
		resources: resources,
	}
}
