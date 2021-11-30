package commands

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
	"":          &HelpCommand{},
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
	"fd":    "find",
	"list":  "find",
	"ls":    "find",
	"cat":   "view",
	"temps": "templates",
	"temp":  "template",
	"req":   "request",
	"csr":   "request",
	"iss":   "issue",
	"rev":   "revoke",
	"crl":   "revoke",
}

var VerboseFlag bool
var TimeRunFlag bool
var OutFileFlag string
var HelpFlag bool

func FlagsMain(f *flag.FlagSet) {
	f.BoolVar(&VerboseFlag, "verbose", false, "Display all logging whilst searching for pems")
	f.BoolVar(&VerboseFlag, "v", false, "same as verbose")
	f.BoolVar(&TimeRunFlag, "t", false, "Times how long the command takes to execute")
	f.StringVar(&OutFileFlag, "out", "", "Specify a filename to write output into. Defaults to stdout")
	f.BoolVar(&HelpFlag, "help", false, "Display help")
}
