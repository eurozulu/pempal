package keys

import (
	"crypto/rand"
	"crypto/x509"
	"fmt"
)

func SignCertificate(template *x509.Certificate, parent *x509.Certificate, key Key) (*x509.Certificate, error) {
	if IsEncrypted(key) {
		return nil, fmt.Errorf("key is encrypted")
	}
	if template.PublicKey == nil {
		return nil, fmt.Errorf("template has no public key")
	}
	by, err := x509.CreateCertificate(rand.Reader, template, parent, template.PublicKey, key.privateKey())
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(by)
}
