package commands

import (
	"bytes"
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
	Recursive    bool
	ShowLocation bool
}

func (cmd *IssuersCommand) Description() string {
	lines := bytes.NewBufferString("lists the issuers available.\n")
	lines.WriteString("The list shows all the certificates which can be used to issue new certificates.\n")
	lines.WriteString("issuers links the available private keys to any certificates with valid properties for being an issuer.\n")
	lines.WriteString("Each certificate will be marked 'IsCA' true and have the KeyUsageCertSign permission\nas well as having a Public Key matching the available Private key\n")
	lines.WriteString("Only the certificate for which the private key is available are listed.\n")
	lines.WriteString("These are certificates which can be used by the caller (defined by their key access) for issuing new certificates.\n")
	return lines.String()
}

func (cmd *IssuersCommand) Flags(f *flag.FlagSet) {
	f.BoolVar(&cmd.Recursive, "r", false, "recursively search args (or keypath) for identity certificates and keys")
	f.BoolVar(&cmd.ShowLocation, "l", false, "show the file location of the certificate")
}

func (cmd *IssuersCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	args = GetKeyPath(args)
	if len(args) == 0 {
		return fmt.Errorf("must provide at least one location to search for issuers or set the %s environment variable with the path(s) to search.", ENV_KeyPath)
	}
	issuers := sortIssuers(issuers(ctx, args, cmd.Recursive, ""))

	//TODO, fix column sizing
	tw := tabwriter.NewWriter(out, 2, 1, 4, ' ', 0)

	for _, issuer := range issuers {
		k := issuer.Key()
		fmt.Fprintf(out, "%s\t%s\t%s", k.String(), k.Type(), issuer.String())
		if cmd.ShowLocation {
			fmt.Fprintf(out, "\t%s", issuer.Location())
		}
		fmt.Fprintln(out)
	}
	return tw.Flush()
}

func issuers(ctx context.Context, keypath []string, recursive bool, dn string) []keytracker.Identity {
	dn = strings.ToLower(dn)
	kt := keytracker.KeyTracker{ShowLogs: VerboseFlag, Recursive: recursive}
	idCh := kt.FindIdentities(ctx, keypath...)
	var found []keytracker.Identity
	for {
		select {
		case <-ctx.Done():
			return nil
		case id, ok := <-idCh:
			if !ok {
				return found
			}
			if !id.Certificate().IsCA || !id.Usage(x509.KeyUsageCertSign) {
				continue
			}
			if dn != "" && !strings.Contains(strings.ToLower(id.Certificate().Subject.String()), dn) {
				continue
			}
			found = append(found, id)
		}
	}
}

func sortIssuers(issuers []keytracker.Identity) []keytracker.Identity {
	sort.Slice(issuers, func(i, j int) bool {
		is := issuers[i].String()
		ij := issuers[i].String()
		ic := []string{is, ij}
		sort.Strings(ic)
		return ic[0] == is
	})
	return issuers
}
