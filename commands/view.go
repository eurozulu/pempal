package commands

import (
	"bytes"
	"context"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"pempal/pemreader"
	"pempal/pemwriter"
)

type ViewCommand struct {
	format string
}

func (cmd *ViewCommand) Description() string {
	lines := bytes.NewBufferString("outputs x509 resources in plain text or formatted as pem or der\n")
	lines.WriteString("Requires one or more locations as parameters.  Location can be a directory or single file.\n")
	lines.WriteString("By default output will be in text (Yaml) format.\n")
	lines.WriteString("Use the -format flag to alter the output format to pem or der.\n")
	return lines.String()
}

func (cmd *ViewCommand) Flags(f *flag.FlagSet) {
	f.StringVar(&cmd.format, "format", "", "defines the output format of the items. One of: pem, der.  When unstated generates a yaml template of the resource")
}

// TODO: Fix bug with ECDSA keys (failed to parse EC private key: asn1: structure error: length too large)
// ViewCommand takes one or more args as pem locations, and outputs them in a given format
func (cmd *ViewCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("must provide at least one location to display as a template.")
	}

	pr := pemreader.PemScanner{
		Verbose:           VerboseFlag,
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

func (cmd ViewCommand) formatPems(ctx context.Context, pemIn <-chan *pem.Block, pemOut pemwriter.PemWriter) error {
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
	if VerboseFlag {
		fmt.Println(err)
	}
	return false
}
