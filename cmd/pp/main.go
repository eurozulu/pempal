package main

import (
	"github.com/eurozulu/mainline"
	"log"
	"os"
)

func main() {
	cmd := mainline.Commands{
		"find, fd, ls": &FindCommand{},
		"view, vw": &ViewCommand{},
	}
	if err := cmd.Run(os.Args...); err != nil {
		log.Fatalln(err)
	}
}
