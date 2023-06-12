package resourceio

import (
	"fmt"
	"strings"
)

// ResourceFormat is a specific format type which resources may be transformed into
type ResourceFormat int

const (
	FormatPEM ResourceFormat = iota
	FormatDER
	FormatYAML
)

var formatNames = []string{
	"FormatPEM",
	"FormatDER",
	"FormatYAML",
}

func (r ResourceFormat) String() string {
	return formatNames[r]
}

func ParseResourceFormat(s string) (ResourceFormat, error) {
	for i, name := range formatNames {
		if strings.EqualFold(s, name) {
			return ResourceFormat(i), nil
		}
	}
	return FormatPEM, fmt.Errorf("%s is an unknown resource format", s)
}
