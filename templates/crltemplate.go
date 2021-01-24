package templates

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"math/big"
	"strings"
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
	crl      pkix.CertificateList
}

func (t CRLTemplate) Location() string {
	return t.FilePath
}

func (t *CRLTemplate) String() string {
	num := " "
	if t.Number != nil {
		num = t.Number.String()
	}
	return strings.Join([]string{TemplateType(t), t.ThisUpdate.String(), num,
		t.Location()}, "\t")
}

func (t CRLTemplate) MarshalBinary() (data []byte, err error) {
	// sync binary copy with struct fields
	t.crl.SignatureValue = t.SignatureValue
	t.crl.SignatureAlgorithm = t.SignatureAlgorithm
	t.crl.TBSCertList.RevokedCertificates = RevokedCertificatesReslice(t.RevokedCertificates)
	t.crl.TBSCertList.ThisUpdate = t.ThisUpdate
	t.crl.TBSCertList.NextUpdate = t.NextUpdate
	t.crl.TBSCertList.Extensions = ExtensionReslice(t.ExtraExtensions)
	return asn1.Marshal(t.crl)
}

func (t *CRLTemplate) UnmarshalBinary(by []byte) error {
	cl, err := x509.ParseDERCRL(by)
	if err != nil {
		return err
	}
	t.crl = *cl
	t.SignatureAlgorithm = cl.SignatureAlgorithm
	t.RevokedCertificates = RevokedCertificatesSlice(cl.TBSCertList.RevokedCertificates)

	t.ThisUpdate = cl.TBSCertList.ThisUpdate
	t.NextUpdate = cl.TBSCertList.NextUpdate
	t.ExtraExtensions = ExtensionSlice(cl.TBSCertList.Extensions)
	return nil
}
