package writers

import (
	"crypto/x509"
	"io"
)

type DERWriter struct {
}

func (dw DERWriter) Write(c *x509.Certificate, out io.Writer) error {
	_, err := out.Write(c.Raw)
	return err
}
