package templates

import "pempal/resources"

type Template interface {
	Type() resources.ResourceType
}

type EmptyTemplate map[string]interface{}

func (e EmptyTemplate) Type() resources.ResourceType {
	return resources.Unknown
}
