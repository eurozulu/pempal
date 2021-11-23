package cmd

import (
	"context"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"pempal/keytracker"
	"sort"
	"strings"
	"text/tabwriter"
)

type IssuersCommand struct {
	listCerts bool
	showHash  bool
}

func (cmd *IssuersCommand) Description() string {
	return fmt.Sprintf("lists the issuers available.")
}

func (cmd *IssuersCommand) Flags(f *flag.FlagSet) {
}

func (cmd *IssuersCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	if KeyPath != "" {
		args = append(args, strings.Split(KeyPath, ":")...)
	}
	if len(args) == 0 {
		return fmt.Errorf("must provide at least one location to search for issuers or set the %s environment variable with the path(s) to search.", ENV_KeyPath)
	}

	issuers := sortIssuers(collectIssuers(ctx, "", args))

	//TODO, fix column sizing
	tw := tabwriter.NewWriter(out, 2, 1, 4, ' ', 0)

	for _, issuer := range issuers {
		k := issuer.Key()
		certs := sortCerts(issuer.Certificates(0, 0))
		fmt.Fprintf(out, "%s\t%s\t%s", k.String(), k.Type(), k.Location())
		for _, c := range certs {
			fmt.Fprintf(out, "\t%s\n", c.Subject.String())
		}
	}
	return tw.Flush()
}

func collectIssuers(ctx context.Context, dn string, keypath []string) []keytracker.Identity {
	kt := &keytracker.KeyTracker{ShowLogs: Verbose}
	issuers := kt.Issuers(ctx, dn, keypath)

	var found []keytracker.Identity
	for {
		select {
		case <-ctx.Done():
			return nil

		case id, ok := <-issuers:
			if !ok {
				return found
			}
			found = append(found, id)
		}
	}
}

func sortIssuers(ids []keytracker.Identity) []keytracker.Identity {
	sort.Slice(ids, func(i, j int) bool {
		is := ids[i].String()
		ij := ids[i].String()
		ic := []string{is, ij}
		sort.Strings(ic)
		return ic[0] == is
	})
	return ids
}

func sortCerts(certs []*x509.Certificate) []*x509.Certificate {
	sort.Slice(certs, func(i, j int) bool {
		is := certs[i].Subject.String()
		ij := certs[i].Subject.String()
		ic := []string{is, ij}
		sort.Strings(ic)
		return ic[0] == is
	})
	return certs
}
