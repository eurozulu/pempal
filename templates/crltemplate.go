package templates

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

type CRLTemplate struct {
	SignatureValue      asn1.BitString           `yaml:"SignatureValue,omitempty"`
	SignatureAlgorithm  pkix.AlgorithmIdentifier `yaml:"SignatureAlgorithm,omitempty"`
	RevokedCertificates []RevokedCertificate     `yaml:"RevokedCertificates,omitempty"`
	Number              *big.Int                 `yaml:"Number,omitempty"`
	ThisUpdate          time.Time                `yaml:"ThisUpdate,omitempty"`
	NextUpdate          time.Time                `yaml:"NextUpdate,omitempty"`
	ExtraExtensions     []Extension              `yaml:"ExtraExtensions,omitempty"`

	FilePath string `yaml:"-"`
	crl      *pkix.CertificateList
}

func (t CRLTemplate) Location() string {
	return t.FilePath
}

func (t *CRLTemplate) String() string {
	return TemplateString(t)
}

func (t CRLTemplate) MarshalBinary() (data []byte, err error) {
	cl := t.crl
	if cl == nil {
		cl = &pkix.CertificateList{}
	}
	cl.SignatureValue = t.SignatureValue
	cl.SignatureAlgorithm = t.SignatureAlgorithm
	cl.TBSCertList.RevokedCertificates = RevokedCertificatesReslice(t.RevokedCertificates)
	cl.TBSCertList.ThisUpdate = t.ThisUpdate
	cl.TBSCertList.NextUpdate = t.NextUpdate
	cl.TBSCertList.Extensions = ExtensionReslice(t.ExtraExtensions)

	return asn1.Marshal(cl)
}

func (t *CRLTemplate) UnmarshalBinary(by []byte) error {
	cl, err := x509.ParseDERCRL(by)
	if err != nil {
		return err
	}
	t.crl = cl
	t.SignatureAlgorithm = cl.SignatureAlgorithm
	t.RevokedCertificates = RevokedCertificatesSlice(cl.TBSCertList.RevokedCertificates)

	t.ThisUpdate = cl.TBSCertList.ThisUpdate
	t.NextUpdate = cl.TBSCertList.NextUpdate
	t.ExtraExtensions = ExtensionSlice(cl.TBSCertList.Extensions)
	return nil
}

func (t *CRLTemplate) UnmarshalPEM(bl *pem.Block) error {
	if bl.Type != "X509 CRL" {
		return fmt.Errorf("%s is not a CRL pem type", bl.Type)
	}
	return t.UnmarshalBinary(bl.Bytes)
}

func (t CRLTemplate) MarshalPEM() (*pem.Block, error) {
	by, err := t.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &pem.Block{
		Type:  "X509 CRL",
		Bytes: by,
	}, nil
}
