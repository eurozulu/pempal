package parsers

import (
	"bytes"
	"crypto/x509"
	"gopkg.in/yaml.v3"
	"pempal/templates"
)

type CSRCertificateParser struct {
}

func (cp CSRCertificateParser) Parse(by []byte) (templates.Template, error) {
	var err error
	var csr *x509.CertificateRequest
	csr, err = x509.ParseCertificateRequest(by)
	if err != nil {
		return nil, err
	}
	var yBy []byte
	yBy, err = csrToYaml(csr)
	if err != nil {
		return nil, err
	}
	yBy, err = cleanYaml(yBy)
	if err != nil {
		return nil, err
	}
	return templates.NewTemplate(yBy), nil
}

func csrToYaml(c *x509.CertificateRequest) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(buf).Encode(c); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
