package main

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"log"
	"os"
	"strconv"
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

func PromptInput(msg string, def string) string {
	in := bufio.NewReader(os.Stdin)
	if msg != "" {
		fmt.Print(msg)
		if def != "" {
			fmt.Printf(" [%s] ", def)
		}
	}
	s, _ := in.ReadString('\n')
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}

func PromptChoice(msg string, choices []string) int {
	in := bufio.NewReader(os.Stdin)
	if msg != "" {
		fmt.Println(msg)
	}
	for i, ch := range choices {
		fmt.Printf("%02d) %s\n", i + 1, ch)
	}
	fmt.Printf("%02d) %s\n", 0, "cancel")
	for {
		fmt.Printf("Select 1 - %02d or 0 to abort: ", len(choices))
		s, err := in.ReadString('\n')
		if err != nil {
			log.Println(err)
			return -1
		}
		i, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("Not a number  %v\n", err)
		}
		if i >= 0 && i <= len(choices) {
			return i - 1
		}
	}
}