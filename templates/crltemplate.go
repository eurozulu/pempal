package templates

import (
	"crypto/x509"
	"github.com/eurozulu/pempal/model"
	"math/big"
	"strings"
	"time"
)

type CRLTemplate struct {
	Issuer                    model.DistinguishedName  `yaml:"issuer"`
	AuthorityKeyId            model.Base64Binary       `yaml:"authority-key-id,omitempty"`
	Signature                 model.Base64Binary       `yaml:"signature,omitempty"`
	SignatureAlgorithm        model.SignatureAlgorithm `yaml:"signature-algorithm,omitempty"`
	RevokedCertificateEntries []RevokedEntryTemplate   `yaml:"revoked-certificate-entries,omitempty"`
	Number                    *big.Int                 `yaml:"number,omitempty"`
	ThisUpdate                model.TimeDTO            `yaml:"this-update"`
	NextUpdate                model.TimeDTO            `yaml:"next-update"`
	//RevokedCertificates       nil
	//Extensions                nil
	//ExtraExtensions           nil
}

type RevokedEntryTemplate struct {
	SerialNumber   *big.Int
	RevocationTime model.TimeDTO
	ReasonCode     int
	//Extensions      nil
	//ExtraExtensions nil
}

func (re RevokedEntryTemplate) ToEntry() x509.RevocationListEntry {
	return x509.RevocationListEntry{
		SerialNumber:   re.SerialNumber,
		RevocationTime: time.Time(re.RevocationTime),
		ReasonCode:     re.ReasonCode,
		//Extensions:      nil,
		//ExtraExtensions: nil,
	}
}

func (rt CRLTemplate) revokedCertificateEntries() []x509.RevocationListEntry {
	rces := make([]x509.RevocationListEntry, len(rt.RevokedCertificateEntries))
	for i, rc := range rt.RevokedCertificateEntries {
		rces[i] = rc.ToEntry()
	}
	return rces
}

func (rt CRLTemplate) ToCRL() *x509.RevocationList {
	return &x509.RevocationList{
		Issuer:                    rt.Issuer.ToName(),
		AuthorityKeyId:            rt.AuthorityKeyId,
		Signature:                 rt.Signature,
		SignatureAlgorithm:        x509.SignatureAlgorithm(rt.SignatureAlgorithm),
		RevokedCertificateEntries: rt.revokedCertificateEntries(),
		Number:                    rt.Number,
		ThisUpdate:                time.Time(rt.ThisUpdate),
		NextUpdate:                time.Time(rt.NextUpdate),
	}
}

func (rt CRLTemplate) Name() string {
	return strings.ToLower(model.RevokationList.String())
}

func newRevokationListEntires(entries []x509.RevocationListEntry) []RevokedEntryTemplate {
	list := make([]RevokedEntryTemplate, len(entries))
	for i, entry := range entries {
		list[i] = RevokedEntryTemplate{
			SerialNumber:   entry.SerialNumber,
			RevocationTime: model.TimeDTO(entry.RevocationTime),
			ReasonCode:     entry.ReasonCode,
		}
	}
	return list
}

func NewCRLTemplate(crl *x509.RevocationList) *CRLTemplate {
	ct := CRLTemplate{}
	ct.Signature = crl.Signature
	ct.SignatureAlgorithm = model.SignatureAlgorithm(crl.SignatureAlgorithm)
	ct.Issuer = model.DistinguishedName(crl.Issuer)
	ct.AuthorityKeyId = crl.AuthorityKeyId
	ct.Signature = crl.Signature
	ct.SignatureAlgorithm = model.SignatureAlgorithm(crl.SignatureAlgorithm)
	ct.RevokedCertificateEntries = newRevokationListEntires(ct.revokedCertificateEntries())
	ct.Number = crl.Number
	ct.ThisUpdate = model.TimeDTO(crl.ThisUpdate)
	ct.NextUpdate = model.TimeDTO(crl.NextUpdate)
	return &ct
}
