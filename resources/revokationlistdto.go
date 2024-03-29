package resources

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
)

var errNoRevokationList = fmt.Errorf("no pem encoded revokation list found")

type RevocationListDTO struct {
	Issuer             string `yaml:"issuer"`
	Signature          string `yaml:"signature"`
	SignatureAlgorithm string `yaml:"signature-algorithm"`

	// RevokedCertificates is used to populate the revokedCertificates
	// sequence in the CRL, it may be empty. RevokedCertificates may be nil,
	// in which case an empty CRL will be created.
	RevokedCertificates []string `yaml:"revoked-certificates" json:"revoked-certificates"`
	Number              int64    `yaml:"number" json:"number"`
	ThisUpdate          string   `yaml:"this-update" json:"this-update"`
	NextUpdate          string   `yaml:"next-update" json:"next-update"`
	Extensions          []string `yaml:"extensions,omitempty" json:"extensions"`

	// ExtraExtensions contains any additional extensions to add directly to
	// the CRL.
	ExtraExtensions []string `yaml:"extra-extensions,omitempty"`
	RevokationList  string   `yaml:"revokation-list,omitempty" json:"-"`
}

func (rvl *RevocationListDTO) String() string {
	return rvl.RevokationList
}

func (rvl *RevocationListDTO) UnmarshalPEM(data []byte) error {
	blk, _ := pem.Decode(data)
	if blk == nil {
		return fmt.Errorf("failed to parse revokation list key from pem")
	}
	return rvl.UnmarshalBinary(blk.Bytes)
}

func (rvl *RevocationListDTO) MarshalBinary() (data []byte, err error) {
	blk, _ := pem.Decode([]byte(rvl.RevokationList))
	if blk == nil {
		return nil, errNoRevokationList
	}
	return blk.Bytes, nil
}

func (rvl *RevocationListDTO) UnmarshalBinary(data []byte) error {
	rlist, err := x509.ParseRevocationList(data)
	if err != nil {
		return err
	}
	rvl.RevokationList = string(pem.EncodeToMemory(&pem.Block{
		Type:  RevocationList.PEMString(),
		Bytes: rlist.Raw,
	}))
	rvl.Issuer = rlist.Issuer.String()
	rvl.Signature = hex.EncodeToString(rlist.Signature)
	rvl.SignatureAlgorithm = rlist.SignatureAlgorithm.String()

	rvl.RevokedCertificates = nil //newRevokedCertificateDTOs(rlist.RevokedCertificates)
	rvl.Number = rlist.Number.Int64()
	rvl.ThisUpdate = rlist.ThisUpdate.Format(CommonDateFormat)
	rvl.NextUpdate = rlist.NextUpdate.Format(CommonDateFormat)
	rvl.Extensions = nil      //newExtentionsDTOs(rlist.Extensions)
	rvl.ExtraExtensions = nil //newExtentionsDTOs(rlist.ExtraExtensions)

	return nil
}
