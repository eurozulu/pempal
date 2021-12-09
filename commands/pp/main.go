package main

import (
	"log"
	"os"
	"pempal/commands"
)

func main() {
	if err := commands.RunCommand(os.Args[1:]...); err != nil {
		log.Println(err)
	}
}
