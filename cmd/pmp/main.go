package main

//go:generate echo "$GOPACKAGE $GOFILE"

import (
	"github.com/eurozulu/commandgo"
	"log"
	"os"
)

func main() {
	cmd := commandgo.Commands{
		"find": FindCommand.Find,
		"ls":   FindCommand.Find,
		"view": ViewCommand.ViewItems,

		"key":     KeyCommand.Key,
		"request": RequestCommand.Request,
		"issue":   IssueCommand.Issue,

		"help": commandgo.HelpCommand.Help,
	}

	// Out sets an output filename. Defaults to the standard output
	commandgo.AddFlag(&Out, "out", "o")

	// Encode sets the output format. valid values are 'pem', 'der' or 'p12'
	commandgo.AddFlag(&Encode, "encode", "e")

	if err := cmd.Run(os.Args...); err != nil {
		log.Fatalln(err)
	}
}
