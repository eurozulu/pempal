package identity

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/utils"
)

type Certificate interface {
	Location() string
	Certificate() *x509.Certificate
	String() string
	Identity() Identity
}

type certificate struct {
	loc  string
	cert *pem.Block
}

func (c certificate) Location() string {
	return c.loc
}

func (c certificate) Identity() Identity {
	cert := c.Certificate()
	if cert == nil || cert.PublicKey == nil {
		return ""
	}
	blk, err := utils.PublicKeyToPEM(cert.PublicKey)
	if err != nil || blk == nil {
		logger.Error("failed to encode public key  %v", err)
		return ""
	}
	return Identity(pem.EncodeToMemory(blk))
}

func (c certificate) Certificate() *x509.Certificate {
	if c.cert == nil {
		return nil
	}
	ct, err := x509.ParseCertificate(c.cert.Bytes)
	if err != nil {
		logger.Error("failed to parse certificate %v", err)
		return nil
	}
	return ct
}

// String returns a PEM encoded string of the certificate
func (c certificate) String() string {
	if c.cert == nil {
		return ""
	}
	return string(pem.EncodeToMemory(c.cert))
}
