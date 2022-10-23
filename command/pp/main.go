package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const ENV_KEY_PATH = "KEY_PATH"
const ENV_PP_PATH = "PP_PATH"

const FLAG_OUTPUT = "output"

var Commands = map[string]NewCommandFnc{
	"find":    newFindCommand,
	"genkey":  newGenKeyCommand,
	"request": newRequestCommand,
	"issue":   newIssueCommand,
	"revoke":  newRevokeCommand,
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "no command given")
		return
	}
	nc, ok := Commands[strings.ToLower(os.Args[1])]
	if !ok {
		fmt.Fprintf(os.Stderr, "%s is an unknown command\n", os.Args[1])
		return
	}
	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}
	out := os.Stdout
	var op string
	op, args = readOutPath(args)
	if op != "" {
		f, err := os.OpenFile(op, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func(out io.WriteCloser) {
			if err := out.Close(); err != nil {
				fmt.Println(err)
			}
		}(f)
		out = f
	}
	cmd := nc()
	params, err := ApplyFlags(cmd, args...)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := cmd.Run(out, params...); err != nil {
		fmt.Println(err)
	}
}

func readOutPath(args []string) (string, []string) {
	if len(os.Args) < 3 {
		return "", args
	}
	s := &struct {
		Out string `flag:"out"`
	}{}
	var err error
	args, err = ApplyFlags(s, args...)
	if err != nil {
		fmt.Println(err)
		return "", args
	}
	return s.Out, args
}
