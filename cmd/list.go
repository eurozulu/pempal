package cmd

import (
	"bytes"
	"context"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"pempal/pemreader"
	"pempal/templates"
	"strings"
	"text/tabwriter"
)

type ListCommand struct {
	quiet       bool
	recursive   bool
	showHeaders bool
	queryString string
	query       ListQuery
}

func (cmd *ListCommand) Description() string {
	return "finds x509 resources in the given path names.  Can filter using 'query' to find specific resources"
}

func (cmd *ListCommand) Flags(f *flag.FlagSet) {
	f.BoolVar(&cmd.quiet, "q", false, "output only the file locations of found resources.")
	f.BoolVar(&cmd.recursive, "r", true, "searches sub directories.")
	f.BoolVar(&cmd.showHeaders, "h", true, "shows any PEM header values found in the block.")
	f.StringVar(&cmd.queryString, "query", "", "comma delimited list of key names, with optional regex expressions to match to searched resources")
}

func (cmd *ListCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("must provide at least one location to search")
	}

	if cmd.queryString != "" {
		q, err := ParseQuery(cmd.queryString)
		if err != nil {
			return err
		}
		cmd.query = q
	}

	pr := pemreader.PemScanner{
		Verbose:           Verbose,
		AddLocationHeader: true,
		Recursive:         cmd.recursive,
		PemTypes:          pemreader.PemTypes,
	}

	tw := tabwriter.NewWriter(out, 4, 8, 1, '\t', 0)
	// Scan each arg in sequence to maintain order
	for _, arg := range args {

		err := cmd.listPems(ctx, pr.Find(ctx, arg), tw)
		if err != nil {
			return err
		}
		if err = tw.Flush(); err != nil {
			return err
		}
	}
	return nil

}

func (cmd ListCommand) listPems(ctx context.Context, pemIn <-chan *pem.Block, out io.Writer) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case b, ok := <-pemIn:
			if !ok {
				return nil
			}

			var cols []string
			if cmd.query != nil {
				// query is present, gather values from matching block or ignore unmathed block
				qs, ok := cmd.queryBlock(b)
				if !ok {
					continue
				}
				cols = append(cols, qs...)
			}

			lineCols := cmd.formatPem(b)
			// insert query values before last linecol (location)
			cols = append(cols, lineCols[len(lineCols)-1])
			lineCols = append(lineCols[:len(lineCols)-1], cols...)

			if _, err := fmt.Fprintln(out, strings.Join(lineCols, "\t")); err != nil {
				return err
			}
		}
	}
}
func (cmd ListCommand) formatPem(b *pem.Block) []string {
	var loc string
	headers := bytes.NewBuffer(nil)
	for k, v := range b.Headers {
		if strings.EqualFold(k, pemreader.LocationHeaderKey) {
			loc = v
			if !cmd.showHeaders {
				break
			}
			continue
		}
		if !cmd.showHeaders {
			continue
		}
		if headers.Len() > 0 {
			headers.WriteString(", ")
		}
		headers.WriteString(k)
		headers.WriteString(" = ")
		headers.WriteString(v)
	}
	if headers.Len() > 0 {
		headers.WriteRune('\t')
	}
	return []string{b.Type, headers.String(), loc}
}

func (cmd ListCommand) queryBlock(b *pem.Block) ([]string, bool) {
	t, err := templates.ParseBlock(b)
	if !handleError(err) {
		return nil, false
	}
	if !cmd.query.Match(t) {
		return nil, false
	}
	return cmd.query.Values(t), true
}
