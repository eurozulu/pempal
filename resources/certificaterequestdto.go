package resources

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/utils"
)

var errNoCertificateRequest = fmt.Errorf("no pem encoded certificate request found")

type CertificateRequestDTO struct {
	Version            int    `yaml:"version" flag:"version,ver"`
	Signature          string `yaml:"signature" flag:"signature,sig"`
	SignatureAlgorithm string `yaml:"signature-algorithm" flag:"signature-algorithm,signaturealgorithm,sig-algo"`
	PublicKeyAlgorithm string `yaml:"public-key-algorithm" flag:"public-key-algorithm,publickeyalgorithm,key-algorithm,keyalgorithm,keyalgo"`
	PublicKey          string `yaml:"public-key" flag:"public-key,publickey,puk,pubkey"`
	Subject            string `yaml:"subject" flag:"subject"`

	CertificateRequest string `yaml:"certificate-request,omitempty"`
}

func (crd CertificateRequestDTO) String() string {
	return crd.CertificateRequest
}

func (crd CertificateRequestDTO) ToCertificateRequest() (*x509.CertificateRequest, error) {
	subject, err := stringToDN(crd.Subject)
	if err != nil {
		logger.Debug("failed to read certificate request subject  %v", err)
	}

	var signature []byte
	if crd.Signature != "" {
		by, err := hex.DecodeString(crd.Signature)
		if err != nil {
			return nil, fmt.Errorf("failed to decodfe signature as hex  %v", err)
		}
		signature = by
	}
	return &x509.CertificateRequest{
		Version:            crd.Version,
		SignatureAlgorithm: utils.ParseSignatureAlgorithm(crd.SignatureAlgorithm),
		PublicKeyAlgorithm: utils.ParsePublicKeyAlgorithm(crd.PublicKeyAlgorithm),
		PublicKey:          stringToPublicKey(crd.PublicKey),
		Subject:            subject,
		Signature:          signature,

		Extensions:      nil,
		ExtraExtensions: nil,

		DNSNames:       nil,
		EmailAddresses: nil,
		IPAddresses:    nil,
		URIs:           nil,
	}, nil
}

func (crd *CertificateRequestDTO) UnmarshalBinary(data []byte) error {
	csr, err := x509.ParseCertificateRequest(data)
	if err != nil {
		return err
	}
	crd.CertificateRequest = string(pem.EncodeToMemory(&pem.Block{
		Type:  CertificateRequest.PEMString(),
		Bytes: csr.Raw,
	}))
	var puk string
	if csr.PublicKey != nil {
		pukt, err := NewPublicKeyDTO(csr.PublicKey)
		if err != nil {
			return fmt.Errorf("Failed to parse certificate requests public key  %v", err)
		} else {
			puk = pukt.String()
		}
	}

	var subject string
	if csr.Subject.String() != "" {
		dto, err := ParseDistinguishedName(csr.Subject.String())
		if err != nil {
			return fmt.Errorf("failed to parse subject  %v", err)
		}
		subject = dto.String()
	}

	crd.Version = csr.Version
	crd.Signature = hex.EncodeToString(csr.Signature)
	crd.SignatureAlgorithm = csr.SignatureAlgorithm.String()
	crd.PublicKeyAlgorithm = csr.PublicKeyAlgorithm.String()
	crd.PublicKey = puk
	crd.Subject = subject
	return nil
}
func (crd *CertificateRequestDTO) MarshalBinary() (data []byte, err error) {
	blk, _ := pem.Decode([]byte(crd.CertificateRequest))
	if blk == nil {
		return nil, errNoCertificateRequest
	}
	return blk.Bytes, nil
}
