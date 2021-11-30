package commands

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
			fmt.Print(" password is too short\n")
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

func PromptInput(msg string, def string, minLen int, maxLen int) string {
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s", msg)
		if def != "" {
			fmt.Print(" [%s]")
		}
		fmt.Print(": ")
		l, _, err := in.ReadLine()
		if err != nil {
			log.Println(err)
			continue
		}
		// hit enter, take default
		if len(l) == 0 && def != "" {
			return def
		}
		if len(l) < minLen {
			fmt.Printf("\nEnter at least %d characters", minLen)
			continue
		}
		if maxLen > 0 && len(l) > maxLen {
			fmt.Printf("\nToo long, must not be longer than %d characters", maxLen)
			def = string(l[:maxLen])
			continue
		}
		return string(l)
	}
}

func PromptConfirm(msg string, def bool) bool {
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [", msg)
		if def {
			fmt.Print("Y,n")
		} else {
			fmt.Print("y,N")
		}
		fmt.Print("] ")
		l, _, err := in.ReadLine()
		if err != nil {
			log.Println(err)
			continue
		}
		// hit enter, take default
		if len(l) == 0 {
			return def
		}
		if strings.EqualFold(string(l), "n") {
			return false
		}
		if strings.EqualFold(string(l), "y") {
			return true
		}
		fmt.Printf("\nenter just 'y' or 'n'\n")
	}
}

func PromptInputNumber(max int) int {
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Select a number between 1 and %d. 0 to abort: ", max)
		s, _ := in.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		si, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("Not a number.  %v", err)
			continue
		}
		if si < 0 || si > max {
			continue
		}
		return si
	}
}

//
func PromptChooseList(msg string, choices []string) int {
	if msg != "" {
		fmt.Println(msg)
	}
	fmt.Println("Enter zero to abort")

	for i, line := range choices {
		fmt.Printf("%d) %s\n", i+1, line)
	}
	return PromptInputNumber(len(choices)) - 1
}
