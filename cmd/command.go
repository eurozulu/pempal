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
	"help": &HelpCommand{},

	"key":  &KeyCommand{},
	"keys": &KeysCommand{},

	"issue":   &IssueCommand{},
	"issuers": &IssuersCommand{},

	"find":      &FindCommand{},
	"view":      &ViewCommand{},
	"templates": &TemplatesCommand{},
}
var Aliases = map[string]string{
	"?":   "help",
	"fd":  "find",
	"cat": "view",
}
