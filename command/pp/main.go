package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"pempal/command"
	"time"
)

var beginTime time.Time

const outFilePermissions = 0640

func main() {
	beginTime = time.Now()
	defer func() {
		fmt.Printf("took %d Milliseconds\n", time.Now().Sub(beginTime).Milliseconds())
	}()

	// First parse any global flags from the command line args
	var gFlags globalFlags
	args, err := command.ApplyArguments(&gFlags, os.Args[1:])
	if err != nil {
		writeError(err)
		return
	}

	out := os.Stdout
	if gFlags.OutputPath != "" {
		f, err := os.OpenFile(gFlags.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, outFilePermissions)
		if err != nil {
			writeError(err)
			return
		}
		defer func(out io.WriteCloser) {
			if err := out.Close(); err != nil {
				writeError(err)
			}
		}(f)
		out = f
	}

	// parse remaining args into command and argument
	cmd, args, err := command.NewCommand(args.Args()...)
	if err != nil {
		writeError(err)
		return
	}

	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	sig := make(chan os.Signal)
	done := make(chan bool)
	go func() {
		defer close(done)
		if err := cmd.Run(ctx, args, out); err != nil {
			writeError(err)
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

func writeError(err error) {
	_, e := fmt.Fprintf(os.Stderr, "%v\n", err)
	if e != nil {
		log.Println(err)
	}
}

type globalFlags struct {
	OutputPath string `flag:"out,o"`
}
