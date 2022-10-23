package main

import (
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"os"
	"pempal/templates"
	"pempal/writers"
	"strings"
)

const FLAG_KEY_ISSUER = "issuer"

const defaultFormat = "pem"

type issueCommand struct {
	Key    string `flag:"key"`
	Issuer string `flag:"issuer"`
}

func (i issueCommand) Run(out io.Writer, args ...string) error {
	ownerKey := keyFromArgs()
	if ownerKey == nil {
		return fmt.Errorf("no key found to sign certificate")
	}

	params, flags := parseArgs(args...)
	if len(params) == 0 {
		return fmt.Errorf("must specifiy at least one template name to create.")
	}
	temps, err := templates.TemplatePath.Find(params...)
	if err != nil {
		return err
	}

	c, err := applyTemplates(temps)
	if err != nil {
		return err
	}

	x509.CreateCertificate(rand.Reader, c, issuer, k.PublicKey(), k.)
	out := os.Stdout
	op, _ := FlagValue(FLAG_OUTPUT, args...)
	if op != "" {
		f, err := os.OpenFile(op, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0640)
		if err != nil {
			return fmt.Errorf("failed to open output file %s  %v", op, err)
		}
		defer func(f *os.File) {
			if err := f.Close(); err != nil {
				log.Printf("Failed to close %s  %v", f.Name(), err)
			}
		}(f)
		out = f
	}
	writer := certWriter(args...)
	if writer == nil {
		return fmt.Errorf("unknown output format")
	}
	return writer.Write(c, out)
}

func applyTemplates(temps []templates.Template) (*x509.Certificate, error) {
	c := &x509.Certificate{}
	for _, t := range temps {
		if err := t.Apply(c); err != nil {
			return nil, fmt.Errorf("failed to apply template  %v", err)
		}
	}
	return c, nil
}

func issuerCert(args ...string) (*x509.Certificate, error) {
	ip, ok := FlagValue(FLAG_KEY_ISSUER, args...)
	if !ok {
		// no issuer stated, assume self signed
		return nil, nil
	}
	r := ResourcePath.Find(ip)
	if r == nil {
		return nil, fmt.Errorf("issuer %s could not be found", ip)
	}
	r.
}

func certWriter(args ...string) writers.CertificateWriter {
	f, ok := FlagValue(FLAG_KEY_FORMAT)
	if ok && f != "" {
		f = strings.ToLower(f)
		for k, v := range writers.CertificateWriters {
			if f == k {
				return v
			}
		}
		return nil
	}
	return writers.CertificateWriters[defaultFormat]
}

func parseArgs(args ...string) (params []string, flags map[string]*string) {
	flags = map[string]*string{}
	var i int
	for ; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			arg := strings.TrimLeft(args[i], "-")
			var v *string
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i++
				*v = args[i]
			}
			flags[arg] = v
		} else {
			params = append(params, args[i])
		}
	}
	return params, flags
}

func newIssueCommand() Command {
	return &issueCommand{}
}
