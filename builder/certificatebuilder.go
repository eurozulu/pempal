package builder

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/keys"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/templates"
	"reflect"
)

type certificateBuilder struct {
	keys keys.Keys
}

func (cb certificateBuilder) Validate(t templates.Template) []error {
	c, errs := cb.buildTemplateCertificate(t)
	if len(errs) > 0 {
		return errs
	}
	if _, err := cb.resolveIssuer(c.Issuer); err != nil {
		return []error{err}
	}
	return nil
}

func (cb certificateBuilder) Build(t templates.Template) (model.Resource, error) {
	tcert, errs := cb.buildTemplateCertificate(t)
	if len(errs) > 0 {
		return nil, CombineErrors(errs)
	}
	issuer, err := cb.resolveIssuer(tcert.Issuer)
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

func (cb certificateBuilder) buildTemplateCertificate(t templates.Template) (*x509.Certificate, []error) {
	dto := model.NewDTOForResourceType(model.Certificate)
	certdto, ok := dto.(*model.CertificateDTO)
	if !ok {
		return nil, []error{fmt.Errorf("unexpected DTO type %s, expected CertificateDTO", reflect.TypeOf(dto).String())}
	}
	if err := t.Apply(certdto); err != nil {
		return nil, []error{err}
	}
	c, err := certdto.ToCertificate()
	if err != nil {
		return nil, []error{err}
	}
	var errs []error
	if c.PublicKey == nil {
		errs = append(errs, fmt.Errorf("public key not found"))
	}
	if c.Subject.CommonName == "" {
		errs = append(errs, fmt.Errorf("common name not found"))
	}
	if c.SignatureAlgorithm == x509.UnknownSignatureAlgorithm {
		errs = append(errs, fmt.Errorf("signature algorithm unknown"))
	}
	if c.SerialNumber.Uint64() == 0 {
		errs = append(errs, fmt.Errorf("serial number not found"))
	}
	if c.NotBefore.IsZero() {
		errs = append(errs, fmt.Errorf("before-now invalid"))
	}
	if !c.NotAfter.After(c.NotBefore) {
		errs = append(errs, fmt.Errorf("after-now invalid, after before-now"))
	}
	if len(errs) > 0 {
		return nil, errs
	}
	return c, nil
}

func (cb certificateBuilder) resolveIssuer(issuer pkix.Name) (keys.User, error) {
	if issuer.CommonName == "" {
		return nil, fmt.Errorf("no issuer common name")
	}
	u, err := cb.keys.UserByName(issuer)
	if err != nil {
		return nil, err
	}
	if !u.Certificate().IsCA {
		return nil, fmt.Errorf("user %s is not an issuer", issuer.CommonName)
	}
	return u, nil
}
