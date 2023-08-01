package help

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/pempal/commandline/commands"
	"github.com/eurozulu/pempal/utils"
	"os"
	"path/filepath"
	"sort"
)

var helpPages = map[string]HelpPage{
	"":         commandsHelpPage,
	"find":     findHelpPage,
	"config":   configHelpPage,
	"template": templateHelpPage,
}

var commandsHelpPage = HelpPage{
	Title:   "Pempal commands",
	Aliases: nil,
	Format: fmt.Sprintf("%s <command> [parameters] [-common flags] [-command specific flags]",
		filepath.Base(os.Args[0])),
	Description: "\nThe command must be the first argument, before any other flags or arguments\n" +
		"parameters and flags may follow, depending on the command. the following flags and parameters can be in any order.\n" +
		"Common Flags may be used with any command.  Common flags are:\n",
	Flags: []HelpPageFlag{
		{
			Name:        "out",
			Description: "specify a file path to write the result into.  Default is the standard console output",
		},
		{
			Name:        "config",
			Description: "specify a specific configuration file to use. See config help for details.",
			Aliases:     []string{"cfg"},
		},
		{
			Name:        "force",
			Description: "used with the -out flag to force overwritting an existing file",
			Aliases:     []string{"f"},
		},
		{
			Name:        "verbose",
			Description: "when present, displays detailed messages of the operation",
			Aliases:     []string{"v"},
		},
		{
			Name:        "debug",
			Description: "when present, displays very detailed messages of the operation",
			Aliases:     []string{"vv"},
		},
		{
			Name:        "quiet",
			Description: "limits to output of each command to the bare minimal",
			Aliases:     []string{"q"},
		},
		{
			Name:        "help",
			Description: "display help.  If a command is given, detailed help on that command is provided.  If no command given, general details about each command is given.",
			Aliases:     []string{"h"},
		},
	},
}

var findHelpPage = HelpPage{
	Title:   "locates certificates and associated resources, keys, requests etc",
	Aliases: nil,
	Description: "Find requires one or more file path arguments of where to search.\n" +
		"These can be directories or files.  Each path is searched for 'resources', certificates, keys, requests or revokation lists\n" +
		"Files may contain multiple resources but all must be encoded in the same manner.\n" +
		"Encodings supported are PEM and der (TODO: Add pk7 etc)\n" +
		"With no flags set, find will display ALL resources found in all the paths.  These may be filtered using the flags.\n\n" +
		"The prompts or fields output for each resource depend on both the 'resource-type' and 'fields' flags\n" +
		"When a single, specific resource type is set with the 'resource-type' flag, fields relating to that resource type are show.\n" +
		"When no type, or multiple types are specified, it shows generic fields for all resources.\n" +
		"Setting the 'fields' flag with specific fields will show values for all resources containing that field and empty for resource which do not.",
	Format: "find <path to directory or file> [<additional paths>...] <-flags>",
	Flags: []HelpPageFlag{
		{
			Name: "resource-type",
			Description: "Set to one or more resource types to show only those types.  Types should be comma deliminated\n" + "" +
				"Valid types are: certificate, privatekey, publickey, certificaterequest, revocationlist\n" + "" +
				"These may be abbreviated to: cert, key, puk, csr/request, crl respectively",
			Aliases: []string{"type"},
		},
		{
			Name: "fields",
			Description: "Optional, comma delimited list of field names to output.\n" +
				"When set with one or more fields, those fields are displayed in the result table, provided the resource contains such a property.\n" +
				"if the resource does not recognise the property it will show blank.\n" +
				"Field names may be optionally proceeded to a '+' sign which signals to add that field to the default fields already being displayed.\n" +
				"Without the +, The field will exclude the defaults and only show the fields listed.\n" +
				"By default, each field is displayed in a 20 character wide column.  This can be adjusted by ending the field name with square brackets, containing a width:\n" +
				"e.g. myfield[40]\nThis will display any values found in 'myfield' in a 40 character wide column\n" +
				"Note: the last column is not bound by any width.  To see a longer column (e.g. identity), simply use '-fields +identity' to add it as the last column\n" +
				"When -fields is empty, a default set of fields is used for each resource type or a generic field set of mixed types.",
			Aliases: []string{"fd"},
		},
		{
			Name: "query",
			Description: "Specify a query to limit to resources show.\n" +
				"Query should be one or more, comma delimited expresions specifying a field name and value.\n" +
				"e.g. myfield = 123, myotherfield = `hahaha`\n" +
				"The first operand (myfield) should be a field in the resource to search.\n" +
				"The second operand (123) may be a constant or another field.  contants may be integer numbers, double quoted strings or single quoted consts.\n" +
				"Operands must be seperated with an operator.  Valid operators are:\n" +
				"=, <, <=, >, >=, !=/NOT, contains, compare\n" +
				"contains should only be used on string fields\n" +
				"compare is like equals, but allows for fields such as serial-number which is too big to compare as an int.\n" +
				"Using compare, the second operand should be given as a single-quoted const: e.g. myfield compare '9949494949494949'",

			Aliases: []string{"qy"},
		},
		{
			Name: "recursive",
			Description: "When used, all sub directories of the given paths will also be searched.\n" +
				"By default only the given directories/files are searched and subdirectories ignored.",
			Aliases: []string{"r"},
		},
	},
}

var configHelpPage = HelpPage{
	Title:   "application configuration display and setting",
	Aliases: []string{"cfg"},
	Format:  "config [<path to a config file]",
	Description: "config control where the application locates and saves resources.  It is controlled by a '.config' file\n" +
		"The config file contains all the paths used to both find and save resources\n" +
		"This file is located by the following rules:\n" +
		"\tif the -config flag is set, this is used\n" +
		"\tif the environment variable CA_CONFIG is set, this is used\n" +
		"\tif the current directory contains a '.config' file, this is used\n" +
		"\tif the default $HOME/.pempal/.config file exists, this is used\n" +
		"\tfinally, if no file is located, the current working directory is assumed\n\n" +
		"When used with no parameters, the full, current configuration is displayed.\n" +
		"Parameters may be one or more of the config key names, when given, will limit the output to just those values.\n" +
		"To change a value, follow the config keyname directly with an equals sign, then followed by the value to set.\n" +
		"e.g. config root-path=./myrootpath",
}

var templateHelpPage = HelpPage{
	Title:   "display template names and contents",
	Aliases: []string{"tp"},
	Format:  "template <template name> [...<template-name] -flags",
	Description: "template manages the templates available to your current configuration.\n" +
		"Used with no parameters, it will list all the template names known to the application," +
		"including built-in/default templates and all of your own 'custom' templates.\n" +
		"One or more parameters may be given, each must be a valid template name.\n" +
		"When names are provided, the template(s) of that name are displayed in their formatted state.\n" +
		"Multiple templates may be merged into one using the merge flag.\n" +
		"This will merge left to right, so all the prompts of all templates are placed into a single template, which is then displayed.\n\n" +
		"To create a new template, use a new template name, followed directly by an equal sign and the property names, comma delimited:\n" +
		"e.g. template mynewtemple=\"subject.common-name: my server certificate, subject.organisation: Acme Ltd\"\n" +
		"The new template may also be read from the standard input, using a dash '-' in place of the template value.\n" +
		"e.g. template mynewtemplate=-",
	Flags: []HelpPageFlag{
		{
			Name: "merge",
			Description: "Merges all of the named templates into a new single template\n" +
				"The merge is performed from left to right, so any template to the right which shares a property name will overwrite the template to the left.\n" +
				"i.e. the last named template takes precedence.\n" +
				"The resulting template is written to the output.\n" +
				"to save the resulting template, pipe it into another template command:\n" +
				"e.g. template firsttemplate secondtemplate thirdtemplate -merge | template mynewtemplate=-",
			Aliases: nil,
		},
		{
			Name: "unformatted",
			Description: "When set, the template is displayed in its unformatted state, prior to any extending or template macros.\n" +
				"This is the state the template is stored, including its #tags and template macros.\n" +
				"Note: raw templates can not be merged.  The merge flag is ignored when unformatted is used.",
			Aliases: []string{"raw"},
		},
		{
			Name: "remove",
			Description: "When set all named templates will be deleted.\n" + "" +
				"Default templates can not be deleted.",
			Aliases: nil,
		},
	},
}

var unknownHelp = HelpPage{
	Title: "UnknownCurve command",
}

func HelpFor(command string) HelpPage {
	hp, ok := helpPages[command]
	if !ok {
		return unknownHelp
	}
	return hp
}
func init() {
	buf := bytes.NewBuffer(nil)
	out := utils.NewColumnOutput(buf)
	out.ColumnWidths = []int{10}
	for _, cmd := range listCommands() {
		s := []string{cmd}
		hp, ok := helpPages[cmd]
		if ok {
			s = append(s, hp.Title)
		}
		out.WriteSlice(s)
		out.WriteString("\n")
	}
	commandsHelpPage.Description = buf.String() + commandsHelpPage.Description
	helpPages[""] = commandsHelpPage
}

func listCommands() []string {
	var cmds []string
	for cmd := range commands.Commands {
		cmds = append(cmds, cmd)
	}
	sort.Strings(cmds)
	return cmds
}
