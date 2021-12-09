package commands

import (
	"context"
	"flag"
	"io"
)

// Command represents a single operation
type Command interface {
	Description() string
	Flags(f *flag.FlagSet)
	Run(ctx context.Context, out io.Writer, args ...string) error
}

// Commands maps the command keyword to the actual Command
var Commands = map[string]Command{
	"help":      &HelpCommand{},
	"keys":      &KeysCommand{},
	"key":       &KeyCommand{},
	"request":   &RequestCommand{},
	"issuers":   &IssuersCommand{},
	"issue":     &IssueCommand{},
	"revoke":    &RevokeCommand{},
	"find":      &FindCommand{},
	"templates": &TemplatesCommand{},
	"template":  &TemplateCommand{},
}

// Aliases lists alternative names for command keywords
var Aliases = map[string]string{
	"?":      "help",
	"search": "find",
	"fd":     "find",
	"list":   "find",
	"ls":     "find",
	"tp":     "template",
	"cat":    "template",
	"view":   "template",
	"temp":   "template",
	"temps":  "templates",
	"tps":    "templates",
	"re":     "request",
	"req":    "request",
	"csr":    "request",
	"is":     "issue",
	"rev":    "revoke",
	"crl":    "revoke",
}
