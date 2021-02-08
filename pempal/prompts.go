package pempal

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func CertPath() []string {
	var ps []string
	p, err := filepath.Abs(".")
	if err != nil {
		log.Println(err)
	}
	if p != "" {
		ps = append(ps, p)
	}
	sp, ok := os.LookupEnv("PEMPAL_CERTPATH")
	if ok {
		ps = append(ps, strings.Split(sp, "|")...)
	}
	return ps
}

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

func PromptInputNumber(max int) int {
	in := bufio.NewReader(os.Stdin)
	for {
		s, _ := in.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return -1
		}
		si, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("Enter a number.  %v", err)
			continue
		}
		if si < 0 || si > max {
			fmt.Println("Enter a number between 0 and %d", max)
			continue
		}

		return si
	}
}
