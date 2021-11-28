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
	"":          &ViewCommand{},
	"help":      &HelpCommand{},
	"key":       &KeyCommand{},
	"request":   &RequestCommand{},
	"issue":     &IssueCommand{},
	"revoke":    &RevokeCommand{},
	"keys":      &KeysCommand{},
	"issuers":   &IssuersCommand{},
	"find":      &FindCommand{},
	"view":      &ViewCommand{},
	"templates": &TemplatesCommand{},
	"template":  &TemplateCommand{},
}
var Aliases = map[string]string{
	"?":     "help",
	"fd":    "find",
	"cat":   "view",
	"temps": "templates",
	"temp":  "template",
	"req":   "request",
	"csr":   "request",
	"iss":   "issue",
	"rev":   "revoke",
	"crl":   "revoke",
}
