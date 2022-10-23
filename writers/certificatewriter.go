package writers

import (
	"crypto/x509"
	"io"
)

var CertificateWriters = map[string]CertificateWriter{
	"der":      &DERWriter{},
	"pem":      &PEMWriter{},
	"template": &TemplateWriter{},
}

type CertificateWriter interface {
	Write(c *x509.Certificate, out io.Writer) error
}
