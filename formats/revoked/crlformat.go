package revoked

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
	"math/big"
	"pempal/formats"
	"pempal/resources"
	"pempal/templates"
)

type crlFormat struct {
}

func (fm crlFormat) Make(t templates.Template) (resources.Resource, error) {
	ct, err := fm.ensureRevokeTemplate(t)
	if err != nil {
		return nil, err
	}

	if err = fm.validateTemplate(*ct); err != nil {
		return nil, err
	}

	var revoked x509.RevocationList
	fm.applyTemplate(&revoked, ct)

	// Establish issuer.
	issuer, err := formats.resolveIssuer(*ct.Issuer, nil)
	if err != nil {
		return nil, err
	}

	// locate issuer private key & ensure it's a Signer key
	prk, err := formats.matchPrivateKeyFromPublic(issuer.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("missing issuer private key for %s", issuer.Subject.CommonName)
	}
	priv, ok := prk.(crypto.Signer)
	if !ok {
		return nil, fmt.Errorf("issuer %s private key is not valid for signing", issuer.Subject.CommonName)
	}
	by, err := x509.CreateRevocationList(rand.Reader, &revoked, issuer, priv)
	if err != nil {
		return nil, err
	}
	return resources.NewResource("", &pem.Block{
		Type:  "X509 CRL",
		Bytes: by,
	}), nil

}

func (fm crlFormat) ensureRevokeTemplate(t templates.Template) (*templates.CRLTemplate, error) {
	ct, ok := t.(*templates.CRLTemplate)
	if !ok {
		ct = &templates.CRLTemplate{}
		by, err := yaml.Marshal(t)
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(by, ct); err != nil {
			return nil, err
		}
	}
	return ct, nil
}

func (fm crlFormat) applyTemplate(crl *x509.RevocationList, t *templates.CRLTemplate) error {
	formats.nameFormat{}.applyTemplate(&crl.Issuer, t.Issuer)
	crl.SignatureAlgorithm = formats.parseSignatureAlgorithm(t.SignatureAlgorithm)
	nt, err := formats.parseTime(t.NextUpdate)
	if err != nil {
		return fmt.Errorf("invalid next-update time")
	}
	crl.NextUpdate = nt
	nt, err = formats.parseTime(t.ThisUpdate)
	if err != nil {
		return fmt.Errorf("invalid this-update time")
	}
	crl.ThisUpdate = nt
	crl.Number = big.NewInt(t.Number)

	rcs, err := parseRevokeCertificates(t.RevokedCertificates)
	if err != nil {
		return fmt.Errorf("invalid revoke certificate %v", err)
	}
	crl.RevokedCertificates = rcs
	return nil
}

func (fm crlFormat) validateTemplate(ct templates.CRLTemplate) error {
	if ct.Issuer == nil {
		return fmt.Errorf("missing issuer")
	}
	if ct.Issuer.IsEmpty() {
		return fmt.Errorf("missing issuer common-name")
	}

	return nil
}

func parseRevokeCertificates(temps []*templates.RevokedCertificateTemplate) ([]pkix.RevokedCertificate, error) {
	var rcs []pkix.RevokedCertificate
	rFmt := formats.revokedCertificateFormat{}
	for _, t := range temps {
		c := pkix.RevokedCertificate{}
		if err := rFmt.applyTemplate(&c, *t); err != nil {
			return nil, err
		}
		rcs = append(rcs, c)
	}
	return rcs, nil
}
