package model

import (
	"crypto"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/utils"
)

type CertificateRequestDTO struct {
	Id                 string               `yaml:"identity"`
	Version            int                  `yaml:"version" flag:"version,ver"`
	Signature          string               `yaml:"signature" flag:"signature,sig"`
	SignatureAlgorithm string               `yaml:"signature-algorithm" flag:"signature-algorithm,signaturealgorithm,sig-algo"`
	PublicKeyAlgorithm string               `yaml:"public-key-algorithm" flag:"public-key-algorithm,publickeyalgorithm,key-algorithm,keyalgorithm,keyalgo"`
	PublicKey          string               `yaml:"public-key" flag:"public-key,publickey,puk,pubkey"`
	Subject            DistinguishedNameDTO `yaml:"subject" flag:"subject"`

	ResourceType string `yaml:"resource-type" flag:"resource-type,resourcetype,type,rt"`
}

func (cd *CertificateRequestDTO) UnmarshalPEM(data []byte) error {
	for len(data) > 0 {
		blk, rest := pem.Decode(data)
		if blk == nil {
			break
		}
		if ParsePEMType(blk.Type) != CertificateRequest {
			data = rest
			continue
		}
		return cd.UnmarshalBinary(blk.Bytes)
	}
	return fmt.Errorf("no pem encoded certificate found")
}

func (cd *CertificateRequestDTO) UnmarshalBinary(data []byte) error {
	csr, err := x509.ParseCertificateRequest(data)
	if err != nil {
		return err
	}

	var puk string
	var id Identity
	if csr.PublicKey != nil {
		pukt, err := NewPublicKeyDTO(csr.PublicKey)
		if err != nil {
			return fmt.Errorf("Failed to parse certificate requests public key  %v", err)
		} else {
			puk = pukt.String()
		}
		id = Identity([]byte(puk))
		if err != nil {
			return err
		}
	}

	subject := newDistinguishedNameDTO(csr.Subject)

	cd.Id = id.String()
	cd.Version = csr.Version
	cd.Signature = hex.EncodeToString(csr.Signature)
	cd.SignatureAlgorithm = csr.SignatureAlgorithm.String()
	cd.PublicKeyAlgorithm = csr.PublicKeyAlgorithm.String()
	cd.PublicKey = puk
	cd.Subject = *subject

	cd.ResourceType = CertificateRequest.String()
	return nil
}

func (cd CertificateRequestDTO) ToCertificateRequest() (*x509.CertificateRequest, error) {
	var puk crypto.PublicKey
	if cd.PublicKey != "" {
		pkdto := &PublicKeyDTO{}
		if err := pkdto.UnmarshalPEM([]byte(cd.PublicKey)); err == nil {
			puk, err = pkdto.ToPublicKey()
			if err != nil {
				return nil, err
			}
		} else {
			logger.Warning("certificate public key failed to parse %v", err)
		}
	}

	var signature []byte
	if cd.Signature != "" {
		by, err := hex.DecodeString(cd.Signature)
		if err != nil {
			return nil, fmt.Errorf("failed to decodfe signature as hex  %v", err)
		}
		signature = by
	}
	return &x509.CertificateRequest{
		Version:            cd.Version,
		SignatureAlgorithm: utils.ParseSignatureAlgorithm(cd.SignatureAlgorithm),
		PublicKeyAlgorithm: utils.ParsePublicKeyAlgorithm(cd.PublicKeyAlgorithm),
		PublicKey:          puk,
		Subject:            cd.Subject.ToName(),
		Signature:          signature,

		Extensions:      nil,
		ExtraExtensions: nil,

		DNSNames:       nil,
		EmailAddresses: nil,
		IPAddresses:    nil,
		URIs:           nil,
	}, nil
}
