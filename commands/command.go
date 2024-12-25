package commands

type Command interface {
	Exec(args ...string) error
}

var Commands = map[string]Command{
	"help":      &HelpCommand{},
	"make":      &MakeCommand{},
	"template":  &TemplateCommand{},
	"templates": &TemplatesCommand{},
	"validate":  &ValidateCommand{},
	"list":      &ListCommand{},
	"show":      &ShowCommand{},
	"config":    &ConfigCommand{},
	"issuers":   &IssuersCommand{},
}

var CommandAliases = map[string][]string{
	"h":     {"help"},
	"mk":    {"make"},
	"ls":    {"list"},
	"tp":    {"template"},
	"tps":   {"templates"},
	"vd":    {"validate"},
	"sh":    {"show"},
	"iss":   {"issuers"},
	"cfg":   {"config"},
	"cf":    {"config"},
	"keys":  {"list", "-type", "privatekey"},
	"certs": {"list", "-type", "certificate"},
}
