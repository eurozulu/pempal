package templates

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/eurozulu/pempal/model"
	"gopkg.in/yaml.v2"
	"math/big"
	"time"
)

type RevocationListTemplate struct {
	// Issuer contains the DN of the issuing certificate.
	Issuer model.DistinguishedName `yaml:"issuer"`

	// AuthorityKeyId is used to identify the public key associated with the
	// issuing certificate. It is populated from the authorityKeyIdentifier
	// extension when parsing a CRL. It is ignored when creating a CRL; the
	// extension is populated from the issuing certificate itself.
	AuthorityKeyId []byte

	Signature          model.Base64Binary       `yaml:"signature,omitempty"`
	SignatureAlgorithm model.SignatureAlgorithm `yaml:"signature-algorithm"`

	// RevokedCertificateEntries represents the revokedCertificates sequence in
	// the CRL. It is used when creating a CRL and also populated when parsing a
	// CRL. When creating a CRL, it may be empty or nil, in which case the
	// revokedCertificates ASN.1 sequence will be omitted from the CRL entirely.
	RevokedCertificateEntries []x509.RevocationListEntry

	// RevokedCertificates is used to populate the revokedCertificates
	// sequence in the CRL if RevokedCertificateEntries is empty. It may be empty
	// or nil, in which case an empty CRL will be created.
	//
	// Deprecated: Use RevokedCertificateEntries instead.
	RevokedCertificates []pkix.RevokedCertificate

	// Number is used to populate the X.509 v2 cRLNumber extension in the CRL,
	// which should be a monotonically increasing sequence number for a given
	// CRL scope and CRL issuer. It is also populated from the cRLNumber
	// extension when parsing a CRL.
	Number *big.Int

	// ThisUpdate is used to populate the thisUpdate field in the CRL, which
	// indicates the issuance date of the CRL.
	ThisUpdate model.TimeDTO
	// NextUpdate is used to populate the nextUpdate field in the CRL, which
	// indicates the date by which the next CRL will be issued. NextUpdate
	// must be greater than ThisUpdate.
	NextUpdate model.TimeDTO

	// Extensions contains raw X.509 extensions. When creating a CRL,
	// the Extensions field is ignored, see ExtraExtensions.
	Extensions []model.Extension

	// ExtraExtensions contains any additional extensions to add directly to
	// the CRL.
	ExtraExtensions []model.Extension
}

func (ct RevocationListTemplate) String() string {
	s := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(s).Encode(&ct); err != nil {
		return err.Error()
	}
	return s.String()
}

func (ct RevocationListTemplate) ApplyTo(list *model.RevocationList) error {
	if ct.Issuer.String() != "" {
		dn, err := model.ParseDistinguishedName(list.Issuer.String())
		if err != nil {
			return err
		}
		dn.Merge(ct.Issuer)
		list.Issuer = pkix.Name(*dn)
	}

	if len(ct.AuthorityKeyId) > 0 {
		list.AuthorityKeyId = ct.AuthorityKeyId
	}
	if len(ct.Signature) > 0 {
		list.Signature = ct.Signature
	}
	if x509.SignatureAlgorithm(ct.SignatureAlgorithm) != x509.UnknownSignatureAlgorithm {
		list.SignatureAlgorithm = x509.SignatureAlgorithm(ct.SignatureAlgorithm)
	}
	if len(ct.RevokedCertificateEntries) > 0 {
		list.RevokedCertificateEntries = ct.RevokedCertificateEntries
	}
	if len(ct.RevokedCertificates) > 0 {
		list.RevokedCertificates = ct.RevokedCertificates
	}
	if ct.Number != nil && ct.Number.Uint64() != 0 {
		list.Number = ct.Number
	}
	if !time.Time(ct.ThisUpdate).IsZero() {
		list.ThisUpdate = time.Time(ct.ThisUpdate)
	}
	if !time.Time(ct.NextUpdate).IsZero() {
		list.NextUpdate = time.Time(ct.NextUpdate)
	}
	return nil
}

func NewRevocationListTemplate(r *model.RevocationList) *RevocationListTemplate {
	return &RevocationListTemplate{
		Issuer:                    model.DistinguishedName(r.Issuer),
		AuthorityKeyId:            r.AuthorityKeyId,
		Signature:                 r.Signature,
		SignatureAlgorithm:        model.SignatureAlgorithm(r.SignatureAlgorithm),
		RevokedCertificateEntries: r.RevokedCertificateEntries,
		Number:                    r.Number,
		ThisUpdate:                model.TimeDTO(r.ThisUpdate),
		NextUpdate:                model.TimeDTO(r.NextUpdate),
		Extensions:                model.ExtensionsToModel(r.Extensions),
		ExtraExtensions:           nil,
	}
}
