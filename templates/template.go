package templates

import "pempal/pemresources"

type Template interface {
	// Type is the PEM type of the template.
	Type() string
	pemresources.Marshaler
	pemresources.Unmarshaler
}
