package builder

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/keys"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/templates"
	"github.com/go-yaml/yaml"
)

type certificateBuilder struct {
	dto  model.CertificateDTO
	keys keys.Keys
}

func (cb *certificateBuilder) ApplyTemplate(tp ...templates.Template) error {
	for _, t := range tp {
		if err := yaml.Unmarshal(t.Bytes(), &cb.dto); err != nil {
			return err
		}
	}
	return nil
}

func (cb certificateBuilder) Validate() error {
	if _, err := cb.buildTemplateCertificate(); err != nil {
		return err
	}
	if _, err := cb.resolveIssuer(); err != nil {
		return err
	}
	return nil
}

func (cb certificateBuilder) Build() (model.Resource, error) {
	tcert, err := cb.buildTemplateCertificate()
	if err != nil {
		return nil, err
	}
	issuer, err := cb.resolveIssuer()
	if err != nil {
		return nil, err
	}

	der, err := x509.CreateCertificate(rand.Reader, tcert, issuer.Certificate(), tcert.PublicKey, issuer.Key())
	if err != nil {
		return nil, err
	}
	return model.NewResource(&pem.Block{
		Type:  model.Certificate.PEMString(),
		Bytes: der,
	}), nil
}

func (cb certificateBuilder) buildTemplateCertificate() (*x509.Certificate, error) {
	c, err := cb.dto.ToCertificate()
	if err != nil {
		return nil, err
	}
	errs := bytes.NewBuffer(nil)
	if c.PublicKey == nil {
		fmt.Fprintln(errs, "no public key found")
	}
	if c.Subject.CommonName == "" {
		fmt.Fprintln(errs, "no common name")
	}

	if c.SignatureAlgorithm == x509.UnknownSignatureAlgorithm {
		fmt.Fprintln(errs, "unknown signature algorithm")
	}
	if c.SerialNumber.Uint64() == 0 {
		fmt.Fprintln(errs, "no serial number")
	}
	if errs.Len() > 0 {
		return nil, fmt.Errorf("%s", errs.String())
	}
	return c, nil
}

func (cb certificateBuilder) resolveIssuer() (keys.User, error) {
	c, err := cb.dto.ToCertificate()
	if err != nil {
		return nil, err
	}
	if c.Issuer.CommonName == "" {
		return nil, fmt.Errorf("no issuer common name")
	}
	u, err := cb.keys.UserByName(c.Issuer)
	if err != nil {
		return nil, err
	}
	if !u.Certificate().IsCA {
		return nil, fmt.Errorf("user %s is not an issuer", c.Issuer.CommonName)
	}
	return u, nil
}
