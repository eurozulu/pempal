package commands

import (
	"bufio"
	"os"
	"strconv"
)

var prompertyPrompts = map[string]Prompt[any]{
	"public-key":          nil,
	"subject.common-name": newStringPrompt(),
	"missing signature":   nil,
	"missing subject":     nil,
	"missing issuer":      nil,
	"not-after":           nil,
}

type Prompt[V any] interface {
	Request(prompt string) V
}

type stringPrompt struct {
	defaultValue string
}

func (sp stringPrompt) Request(prompt string) string {
	writePrompt(prompt)
	return readStdInputLine()
}
func newStringPrompt() Prompt[string] {
	return &stringPrompt{}
}

func (sp stringPrompt) Requestx(prompt string) string {
	writePrompt(prompt)
	s := readStdInputLine()
	if s == "" {
		return sp.defaultValue
	}
	return s
}

func writePrompt(p string) {
	if p != "" {
		os.Stdout.WriteString(p)
		os.Stdout.WriteString(" ")
	}
}

func readStdInputLine() string {
	buf := bufio.NewReader(os.Stdin)
	s, _ := buf.ReadString('\n')
	return s
}

func readStdInputNumber() int {
	var by []byte
	for {
		b := readStdInputChoice('0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '\n')
		if b == '\n' {
			break
		}
		os.Stdout.Write([]byte{b})
		by = append(by, b)
	}
	i, _ := strconv.Atoi(string(by))
	return i
}

func readStdInputChar() byte {
	buf := bufio.NewReader(os.Stdin)
	s, _ := buf.ReadByte()
	return s
}
func readStdInputChoice(b ...byte) byte {
	for {
		i := readStdInputChar()
		for _, by := range b {
			if by == i {
				return by
			}
		}
	}
}

func NewPrompt() Prompt {
	return &stringPrompt{}
}
