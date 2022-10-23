package writers

import (
	"crypto/x509"
	"gopkg.in/yaml.v3"
	"io"
	"pempal/templates/parsers"
)

type TemplateWriter struct {
}

func (tw TemplateWriter) Write(c *x509.Certificate, out io.Writer) error {
	dw := parsers.DERCertificateParser{}
	t, err := dw.Parse(c.Raw)
	if err != nil {
		return err
	}
	return yaml.NewEncoder(out).Encode(t)
}
