package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"pempal/commands"
	"strings"
	"time"
)

func main() {
	started := time.Now()
	args := os.Args[1:]

	// establish the command
	// If no args given, or first arg a flag, insert empty string to force help
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		args = append([]string{""}, args...)
	}
	firstArg := args[0]

	// remap if its an alias command
	if als, ok := commands.Aliases[firstArg]; ok {
		firstArg = als
	}
	command, ok := commands.Commands[firstArg]
	if !ok {
		log.Fatalf("%s is an unknown command\n", firstArg)
	}
	// add flags, both general and command specific
	commands.FlagsMain(flag.CommandLine)
	command.Flags(flag.CommandLine)

	// parse the command line args, (again) trimmed of the leading command
	if err := flag.CommandLine.Parse(args[1:]); err != nil {
		log.Fatalln(err)
	}
	args = flag.Args()
	if commands.HelpFlag {
		command = commands.Commands[""]
		// Add on first command for command specific help
		if firstArg != "" {
			args = []string{firstArg}
		}
	} else {

	}

	if commands.TimeRunFlag {
		defer func(s time.Time) {
			fmt.Printf("took: %v\n", time.Now().Sub(s))
		}(started)
	}

	// establish the output stream
	out := os.Stdout
	if commands.OutFileFlag != "" {
		f, err := os.OpenFile(commands.OutFileFlag, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
		if err != nil {
			log.Fatalln(err)
		}
		defer func(f *os.File) {
			if err := f.Close(); err != nil {
				log.Println(err)
			}
		}(f)
		out = f
	}

	// if '-' in the args, insert the std in at that point
	// TODO: Change pipe in to accept PEM blocks :  Allowing one PP command to stream into another.
	if i := stringIndex("-", args); i >= 0 {
		as := args[:i]
		var ae []string
		if i+1 < len(args) {
			ae = args[i+1:]
		}
		args = append(as, readInput()...)
		args = append(args, ae...)
	}

	// monitor OS signal so trigger ctx cancel to all on-going routines when killed or interrupted
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Kill, os.Interrupt)
	done := make(chan struct{})

	go func() {
		defer close(done)
		if err := command.Run(ctx, out, args...); err != nil {
			log.Fatalln(err)
		}
	}()
	for {
		select {
		case <-sig:
			return
		case <-done:
			return
		}
	}
}

func readInput() []string {
	var lines []string
	scn := bufio.NewScanner(os.Stdin)
	for scn.Scan() {
		line := strings.TrimSpace(scn.Text())
		if line == "" {
			break
		}
		lines = append(lines, line)
	}
	return lines
}

func stringIndex(s string, ss []string) int {
	for i, sz := range ss {
		if sz == s {
			return i
		}
	}
	return -1
}
