package commands

import (
	"crypto/x509"
	"fmt"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/resources"
	"io"
)

type IssuersCommand struct {
	Output       io.Writer
	KeyedIssuers bool `flag:"with-key,k,own"`
}

func (i IssuersCommand) Exec(args ...string) error {
	issurs := resources.NewIssuers(config.Config.CertPath)
	var certs []*x509.Certificate
	if !i.KeyedIssuers {
		certs = issurs.Issuers()
	} else {
		certs = issurs.IssuersWithKeys()
	}
	if len(certs) == 0 {
		if !CommonFlags.Quiet {
			var k string
			if i.KeyedIssuers {
				k = " keyed"
			}
			fmt.Printf("No%s issuer certificates found\n", k)
		}
		return nil
	}

	return nil
}
