package writers

import (
	"crypto/x509"
	"encoding/pem"
	"io"
)

type PEMWriter struct{}

func (pw PEMWriter) Write(c *x509.Certificate, out io.Writer) error {
	blk := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: c.Raw,
	}
	return pem.Encode(out, blk)
}
