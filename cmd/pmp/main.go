package main

import (
	"log"
	"os"

	"github.com/eurozulu/mainline"
)

func main() {
	cmd := mainline.Commands{
		"find": FindCommand.Find,
		"ls":   FindCommand.Find,
		"view": ViewCommand.ViewItems,

		"key":     KeyCommand.Key,
		"request": RequestCommand.Request,
		"issue":   IssueCommand.Issue,

		"help": mainline.HelpCommand.Help,
	}
	if err := cmd.Run(os.Args...); err != nil {
		log.Fatalln(err)
	}
}
