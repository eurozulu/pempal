package resources

import (
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/utils"
	"math/big"
	"time"
)

var CommonDateFormat = time.RFC850
var errNoCertificate = fmt.Errorf("no pem encoded certificate found")

type CertificateDTO struct {
	Version               int      `yaml:"version" json:"version"`
	SerialNumber          int64    `yaml:"serial-number" json:"serial-number"`
	Signature             string   `yaml:"signature" json:"signature"`
	SignatureAlgorithm    string   `yaml:"signature-algorithm" json:"signature-algorithm"`
	PublicKeyAlgorithm    string   `yaml:"public-key-algorithm" json:"public-key-algorithm"`
	PublicKey             string   `yaml:"public-key" json:"public-key"`
	Issuer                string   `yaml:"issuer" json:"issuer"`
	Subject               string   `yaml:"subject" json:"subject"`
	NotBefore             string   `yaml:"not-before" json:"not-before"`
	NotAfter              string   `yaml:"not-after" json:"not-after"`
	IsCA                  bool     `yaml:"is-ca" json:"is-ca"`
	BasicConstraintsValid bool     `yaml:"basic-constraints-valid,omitempty" json:"basic-constraints-valid"`
	MaxPathLen            int      `yaml:"max-path-len,omitempty" json:"max-path-len,omitempty"`
	MaxPathLenZero        bool     `yaml:"max-path-len-zero,omitempty" json:"max-path-len-zero"`
	KeyUsage              []string `yaml:"key-usage,omitempty" json:"key-usage"`
	ExtKeyUsage           []string `yaml:"extended-key-usage,omitempty" json:"extended-key-usage"`

	Certificate string `yaml:"certificate,omitempty" json:"-"`
}

func (c CertificateDTO) ToCertificate() *x509.Certificate {
	sig, err := hex.DecodeString(c.Signature)
	if err != nil {
		logger.Debug("signature invalid %v", err)
	}
	notBefore, err := time.Parse(CommonDateFormat, c.NotBefore)
	if err != nil {
		logger.Debug("failed to read certificate notBefore time  %v", err)
	}
	notAfter, err := time.Parse(CommonDateFormat, c.NotAfter)
	if err != nil {
		logger.Debug("failed to read certificate notAfter time  %v", err)
	}
	ku, err := utils.ParseKeyUsage(c.KeyUsage)
	if err != nil {
		logger.Debug("failed to read certificate KeyUsage  %v", err)
	}
	eku, err := utils.ParseExtKeyUsage(c.ExtKeyUsage)
	if err != nil {
		logger.Debug("failed to read certificate Extended KeyUsage  %v", err)
	}

	issuer, err := stringToDN(c.Issuer)
	if err != nil {
		logger.Debug("failed to read certificate issuer  %v", err)
	}
	subject, err := stringToDN(c.Subject)
	if err != nil {
		logger.Debug("failed to read certificate subject  %v", err)
	}

	return &x509.Certificate{
		Signature:                   sig,
		SignatureAlgorithm:          utils.ParseSignatureAlgorithm(c.SignatureAlgorithm),
		PublicKeyAlgorithm:          utils.ParsePublicKeyAlgorithm(c.PublicKeyAlgorithm),
		PublicKey:                   stringToPublicKey(c.PublicKey),
		Version:                     c.Version,
		SerialNumber:                big.NewInt(c.SerialNumber),
		Issuer:                      issuer,
		Subject:                     subject,
		NotBefore:                   notBefore,
		NotAfter:                    notAfter,
		BasicConstraintsValid:       c.BasicConstraintsValid,
		IsCA:                        c.IsCA,
		MaxPathLen:                  c.MaxPathLen,
		MaxPathLenZero:              c.MaxPathLenZero,
		KeyUsage:                    ku,
		ExtKeyUsage:                 eku,
		Extensions:                  nil,
		ExtraExtensions:             nil,
		UnknownExtKeyUsage:          nil,
		SubjectKeyId:                nil,
		AuthorityKeyId:              nil,
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
		OCSPServer:                  nil,
		IssuingCertificateURL:       nil,
	}
}

func (c *CertificateDTO) UnmarshalBinary(data []byte) error {
	cer, err := x509.ParseCertificate(data)
	if err != nil {
		return err
	}
	c.Certificate = string(pem.EncodeToMemory(&pem.Block{
		Type:  Certificate.PEMString(),
		Bytes: cer.Raw,
	}))
	c.Version = cer.Version
	c.SerialNumber = cer.SerialNumber.Int64()
	c.Signature = hex.EncodeToString(cer.Signature)
	c.SignatureAlgorithm = cer.SignatureAlgorithm.String()
	c.PublicKeyAlgorithm = cer.PublicKeyAlgorithm.String()
	c.PublicKey = publickKeyToString(cer.PublicKey)
	c.Issuer = cer.Issuer.String()
	c.Subject = cer.Subject.String()
	c.NotAfter = cer.NotAfter.Format(CommonDateFormat)
	c.NotBefore = cer.NotBefore.Format(CommonDateFormat)
	c.IsCA = cer.IsCA
	c.BasicConstraintsValid = cer.BasicConstraintsValid
	c.MaxPathLen = cer.MaxPathLen
	c.MaxPathLenZero = cer.MaxPathLenZero
	c.KeyUsage = utils.KeyUsageToStrings(cer.KeyUsage)
	c.ExtKeyUsage = utils.ExtKeyUsageToStrings(cer.ExtKeyUsage)
	return nil
}

func (c CertificateDTO) MarshalBinary() (data []byte, err error) {
	blk, _ := pem.Decode([]byte(c.Certificate))
	if blk == nil {
		return nil, errNoCertificate
	}
	return blk.Bytes, nil
}

func (c CertificateDTO) String() string {
	return c.Certificate
}

func publickKeyToString(puk crypto.PublicKey) string {
	pukBlk, err := utils.PublicKeyToPEM(puk)
	if err != nil || pukBlk == nil {
		logger.Debug("Failed to parse certificate public key %v", err)
		return ""
	}
	return string(pem.EncodeToMemory(pukBlk))
}

func stringToPublicKey(s string) crypto.PublicKey {
	blk, _ := pem.Decode([]byte(s))
	if blk == nil {
		return nil
	}
	puk, err := x509.ParsePKIXPublicKey(blk.Bytes)
	if err != nil {
		logger.Debug("Failed to parse certificate public key %v", err)
		return nil
	}
	return puk
}

func stringToDN(s string) (pkix.Name, error) {
	if s == "" {
		return pkix.Name{}, nil
	}
	d := &DistinguishedNameDTO{}
	if err := d.UnmarshalBinary([]byte(s)); err != nil {
		return pkix.Name{}, err
	}
	return d.ToName(), nil
}
