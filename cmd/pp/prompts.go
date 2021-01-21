package main

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"log"
	"os"
	"strings"
	"syscall"
)

// PromptPassword requests a single line of text from the stdin without echoing the string
func PromptPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	by, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	return string(by), nil
}

func PromptCreatePassword(prompt string, minLength int) (string, error) {
	for {
		if minLength > 0 {
			fmt.Printf("(minimum %d charactors)  ", minLength)
		}
		pass, err := PromptPassword(prompt)
		if err != nil {
			return "", err
		}
		if len(pass) < minLength {
			fmt.Print("password is too short\n")
			continue
		}
		fmt.Println()

		fmt.Printf("\nReenter password to confirm: ")
		pass2, err := PromptPassword("")
		if err != nil {
			return "", err
		}
		fmt.Println()
		if pass != pass2 {
			fmt.Println("passwords do not match")
			continue
		}
		return pass, nil
	}
}

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
