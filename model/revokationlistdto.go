package model

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"time"
)

type RevocationListDTO struct {
	Issuer             DistinguishedNameDTO `yaml:"issuer"`
	Signature          string               `yaml:"signature"`
	SignatureAlgorithm string               `yaml:"signature-algorithm"`

	// RevokedCertificates is used to populate the revokedCertificates
	// sequence in the CRL, it may be empty. RevokedCertificates may be nil,
	// in which case an empty CRL will be created.
	RevokedCertificates []RevokedCertificateDTO `yaml:"revoked-certificates"`
	Number              uint64                  `yaml:"number"`

	ThisUpdate time.Time      `yaml:"this-update"`
	NextUpdate time.Time      `yaml:"next-update"`
	Extensions []ExtensionDTO `yaml:"extensions"`

	// ExtraExtensions contains any additional extensions to add directly to
	// the CRL.
	ExtraExtensions []ExtensionDTO `yaml:"extraExtensions"`

	ResourceType string `yaml:"resource-type"`
}

func (rvl *RevocationListDTO) UnmarshalPEM(data []byte) error {
	for len(data) > 0 {
		blk, rest := pem.Decode(data)
		if blk == nil {
			break
		}
		if ParsePEMType(blk.Type) != RevokationList {
			data = rest
			continue
		}
		return rvl.UnmarshalBinary(blk.Bytes)
	}
	return fmt.Errorf("no pem encoded public key found")
}

func (rvl *RevocationListDTO) UnmarshalBinary(data []byte) error {
	rlist, err := x509.ParseRevocationList(data)
	if err != nil {
		return err
	}

	issuer := newDistinguishedNameDTO(rlist.Issuer)
	rvl.Issuer = *issuer
	rvl.Signature = hex.EncodeToString(rlist.Signature)
	rvl.SignatureAlgorithm = rlist.SignatureAlgorithm.String()

	rvl.RevokedCertificates = newRevokedCertificateDTOs(rlist.RevokedCertificates)
	rvl.Number = rlist.Number.Uint64()
	rvl.ThisUpdate = rlist.ThisUpdate
	rvl.NextUpdate = rlist.NextUpdate
	rvl.ExtraExtensions = newExtentionsDTOs(rlist.Extensions)
	rvl.ExtraExtensions = newExtentionsDTOs(rlist.ExtraExtensions)

	rvl.ResourceType = RevokationList.String()
	return nil
}
