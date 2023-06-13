package model

import (
	"crypto"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/utils"
	"math/big"
	"strings"
	"time"
)

type CertificateDTO struct {
	Id                    string               `yaml:"identity,omitempty"`
	Version               int                  `yaml:"version" flag:"version,ver"`
	SerialNumber          uint64               `yaml:"serial-number" flag:"serial-number,serialnumber,sn"`
	Signature             string               `yaml:"signature" flag:"signature,sig"`
	SignatureAlgorithm    string               `yaml:"signature-algorithm" flag:"signature-algorithm,signaturealgorithm,sig-algo"`
	PublicKeyAlgorithm    string               `yaml:"public-key-algorithm" flag:"public-key-algorithm,publickeyalgorithm,key-algorithm,keyalgorithm,keyalgo"`
	PublicKey             string               `yaml:"public-key" flag:"public-key,publickey,puk,pubkey"`
	Issuer                DistinguishedNameDTO `yaml:"issuer" flag:"issuer"`
	Subject               DistinguishedNameDTO `yaml:"subject" flag:"subject"`
	NotBefore             time.Time            `yaml:"not-before" flag:"not-before,notbefore,before"`
	NotAfter              time.Time            `yaml:"not-after" flag:"not-after,notafter,after"`
	IsCA                  bool                 `yaml:"is-ca" flag:"is-ca,isca"`
	BasicConstraintsValid bool                 `yaml:"basic-constraints-valid,omitempty" flag:"basic-constraints-valid,basicconstraintsvalid,constraints-valid,constraintsvalid"`
	MaxPathLen            int                  `yaml:"max-path-len,omitempty" flag:"max-path-len,maxpathlen,path-len,pathlen"`
	MaxPathLenZero        bool                 `yaml:"max-path-len-zero,omitempty" flag:"max-path-len-zero,maxpathlenzero,path-len-zero,pathlenzero"`
	KeyUsage              string               `yaml:"key-usage,omitempty" flag:"key-usage,keyusage"`

	ResourceType string `yaml:"resource-type" flag:"resource-type,resourcetype,type,rt"`
}

func (cd *CertificateDTO) UnmarshalPEM(data []byte) error {
	for len(data) > 0 {
		blk, rest := pem.Decode(data)
		if blk == nil {
			break
		}
		if ParsePEMType(blk.Type) != Certificate {
			data = rest
			continue
		}
		return cd.UnmarshalBinary(blk.Bytes)
	}
	return fmt.Errorf("no pem encoded certificate found")
}

func (cd *CertificateDTO) UnmarshalBinary(data []byte) error {
	cert, err := x509.ParseCertificate(data)
	if err != nil {
		return err
	}

	var puk string
	var id Identity
	if cert.PublicKey != nil {
		pukt, err := NewPublicKeyDTO(cert.PublicKey)
		if err != nil {
			return fmt.Errorf("Failed to parse certificates public key  %v", err)
		} else {
			puk = pukt.String()
		}
		id = Identity([]byte(puk))
		if err != nil {
			return err
		}
	}
	subject := newDistinguishedNameDTO(cert.Subject)
	issuer := newDistinguishedNameDTO(cert.Issuer)

	cd.Id = id.String()
	cd.Version = cert.Version
	cd.SerialNumber = cert.SerialNumber.Uint64()
	cd.Signature = hex.EncodeToString(cert.Signature)
	cd.SignatureAlgorithm = cert.SignatureAlgorithm.String()
	cd.PublicKeyAlgorithm = cert.PublicKeyAlgorithm.String()
	cd.PublicKey = puk
	cd.Issuer = *issuer
	cd.Subject = *subject
	cd.NotBefore = cert.NotBefore
	cd.NotAfter = cert.NotAfter
	cd.IsCA = cert.IsCA
	cd.BasicConstraintsValid = cert.BasicConstraintsValid
	cd.MaxPathLen = cert.MaxPathLen
	cd.MaxPathLenZero = cert.MaxPathLenZero
	cd.KeyUsage = strings.Join(utils.KeyUsageToStrings(cert.KeyUsage), "|")
	cd.ResourceType = Certificate.String()
	return nil
}

func (cd CertificateDTO) ToCertificate() (*x509.Certificate, error) {
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
			return nil, fmt.Errorf("failed to decode signature as hex  %v", err)
		}
		signature = by
	}
	var keyUsage x509.KeyUsage
	if cd.KeyUsage != "" {
		ku, err := utils.ParseKeyUsage(strings.Split(cd.KeyUsage, "|"))
		if err != nil {
			return nil, fmt.Errorf("failed to parse key usage  %v", err)
		}
		keyUsage = ku
	}
	return &x509.Certificate{
		Version:            cd.Version,
		SerialNumber:       new(big.Int).SetUint64(cd.SerialNumber),
		SignatureAlgorithm: utils.ParseSignatureAlgorithm(cd.SignatureAlgorithm),
		PublicKeyAlgorithm: utils.ParsePublicKeyAlgorithm(cd.PublicKeyAlgorithm),
		PublicKey:          puk,
		Issuer:             cd.Issuer.ToName(),
		Subject:            cd.Subject.ToName(),
		NotBefore:          cd.NotBefore,
		NotAfter:           cd.NotAfter,
		Signature:          signature,

		KeyUsage:                    keyUsage,
		Extensions:                  nil,
		ExtraExtensions:             nil,
		UnhandledCriticalExtensions: nil,
		ExtKeyUsage:                 nil,
		UnknownExtKeyUsage:          nil,

		BasicConstraintsValid:       cd.BasicConstraintsValid,
		IsCA:                        cd.IsCA,
		MaxPathLen:                  cd.MaxPathLen,
		MaxPathLenZero:              cd.MaxPathLenZero,
		SubjectKeyId:                nil,
		AuthorityKeyId:              nil,
		OCSPServer:                  nil,
		IssuingCertificateURL:       nil,
		DNSNames:                    nil,
		EmailAddresses:              nil,
		IPAddresses:                 nil,
		URIs:                        nil,
		PermittedDNSDomainsCritical: false,
		PermittedDNSDomains:         nil,
		ExcludedDNSDomains:          nil,
		PermittedIPRanges:           nil,
		ExcludedIPRanges:            nil,
		PermittedEmailAddresses:     nil,
		ExcludedEmailAddresses:      nil,
		PermittedURIDomains:         nil,
		ExcludedURIDomains:          nil,
		CRLDistributionPoints:       nil,
		PolicyIdentifiers:           nil,
	}, nil
}
