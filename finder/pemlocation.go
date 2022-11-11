package finder

import (
	"pempal/pemtypes"
)

// PemLocation represents a single data source (such as a file) and any PemResources found in it.
type PemLocation interface {
	Location() string
	Resources() []pemtypes.PemResource
}

type pemLocation struct {
	location  string
	resources []pemtypes.PemResource
}

func (p pemLocation) Location() string {
	return p.location
}
func (p pemLocation) Resources() []pemtypes.PemResource {
	return p.resources
}
