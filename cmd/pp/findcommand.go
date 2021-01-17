package main

import (
	"context"
	"fmt"
	"github.com/eurozulu/pempal/filescan"
	"github.com/eurozulu/pempal/templates"
	"io"
	"log"
	"os"
	"reflect"
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

	// NoHeaders surpresses the output of column header names
	NoHeaders bool `flag:"noheaders,nohead,n"`

	// Verbose adds additional fields about the resource to the output.
	// When true, implies -table to format the output in columns
	Verbose bool `flag:"verbose,v"`
	// VeryVerbose adds to verbose output by listing each file visited
	VeryVerbose bool `flag:"vv"`

	// Certificate flag, when present, will find certificate resources matching the given properties
	// specify query properties in curly crackets:  {"issuedby": "My CA certificate"}
	Certificate FilterCriteria `flag:"certificate,cer,optionalvalue"`

	// Request flag, when present, will find certificate Request resources matching the given properties
	// specify query properties in curly crackets:  {"commonname": "My requested certificate"}
	Request FilterCriteria `flag:"request,csr,optionalvalue"`

	// Revocation flag, when present, will find Revocation list resources matching the given properties
	Revocation FilterCriteria `flag:"revocation,crl,optionalvalue"`

	// PublicKey flag, when present, will find PublicKey resources matching the given properties
	PublicKey FilterCriteria `flag:"publickey,puk,optionalvalue"`

	// PrivateKey flag, when present, will find PrivateKey resources matching the given properties
	PrivateKey FilterCriteria `flag:"privatekey,prk,optionalvalue"`

	// Any flag, when present, will find any resource containung the given properties
	Any FilterCriteria `flag:"any,optionalvalue"`

	// Chain attempts to locate the full certificate chain of each certificate found.
	// Any certificate not self signed will invoke a secondary find for its issuer.
	// These continue until the self signed root is found.
	Chain bool `flag:"chain,c"`

	// Keys attempts to locate the related keys to any found resource.
	// Any signed resource (cert, csr, crl) starts a find for its public key file.
	Keys bool `flag:"keys,k"`

	tFilters   map[string]FilterCriteria
	fieldNames []string
}

// Find locates filescan (Certificates, CSRs, CRLs, Private & Public keys) based on query criteria.
// With no flags, it simply lists the resource location, depending on its container:
// A file in a directory has the full path listed if it contains a single item.
// files containing more than one item are listed with a sub index:
// ./thisdir/certs/servertcerts.p12#3 Indicating it is the 3rd resource in this file.
// Arguments must include at least one location to search.  Any number of space delimited locations can be given,
// each will be searched.
// The argument may be the path to a specific file or a container such as a directory or pkcs container file.
// Optional Flags:
// -recursive  Search the named locations and any sub location found within them. default false
// -chain search for an issuer chain of certificates for any certificate found.
// -keys search for the related private and public keys for any signed resource found.
//
// Resources can be 'filtered' by their types and properties.
// Filter flags are:
// -cer		finds only certificates
// -csr		Finds only certiifcate requests
// -crl		Finds only certiifcate revokation lists
// -puk		Finds only public keys
// -prk		Finds only private keys
// -any 	Finds any and all resources.
// These flags have an optional value which can list one or more properties within the resource and a value
// of that property to match against.  Properties are expresed as {<key name>: <key value>,...}
// e.g. To find all certificates issued by a certain CA:
// Find -cer {issuedby: "my root ca certificate"}
// Certificate dates can be expressed as yyyy/mm/dd hh/MM/ss format or [+|-] duration.
// e.g. to find certificates that will expire in the next month:
// find -cer{notafter: "+1month"}
func (fc FindCommand) Find(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Find: must provide at least one path")
	}

	// Collect any template filters
	fc.tFilters = map[string]FilterCriteria{}
	if fc.Certificate != nil {
		fc.tFilters["CERTIFICATE"] = fc.Certificate
	}
	if fc.Request != nil {
		fc.tFilters["CERTIFICATE REQUEST"] = fc.Request
	}
	if fc.Revocation != nil {
		fc.tFilters["X509 CRL"] = fc.Revocation
	}
	if fc.PublicKey != nil {
		fc.tFilters["PUBLIC KEY"] = fc.PublicKey
		fc.tFilters["SSH PUBLIC KEY"] = fc.PublicKey
	}
	if fc.PrivateKey != nil {
		fc.tFilters["PRIVATE KEY"] = fc.PrivateKey
		fc.tFilters["OPENSSH PRIVATE KEY"] = fc.PrivateKey
	}
	if len(fc.Any) > 0 {
		fc.tFilters["ANY"] = fc.Any
	}

	// If no filters set, add an empty Any
	if len(fc.tFilters) == 0 {
		fc.tFilters["ANY"] = FilterCriteria{}
	}

	out := tabwriter.NewWriter(os.Stdout, 8, 4, 1, ' ', 0)
	defer func() {
		if err := out.Flush(); err != nil {
			log.Println(err)
		}
	}()

	// Build a list of the field names to display
	fc.fieldNames = fc.titleNames()
	if !fc.NoHeaders {
		_, err := fmt.Fprintln(out, strings.Join(fc.fieldNames, "\t"))
		if err != nil {
			return err
		}
	}

	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()

	ds := filescan.DirectoryScanner{
		Recursive:   fc.Recursive,
		PrintErrors: fc.VeryVerbose,
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
			tpls = fc.filterTemplates(tpls)
			if len(tpls) > 0 {
				if err := fc.listTemplates(tpls, out); err != nil {
					return err
				}
			}
		}
	}
}

// listTemplates prints out the given list of templates
func (fc FindCommand) listTemplates(tps []templates.Template, out io.Writer) error {
	tVals := map[string]interface{}{}
	for _, t := range tps {
		if fc.Verbose || fc.VeryVerbose {
			// Get template values if verbose output
			tVals = templates.TemplateValues(t, fc.fieldNames)
		}
		tVals["Location"] = t.Location()

		// line up values according to titles
		dVals := make([]string, len(fc.fieldNames))
		for f, v := range tVals {
			pos := indexOf(f, fc.fieldNames)
			if pos < 0 {
				continue
			}
			dVals[pos] = reflect.ValueOf(v).String()
		}
		_, err := fmt.Fprintln(out, strings.Join(dVals, "\t"))
		if err != nil {
			return err
		}
	}
	return nil
}

func (fc FindCommand) filterTemplates(tps []templates.Template) []templates.Template {
	var ntps []templates.Template
	for _, t := range tps {
		flts := fc.findTemplateFilters(t)
		// no filters for this template type, ignore template
		if len(flts) == 0 {
			continue
		}
		// check template matches the filters found
		for _, f := range flts {
			if f.Match(t) {
				ntps = append(ntps, t)
			}
		}
	}
	return ntps
}

func (fc FindCommand) findTemplateFilters(t templates.Template) []FilterCriteria {
	var cr []FilterCriteria
	tp, ok := fc.tFilters["ANY"]
	if ok {
		cr = append(cr, tp)
	}
	tp, ok = fc.tFilters[templates.TemplateType(t)]
	if ok {
		cr = append(cr, tp)
	}
	return cr
}

func (fc FindCommand) titleNames() []string {
	var fields []string
	if fc.Verbose || fc.VeryVerbose {
		for tt, f := range fc.tFilters {
			// Empty filters = All items, display first 3 fields
			flds := f.FieldNames()
			if len(flds) == 0 && tt != "ANY" {
				flds = firstTemplateFields(tt, 5)
			}
			fields = appendUnique(fields, flds...)
		}
	}
	fields = appendUnique(fields, "Location")
	return fields
}

// FilterCriteria is a map of property name keys with the values to match against the resource value.
type FilterCriteria map[string]interface{}

// Gets the field names this criteria will filter on.
func (fc FilterCriteria) FieldNames() []string {
	var n []string
	for k := range fc {
		n = append(n, k)
	}
	return n
}

// Match the given template with the criteria.
// All properties in the criteria must be present in the template and
// the value of the template property must be equal to the criteria value.
// If critiera has no named properties, it will match all templates.
func (fc FilterCriteria) Match(t templates.Template) bool {
	// Empty filter allows all templates
	if len(fc) == 0 {
		return true
	}

	vals := templates.TemplateValues(t, fc.FieldNames())
	// If template doesn't have all the filter names, not a match
	if len(vals) < len(fc) {
		return false
	}

	// compare each value with filter value
	for k, cv := range fc {
		vv, ok := vals[k]
		if !ok {
			return false
		}
		if cv != vv {
			return false
		}
	}
	return true
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

// Gets the first 'n' field names of a Template from the given type
func firstTemplateFields(tType string, count int) []string {
	t, err := templates.NewTemplate("", tType)
	if err != nil {
		log.Println(err)
		return nil
	}
	flds := templates.TemplateFields(t)
	if count > len(flds) {
		count = len(flds)
	}
	return flds[0:count]
}
