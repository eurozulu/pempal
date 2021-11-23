package cmd

import (
	"context"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"pempal/pemreader"
	"pempal/pemwriter"
)

var Verbose bool

type ShowCommand struct {
	format string
}

func (cmd *ShowCommand) Description() string {
	return "outputs x509 resource properties in plain text or formatted as pem or der"
}

func (cmd *ShowCommand) Flags(f *flag.FlagSet) {
	flag.StringVar(&cmd.format, "format", "", "defines the output format of the items. One of: pem, der.  Unstated generates a template")
}

// ShowCommand takes one or more args as pem locations, and outputs them in a given format
func (cmd *ShowCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("must provide at least one location to display as a template.")
	}

	pr := pemreader.PemReader{
		Verbose:           Verbose,
		AddLocationHeader: cmd.format == "", // Add location if output is template
		Recursive:         true,
	}
	pw, err := pemwriter.NewPemWriter(cmd.format, out)
	if err != nil {
		return err
	}
	// Scan each arg in sequence to maintain order
	for _, arg := range args {
		err := cmd.formatPems(ctx, pr.Find(ctx, arg), pw)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cmd ShowCommand) formatPems(ctx context.Context, pemIn <-chan *pem.Block, pemOut pemwriter.PemWriter) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case b, ok := <-pemIn:
			if !ok {
				return nil
			}
			if err := pemOut.Write(b); err != nil {
				return err
			}
		}
	}
}

func handleError(err error) bool {
	if err == nil {
		return true
	}
	if Verbose {
		fmt.Println(err)
	}
	return false
}
