package commands

import (
	"bytes"
	"context"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"pempal/fileformats"
	"pempal/pemresources"
	"pempal/templates"
	"strings"
	"text/tabwriter"
)

// FindCommand locates x509 resources
type FindCommand struct {
	Recursive   bool
	ShowHeaders bool
	Types       string
	QueryString string
	Query       FindQuery
}

func (cmd *FindCommand) Description() string {
	fn := path.Base(os.Args[0])
	line := bytes.NewBufferString("finds certificates and related resources in the given path names.  Can filter search using most properties of those resources.\n")
	line.WriteString("Requires at least one argument of the location to search.  Multiple locations can be given as space delimited argument.\n")
	line.WriteString("Locations can be a directory path, or a single file.  Use -r recursive flag to search sub directories\n")
	line.WriteString(fmt.Sprintf("e.g. %s find -r /etc/ssl /etc/ssh ~/.ssh\n", fn))
	line.WriteString("Will list all the pem resources in those three directories\n")
	line.WriteString("\n")
	line.WriteString("By default find outputs basic information, pem type and location.\n")
	line.WriteString("Using the -query flag enables additional properties to be added to the output list as well as optionally filtering the list.\n")
	line.WriteString("\n")
	line.WriteString("the query flag uses the format:  <property name>[=<value or expression>] [,<property name>[=<value or expression>]...\n")
	line.WriteString("Simply adding a property name will output that property value in the list, assuming the resource has such a property.\n")
	line.WriteString("If the resource has no such property a blank is output in that property, for that resource.\n")
	line.WriteString(fmt.Sprintf("e.g. %s find -r -query \"Subject.CommonName\" /etc/ssl\nwill output the CommonName of certificates and CSRs found in /etc/ssl\n", fn))
	line.WriteString("\n")
	line.WriteString("To limit the output list to resources only containing that property, place an expression after the name.\n")
	line.WriteString(fmt.Sprintf("e.g. %s find . -query \"isCA=true\" will only list certificates (IsCA is unique to certificates) that are Certificate Authorities\n", fn))
	line.WriteString("The expression or value must be precceded with an '=' and may be simple text or a regular expression.\n")
	line.WriteString("Multiple properties can be stated, delimited with a comma\n")
	line.WriteString(fmt.Sprintf("e.g. %s find . -query \"Subject.CommonName=^Dev.*Server.*, IsCA=false\" will only list non CA certificates with a common name matching that regex.\n", fn))
	line.WriteString("To limit output to only resources with a named property with any value, use the '*' wildcard\n")
	line.WriteString(fmt.Sprintf("e.g. %s find . -query \"PublicKey=*\" will only list resources with an identifyable public key (Certificates, csrs, unencrypted or 'linked' keys\n", fn))
	line.WriteString("\n")
	line.WriteString("Most properties within resources can be queries.  Use 'template' or 'view' on your existing properties to see the names they contain.\n")
	line.WriteString("Most binary fields (PublicKey, Signature etc) have a 'xxxHash' field to use for comparison. This is a SHA1 hash of the bytes. e.g. PublicKeyHash will identify any resource with the same public key\n")
	line.WriteString("In addition to the properties in the resource, some additional properties can be queries:\n")
	line.WriteString("location: this is the full filepath of the resource, including any index.\n\t" +
		"An index is added to file locations containing more than one resource. Following the filepath, a colon and a number will indicate the position within the file.\n\t")
	line.WriteString("Files containing a single resource have the index ommited\n")
	line.WriteString("pem_type:  This property contains the PEM type of the resource.  'CERTIFICATE', 'PRIVATE KEY' etc")
	line.WriteString("A few properties are shared by all resources.  PublicKey is a one of them.  Using the PublicKeyHash query will list all resources using that key, including the private key, any certificates or csrs its signed etc.\n")
	line.WriteString("\n")
	return line.String()
}

func (cmd *FindCommand) Flags(f *flag.FlagSet) {
	f.BoolVar(&cmd.Recursive, "r", false, "searches sub directories.")
	f.BoolVar(&cmd.ShowHeaders, "h", false, "shows any PEM header values found in the resource.")
	f.StringVar(&cmd.Types, "types", "", "comma delimited list of pem types. CERTIFICATE, RSA PRIVATE KEY etc, limits results to just the listed types")
	f.StringVar(&cmd.QueryString, "query", "", "comma delimited list of key names, with optional regex expressions to match to searched resources")
}

func (cmd *FindCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("must provide at least one location to search")
	}

	if cmd.QueryString != "" {
		q, err := ParseQuery(cmd.QueryString)
		if err != nil {
			return err
		}
		cmd.Query = q
	}
	ps := pemresources.PemScanner{
		Recursive:         cmd.Recursive,
		Verbose:           Verbose,
		Reader:            fileformats.NewFormatReader(),
		TypeFilter:        cmd.typeFilterMap(),
		AddLocationHeader: true,
	}
	pemCh := ps.Scan(ctx, args...)
	tw := tabwriter.NewWriter(out, 4, 8, 1, '\t', 0)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case b, ok := <-pemCh:
			if !ok {
				return nil
			}

			lineCols := cmd.formatPem(b)
			if cmd.Query != nil {
				// query is present, gather values from matching block or ignore unmathed block
				qs, ok := cmd.queryBlock(b)
				if !ok {
					continue
				}
				// insert query values before last linecol (location)
				qs = append(qs, lineCols[len(lineCols)-1]) // query with location on the end
				lineCols = append(lineCols[:len(lineCols)-1], qs...)
			}
			if _, err := fmt.Fprintln(tw, strings.Join(lineCols, "\t")); err != nil {
				return err
			}
		}
	}
}

func (cmd FindCommand) formatPem(b *pem.Block) []string {
	fields := []string{b.Type}

	// gather location header and, optionally the rest of them.
	if cmd.ShowHeaders {
		headers := bytes.NewBuffer(nil)
		for k, v := range b.Headers {
			if headers.Len() > 0 {
				headers.WriteString(", ")
			}
			headers.WriteString(k)
			headers.WriteString(" = ")
			headers.WriteString(v)
		}
		fields = append(fields, headers.String())
	}
	fields = append(fields, b.Headers[pemresources.LocationHeaderKey])
	return fields
}

func (cmd FindCommand) queryBlock(b *pem.Block) ([]string, bool) {
	t, err := templates.BlockToTemplate(b)
	if !handleError(err) {
		return nil, false
	}
	if !cmd.Query.Match(t) {
		return nil, false
	}
	return cmd.Query.Values(t), true
}

func (cmd FindCommand) typeFilterMap() map[string]bool {
	if cmd.Types == "" {
		return nil
	}
	types := strings.Split(cmd.Types, ",")
	m := map[string]bool{}
	for _, t := range types {
		t = strings.ToUpper(t)
		if fileformats.PemTypesCertificate[t] {
			m = fileformats.CombineMaps(m, fileformats.PemTypesCertificate)
			continue
		}
		if fileformats.PemTypesPrivateKey[t] {
			m = fileformats.CombineMaps(m, fileformats.PemTypesPrivateKey)
			continue
		}
		if fileformats.PemTypesPublicKey[t] {
			m = fileformats.CombineMaps(m, fileformats.PemTypesPublicKey)
			continue
		}
		if fileformats.PemTypesCertificateRequest[t] {
			m = fileformats.CombineMaps(m, fileformats.PemTypesCertificateRequest)
			continue
		}
		m[t] = true
	}
	return m
}
