package cmd

import (
	"context"
	"flag"
	"io"
)

type Command interface {
	Description() string
	Flags(f *flag.FlagSet)
	Run(ctx context.Context, out io.Writer, args ...string) error
}

var Commands = map[string]Command{
	"?":    &HelpCommand{},
	"help": &HelpCommand{},
	"key":  &KeyCommand{},
	"keys": &KeysCommand{},

	"issue":   &IssueCommand{},
	"issuers": &IssuersCommand{},

	"list":      &ListCommand{},
	"show":      &ShowCommand{},
	"templates": &TemplatesCommand{},
}
