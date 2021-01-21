package main

import (
	"github.com/eurozulu/mainline"
	"log"
	"os"
)

func main() {
	cmd := mainline.Commands{
		"find":    FindCommand.Find,
		"ls":      FindCommand.Find,
		"view":    ViewCommand.View,
		"vw":      ViewCommand.View,
		"makekey": MakeKeyCommand.MakeKey,
		"mk":      MakeKeyCommand.MakeKey,
		"help":    mainline.HelpCommand.Help,
	}
	if err := cmd.Run(os.Args...); err != nil {
		log.Fatalln(err)
	}
}
