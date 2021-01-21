package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func PromptConfirm(msg string, def bool) bool {
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [", msg)
		if def {
			fmt.Print("y,N")
		} else {
			fmt.Print("Y,n")
		}
		fmt.Print("] ")
		r, _, err := in.ReadRune()
		if err != nil {
			log.Println(err)
			r = 'n'
		}
		if strings.EqualFold(string(r), "n") {
			return false
		}
		if strings.EqualFold(string(r), "y") {
			return true
		}
		fmt.Printf("\nenter just 'y' or 'n'\n")
	}
}
