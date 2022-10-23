package parsers

import (
	"bytes"
	"crypto/x509"
	"gopkg.in/yaml.v3"
	"pempal/templates"
)

type DERCertificateParser struct {
}

func (D DERCertificateParser) Parse(by []byte) (templates.Template, error) {
	var err error
	var c *x509.Certificate
	c, err = x509.ParseCertificate(by)
	if err != nil {
		return nil, err
	}
	var yBy []byte
	yBy, err = certToYaml(c)
	if err != nil {
		return nil, err
	}
	yBy, err = cleanYaml(yBy)
	if err != nil {
		return nil, err
	}
	return templates.NewTemplate(yBy), nil
}

func certToYaml(c *x509.Certificate) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(buf).Encode(c); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func cleanYaml(by []byte) ([]byte, error) {
	// Read as a map
	mIn := map[string]interface{}{}
	if err := yaml.NewDecoder(bytes.NewBuffer(by)).Decode(&mIn); err != nil {
		return nil, err
	}
	// copy map, ignoring all null values
	mClean := map[string]interface{}{}
	for k, v := range mIn {
		if v == nil {
			continue
		}
		mClean[k] = v
	}
	// encode back into yaml
	buf := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(buf).Encode(mClean); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
