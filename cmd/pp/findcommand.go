package main

import (
	"context"
	"fmt"
	"github.com/eurozulu/pempal/filescan"
	"github.com/eurozulu/pempal/templates"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
	"text/tabwriter"
)

// FindCommand finds x509 resources based on given criteria
// Can be used to locate, group and bundle resources
type FindCommand struct {
	// Recursive, when true, searches sub containers found in any location being searched.
	// By default false, so only the named resource/containers are searched.
	// When true, searches sub containers such as sub-directories
	Recursive bool `flag:"recursive,r"`

	// Verbose adds additional fields about the resource to the output.
	// When true, implies -table to format the output in columns
	Verbose bool `flag:"verbose,v"`
	// VeryVerbose adds to verbose output by listing each file visited
	VeryVerbose bool `flag:"vv"`

	// Chain attempts to locate the full certificate chain of each certificate found.
	// Any certificate not self signed will invoke a secondary find for its issuer.
	// These continue until the self signed root is found.
	Chain bool `flag:"chain,c"`

	// Keys attempts to locate the related keys to any found resource.
	// Any signed resource (cert, csr, crl) starts a find for its public key file.
	Keys bool `flag:"keys,k"`

	// Query is a free text search of a resource.
	// If the given string appears in the resource (As it would appear on screen)
	// It will be included in the results, otherwise it will be omitted.
	// e.g. -t cer -q "Root CA"  Will find any certificate with "Root CA" in any of its fields.
	// Query is case sensitive unless -insensitive is set
	Query string `flag:"query,q"`

	// Insensative effects the -query flag result, making it case insensitive when specified.
	Insensitive bool `flag:"insensitive,i"`

	// Type specifies the type of resource to find.  Can be a comma delimited list of one or more of th efollowing:
	// 			cer		finds only certificates
	// 			csr		Finds only certiifcate requests
	// 			crl		Finds only certiifcate revokation lists
	// 			puk		Finds only public keys
	// 			prk		Finds only private keys
	// When not specified, finds All types.
	Type []string `flag:"type,t"`

	Password string `flag:"pass,p"`
}

// Find locates files (Certificates, CSRs, CRLs, Private & Public keys) based on query criteria.
// With no flags, it simply lists the resource location, depending on its container:
// A file in a directory has the full path listed if it contains a single item.
// files containing more than one item are listed with a sub index:
// ./thisdir/certs/servertcerts.p12#3 Indicating it is the 3rd resource in this file.
// The path arguments may be a path to a specific file or a container such as a directory or pkcs container file.
func (fc FindCommand) Find(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Find: must provide at least one path")
	}

	var err error
	fc.Type, err = formatTypeNames(fc.Type)
	if err != nil {
		return err
	}

	if fc.VeryVerbose {
		fc.Verbose = true
	}

	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()

	ds := filescan.DirectoryScanner{
		Recursive:   fc.Recursive,
		PrintErrors: fc.VeryVerbose,
		Password:    fc.Password,
	}
	ch := ds.ScanDirectories(ctx, args)
	for {
		select {
		case <-ctx.Done():
			return nil
		case tpls, ok := <-ch:
			if !ok {
				return nil
			}
			pts := fc.filterTemplates(tpls)
			if len(pts) > 0 {
				if err := fc.listTemplates(pts); err != nil {
					return err
				}
			}
		}
	}
}

// listTemplates prints out the given list of templates
func (fc FindCommand) listTemplates(tps []*printedTemplate) error {
	out := tabwriter.NewWriter(os.Stdout, 24, 8, 4, ' ', 0)
	defer func() {
		if err := out.Flush(); err != nil {
			log.Println(err)
		}
	}()
	for _, pt := range tps {
		t := pt.Template
		s := t.String()
		if !fc.Verbose {
			ss := strings.Split(s, "\t")
			s = ss[len(ss)-1]
			if len(ss) > 1 {
				s = strings.Join([]string{ss[0], s}, "\t")
			}
		}
		if pt.QueryResult != "" {
			_, _ = fmt.Fprintf(out, "(%s)\t", pt.QueryResult)
		}
		_, _ = fmt.Fprint(out, s)

		_, _ = fmt.Fprintln(out)
	}
	return nil
}

// filterTemplates filters the list of templates, first by type, then by and query string.
// returns a subset of the given templates, which qualify by type and query
func (fc FindCommand) filterTemplates(tps []templates.Template) []*printedTemplate {
	var ntps []*printedTemplate
	for _, t := range tps {
		// filter by template type
		tt := templates.TemplateType(t)
		if len(fc.Type) > 0 && indexOf(tt, fc.Type) < 0 {
			continue
		}
		qr, ok := fc.queryTemplate(t)
		if !ok {
			continue
		}
		ntps = append(ntps, &printedTemplate{
			Template:    t,
			QueryResult: qr,
		})
	}
	return ntps
}

// queryTemplate queries the given template against the command Query string.
// when present the temaplte is searched for that string, returning true if its found.
func (fc FindCommand) queryTemplate(t templates.Template) (string, bool) {
	if fc.Query == "" {
		return "", true
	}
	by, err := yaml.Marshal(t)
	if err != nil {
		log.Println(err)
		return err.Error(), false
	}
	i := findLineIndex(string(by), fc.Query, !fc.Insensitive)
	if i < 0 {
		return "", false
	}
	s := strings.Split(string(by), "\n")
	return strings.TrimSpace(s[i]), true
}

func findLineIndex(s string, q string, cs bool) int {
	if !cs {
		s = strings.ToLower(s)
		q = strings.ToLower(q)
	}
	ss := strings.Split(s, "\n")
	for i, l := range ss {
		if strings.Contains(l, q) {
			if strings.Contains(l, "ocation:") {
				continue
			}
			return i
		}
	}
	return -1
}

func indexOf(s string, ss []string) int {
	for i, sz := range ss {
		if s == sz {
			return i
		}
	}
	return -1
}

func appendUnique(ss []string, s ...string) []string {
	for _, sz := range s {
		if indexOf(sz, ss) >= 0 {
			continue
		}
		ss = append(ss, sz)
	}
	return ss
}

// formatTypeNames takes the abbreviated names from the -type flag and converts them into template type names.
func formatTypeNames(types []string) ([]string, error) {
	var ts []string

	for _, tt := range types {
		if tt == "" {
			continue
		}
		switch strings.ToLower(tt) {
		case "*", "any":
			return nil, nil

		case "cer", "certificate", "certificates":
			ts = appendUnique(ts, "CERTIFICATE")

		case "csr", "req", "request", "certificaterequest":
			ts = appendUnique(ts, "CERTIFICATE REQUEST")

		case "puk", "public", "publickey", "publickeys":
			ts = appendUnique(ts, "PUBLIC KEY", "SSH PUBLIC KEY")

		case "prk", "private", "privatekey", "privatekeys":
			ts = appendUnique(ts, "PRIVATE KEY", "OPENSSH PRIVATE KEY")

		case "keys":
			ts = appendUnique(ts, "PUBLIC KEY", "SSH PUBLIC KEY",
				"PRIVATE KEY", "OPENSSH PRIVATE KEY")

		case "crl", "rev", "revocation", "revocationlist":
			ts = appendUnique(ts, "X509 CRL")

		default:
			return nil, fmt.Errorf("'%s' is not a known resource type.  Use 'cer', 'csr', 'crl', 'prk', 'puk' or 'any'", tt)
		}
	}
	return ts, nil
}

type printedTemplate struct {
	Template    templates.Template
	QueryResult string
}
