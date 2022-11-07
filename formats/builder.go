package formats

import (
	"pempal/resources"
	"pempal/templates"
)

type Builder interface {
	SetLocation(l string)

	// Adds the (relevant) properties from the given template, into the builders collection of properties.
	AddTemplate(ts ...templates.Template) error

	// Template will return the current state of the template to be used to build the resource
	Template() templates.Template

	Validate() []error

	// Build will attempt to generate the new resource(s), based on the template.
	Build() (resources.Resources, error)
}
