package main

import (
	"fmt"
	"github.com/pempal/pempal"
	"github.com/pempal/templates"
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

	// Query is a free text search of a resource.
	// If the given string appears in the resource (As it would appear on screen)
	// It will be included in the results, otherwise it will be omitted.
	// e.g. -t cer -q "Root CA"  Will find any certificate with "Root CA" in any of its fields.
	// Query is case sensitive unless -insensitive is set
	Query []string `flag:"query,q"`

	// CaseSensitive effects the -query flag result, making it case sensitive when specified.
	CaseSensitive bool `flag:"casesensitive,c"`

	// Types specifies the type of resource to find.  Can be a comma delimited list of one or more of th efollowing:
	// 			cer		finds only certificates
	// 			csr		Finds only certiifcate requests
	// 			crl		Finds only certiifcate revokation lists
	// 			puk		Finds only public keys
	// 			prk		Finds only private keys
	// When not specified, finds All types.
	Types []string `flag:"type,t"`

	Password string `flag:"pass,p"`
}

// Find locates files (Certificates, CSRs, CRLs, Private & Public keys) based on query criteria.
//
// With no flags, it simply lists the resource location, depending on its container:
// A file in a directory has the full path listed if it contains a single item.
// files containing more than one item are listed with a sub index:
// ./thisdir/certs/servertcerts.p12#3 Indicating it is the 3rd resource in this file.
// The path arguments may be a path to a specific file or a container such as a directory or pkcs container file.
func (fc FindCommand) Find(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Find: must provide at least one path")
	}
	if fc.VeryVerbose {
		fc.Verbose = true
	}
	ts, err := formatTypeNames(fc.Types)
	if err != nil {
		return err
	}
	fc.Types = ts

	var count int
	for _, p := range args {
		ps := &pempal.PEMQuery{
			Types:    fc.Types,
			Query:    fc.Query,
			Password: fc.Password,
			Verbose:  fc.VeryVerbose,
		}
		qrs, err := ps.QueryPath(p, fc.Recursive)
		if err != nil {
			return err
		}
		count += len(qrs)
		if len(qrs) == 0 {
			continue
		}
		ListPems(qrs, fc.Verbose, fc.VeryVerbose, false)
	}
	if fc.Verbose {
		fmt.Printf("found %d resources\n", count)
	}
	return nil
}

// listPems prints out the given list of PEMs
func ListPems(queryResults []*pempal.QueryResult, verbose, showErrors, showlines bool) {
	out := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', 0)
	defer func() {
		if err := out.Flush(); err != nil {
			log.Println(err)
		}
	}()
	for i, qr := range queryResults {
		var s string
		if showlines {
			s = fmt.Sprintf("%d)\t", i+1)
		}
		s = strings.Join([]string{s, strings.Title(qr.Block.Type)}, "")

		if verbose {
			t, err := templates.NewTemplate(qr.Block)
			if err != nil {
				if showErrors {
					fmt.Println(err)
				}
				continue
			}
			s = strings.Join([]string{s, t.String()}, "\t")
		}

		if len(qr.QueryMatch) > 0 {
			qm := strings.Join(qr.QueryMatch, "\t")
			s = strings.Join([]string{s, qm}, "\t")
		}

		s = strings.Join([]string{s, qr.FilePath}, "\t")
		_, _ = fmt.Fprint(out, s)
		_, _ = fmt.Fprintln(out)
	}
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
			ts = appendUnique(ts, "PUBLIC KEY")

		case "prk", "private", "privatekey", "privatekeys":
			ts = appendUnique(ts, "PRIVATE KEY")

		case "keys":
			ts = appendUnique(ts, "PUBLIC KEY", "PRIVATE KEY")

		case "crl", "rev", "revocation", "revocationlist":
			ts = appendUnique(ts, "X509 CRL")

		default:
			return nil, fmt.Errorf("'%s' is not a known resource type.  Use 'cer', 'csr', 'crl', 'prk', 'puk' or 'any'", tt)
		}
	}
	return ts, nil
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
func indexOf(s string, ss []string) int {
	for i, sz := range ss {
		if s == sz {
			return i
		}
	}
	return -1
}
