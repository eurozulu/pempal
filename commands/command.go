package commands

import (
	"context"
	"io"
)

type Command interface {
	Run(cxt context.Context, args Arguments, out io.Writer) error
}

type Arguments interface {
	FlagNames() []string
	FlagValue(name string) string
	FlagBoolValue(name string) string

	Parameters() []string
}
