package validation

import "fmt"

type ValidationError struct {
	PropertyName string
	Message      string
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("invalid %q: %s", v.PropertyName, v.Message)
}
