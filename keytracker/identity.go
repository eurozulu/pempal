package keytracker

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"pempal/pemreader"
)

// Identity represents a public key Identity, in the form of its Private key and a signed certificate sharing the same public key
// The certificate defines the usage of that key. Each identity shares the same key with any other identity signed by the same key
type Identity interface {
	fmt.Stringer

	// Key prepresents the private key which signed the Identity Certificate
	Key() Key

	// Certificate gets the key identity certificate
	Certificate() *x509.Certificate

	// Location is the location of the identity certificate
	Location() string

	// Usage checks if this identity can perform the given usages, based on the properties of its certificate.
	Usage(ku x509.KeyUsage, eku ...x509.ExtKeyUsage) bool
}

type identity struct {
	k    Key
	cert *x509.Certificate
	loc  string
}

func (id identity) String() string {
	n := id.cert.Subject.CommonName
	if n == "" {
		n = id.cert.Subject.String()
	}
	return n
}

func (id identity) Key() Key {
	return id.k
}

func (id identity) Certificate() *x509.Certificate {
	return id.cert
}

func (id identity) Location() string {
	return id.loc
}

func (id identity) Usage(ku x509.KeyUsage, eku ...x509.ExtKeyUsage) bool {
	if ku != 0 && id.cert.KeyUsage&ku != ku {
		return false
	}
	if len(eku) > 0 && !id.containsExtKeyUsages(eku) {
		return false
	}
	return true
}

func (id identity) containsExtKeyUsages(ekus []x509.ExtKeyUsage) bool {
	c := id.cert
	for _, eku := range ekus {
		if !containsExtKeyUsage(c, eku) {
			return false
		}
	}
	return true
}

func containsExtKeyUsage(c *x509.Certificate, eku x509.ExtKeyUsage) bool {
	for _, cku := range c.ExtKeyUsage {
		if cku == eku {
			return true
		}
	}
	return false
}

func NewIdentity(k Key, cert *pem.Block) (Identity, error) {
	// Check now so can skip when parsing later on
	c, err := x509.ParseCertificate(cert.Bytes)
	if err != nil {
		return nil, fmt.Errorf("invalid certificate pem creating new identity  %w", err)
	}
	loc, _ := readBlockHeader(pemreader.LocationHeaderKey, cert)
	return &identity{
		k:    k,
		cert: c,
		loc:  loc,
	}, nil
}
