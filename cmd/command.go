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
	"help":      &HelpCommand{},
	"key":       &KeyCommand{},
	"keys":      &KeysCommand{},
	"issue":     &IssueCommand{},
	"issuers":   &IssuersCommand{},
	"find":      &FindCommand{},
	"view":      &ViewCommand{},
	"templates": &TemplatesCommand{},
	"template":  &TemplateCommand{},
}
var Aliases = map[string]string{
	"":    "view", // empty alias is the 'default' command, used when no command given
	"?":   "help",
	"fd":  "find",
	"cat": "view",
}
