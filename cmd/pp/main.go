package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"pempal/cmd"
	"pempal/pemreader"
	"strings"
	"time"
)

func main() {
	started := time.Now()
	var help bool
	var timeRun bool
	var pemTypes string
	var outFile string

	if len(os.Args) < 2 {
		log.Println("must provide a command as the first argument")
		return
	}

	fs := flag.CommandLine
	fs.BoolVar(&cmd.Verbose, "verbose", false, "Display all logging whilst searching for pemreader")
	fs.BoolVar(&timeRun, "t", false, "Times how long the command takes to execute")
	fs.BoolVar(&help, "help", false, "Display help")
	fs.StringVar(&outFile, "out", "", "Specify a filename to write the output to. Defaults to stdout")
	fs.StringVar(&pemTypes, "types", "", "A comma delimited list of PEM types.")

	command, ok := cmd.Commands[os.Args[1]]
	if !ok {
		log.Fatalf("%s unknown command", os.Args[1])
	}
	command.Flags(fs)
	if err := flag.CommandLine.Parse(os.Args[2:]); err != nil {
		log.Fatalln(err)
	}
	args := flag.Args()

	// if '-' in the args, insert the std in at that point
	if i := stringIndex("-", args); i >= 0 {
		as := args[:i]
		var ae []string
		if i+1 < len(args) {
			ae = args[i+1:]
		}
		args = append(as, readInput()...)
		args = append(args, ae...)
	}

	if timeRun {
		defer func(s time.Time) {
			fmt.Printf("took: %v\n", time.Now().Sub(s))
		}(started)
	}

	if pemTypes != "" {
		pemreader.AddPemTypes(strings.Split(pemTypes, ","))
	}

	// establish the output stream
	out := os.Stdout
	if outFile != "" {
		f, err := os.OpenFile(outFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
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
