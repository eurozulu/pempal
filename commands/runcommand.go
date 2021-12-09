package commands

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"pempal/fileformats"
	"pempal/keycache"
	"strings"
	"time"
)

const defaultOutputPermissions = os.FileMode(0640)
const defaultCommand = "help"

const PPName = "Pempal, the certificate assistant"
const PPVersion = "0.0.0"

const ENV_KEYPATH = "PEMPAL_KEYPATH"

var KeyPath = os.ExpandEnv(os.Getenv(ENV_KEYPATH))

var Verbose bool
var Version bool
var Script bool
var TimeRun bool
var OutFileName string
var HelpFlag bool
var Format string

var formatWriter fileformats.FormatWriter

func Flags(f *flag.FlagSet) {
	f.BoolVar(&Verbose, "verbose", false, "Display all logging whilst searching for pems")
	f.BoolVar(&Version, "v", false, "same as verbose")
	f.BoolVar(&Version, "version", false, "outputs the version of the application")
	f.BoolVar(&TimeRun, "t", false, "Times how long the command takes to execute")
	f.StringVar(&OutFileName, "out", "", "Specify a filename to write output into. Defaults to stdout")
	f.BoolVar(&HelpFlag, "help", false, "Display help")
	f.StringVar(&Format, "format", "yaml", "The output format. should be yaml, pem or der")
}

// Runs the command given as the first argument
func RunCommand(args ...string) error {
	// establish the command.

	// If no command (no args given, or first arg a flag) insert defaultcommandd as first argument
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		args = append([]string{defaultCommand}, args...)
	}
	cmdArg := args[0]

	// remap command if its an alias
	if als, ok := Aliases[cmdArg]; ok {
		cmdArg = als
	}
	command, ok := Commands[cmdArg]
	if !ok {
		return fmt.Errorf("%s is an unknown command\n", cmdArg)
	}
	// add flags, both general and command specific
	Flags(flag.CommandLine)
	command.Flags(flag.CommandLine)

	// parse the command line args, trimmed of the leading command
	if err := flag.CommandLine.Parse(args[1:]); err != nil {
		return fmt.Errorf("Invalid command: %w", err)
	}
	args = flag.Args()

	// post parse, check the general flags

	if Version {
		showVersion()
		return nil
	}

	if HelpFlag {
		command = Commands["help"]
		// reinsert first argument for command specific help
		if cmdArg != "help" {
			args = append([]string{cmdArg}, args...)
		}
	}

	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	if keyCmd, ok := command.(SigningCommand); ok {
		keyCmd.SetKeys(keycache.NewKeyCache(ctx, GetKeyPath([]string{})...))
	}

	if fr, err := fileformats.NewFormatWriter(Format); err != nil {
		return err
	} else {
		formatWriter = fr
	}

	// establish the output stream
	out := os.Stdout
	if OutFileName != "" {
		f, err := os.OpenFile(OutFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, defaultOutputPermissions)
		if err != nil {
			return fmt.Errorf("Failed to open output file %s  %w", OutFileName, err)
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

	if TimeRun {
		defer func(s time.Time) {
			fmt.Printf("took: %v\n", time.Now().Sub(s))
		}(time.Now())
	}

	// main loop starts

	// monitor OS signal to trigger ctx cancel to all on-going routines when killed
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Kill, os.Interrupt)
	// signal channel to indicate command has completed
	done := make(chan struct{})

	go func() {
		defer close(done)
		if err := command.Run(ctx, out, args...); err != nil {
			log.Println(err)
		}
	}()
	for {
		select {
		case <-sig:
			return fmt.Errorf("cancelled")
		case <-done:
			return nil
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

func showVersion() {
	fmt.Printf("\n%s, Version: %s\n", PPName, PPVersion)
}

func GetKeyPath(p []string) []string {
	if len(p) == 0 {
		p = []string{os.ExpandEnv("$PWD")}
	}
	if KeyPath == "" {
		return p
	}
	return append(p, strings.Split(KeyPath, ":")...)
}
