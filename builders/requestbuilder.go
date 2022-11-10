package builders

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
	"pempal/encoders"
	"pempal/pemtypes"
	"pempal/templates"
)

type requestBuilder struct {
	buildTemplate  templates.CSRTemplate
	localResources []*pem.Block
}

func (r *requestBuilder) AddResource(p *pem.Block) {
	if r.ContainsResource(p) {
		return
	}
	r.localResources = append(r.localResources, p)
}

func (r requestBuilder) AddTemplate(temps ...templates.Template) error {
	for _, t := range temps {
		txt, err := t.(encoding.TextMarshaler).MarshalText()
		if err != nil {
			return err
		}
		if err = yaml.NewDecoder(bytes.NewBuffer(txt)).Decode(&r.buildTemplate); err != nil {
			return err
		}
	}
	return nil
}

func (r requestBuilder) Validate() []error {
	var errs []error
	if r.buildTemplate.PublicKey == "" {
		errs = append(errs, fmt.Errorf(ErrMissingValue.Error(), "public-key"))
	}
	if r.buildTemplate.SignatureAlgorithm == "" {
		errs = append(errs, fmt.Errorf(ErrMissingValue.Error(), "signature-algorithm"))
	}
	if r.buildTemplate.Subject == nil {
		errs = append(errs, fmt.Errorf(ErrMissingValue.Error(), "subject"))
	} else if r.buildTemplate.Subject.IsEmpty() {
		errs = append(errs, fmt.Errorf(ErrMissingValue.Error(), "subject.common-name"))
	}
	errs = append(errs, r.resolveResources())
	return errs
}

func (r requestBuilder) resolveResources() error {
	keys := r.ResourcesByType(pemtypes.PrivateKey)
	if len(keys) == 0 {
		return fmt.Errorf(ErrMissingValue.Error(), "private-key")
	}
	if len(keys) > 1 {
		return fmt.Errorf("multiple private-keys found")
	}
	return nil
}

func (r requestBuilder) ContainsResource(p *pem.Block) bool {
	for _, b := range r.localResources {
		if bytes.Equal(b.Bytes, p.Bytes) {
			return true
		}
	}
	return false
}

func (r requestBuilder) ResourcesByType(pemType pemtypes.PEMType) []*pem.Block {
	var blocks []*pem.Block
	for _, b := range r.localResources {
		pt := pemtypes.ParsePEMType(b.Type)
		if pt != pemType {
			continue
		}
		blocks = append(blocks, b)
	}
	return blocks
}

func (r requestBuilder) Build() ([]*pem.Block, error) {
	errs := r.Validate()
	if len(errs) > 0 {
		return nil, fmt.Errorf("%v", errs)
	}

	prk, err := r.resolvePrivateKey()
	if err != nil {
		return nil, err
	}

	var csr x509.CertificateRequest
	d := encoders.RequestDecoder{}
	d.ApplyTemplate(&r.buildTemplate, &csr)
	der, err := x509.CreateCertificateRequest(rand.Reader, &csr, prk)
	if err != nil {
		return nil, err
	}
	return []*pem.Block{&pem.Block{
		Type:  pemtypes.Request.String(),
		Bytes: der,
	}}, nil
}

func (r requestBuilder) resolvePrivateKey() (crypto.PrivateKey, error) {
	prks := r.ResourcesByType(pemtypes.PrivateKey)
	if len(prks) != 1 {
		if len(prks) == 0 {
			return nil, fmt.Errorf("no private key found")
		}
		return nil, fmt.Errorf("multiple private keys")
	}
	return encoders.ParsePrivateKey(prks[0].Bytes)
}
