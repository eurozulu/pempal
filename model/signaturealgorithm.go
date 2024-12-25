package model

import (
	"crypto/x509"
	"fmt"
	"strings"
)

type SignatureAlgorithm x509.SignatureAlgorithm

var signatureNames = []string{
	"UnknownSignatureAlgorithm",
	"MD2-RSA",
	"MD5-RSA",
	"SHA1-RSA",
	"SHA256-RSA",
	"SHA384-RSA",
	"SHA512-RSA",
	"DSA-SHA1",
	"DSA-SHA256",
	"ECDSA-SHA1",
	"ECDSA-SHA256",
	"ECDSA-SHA384",
	"ECDSA-SHA512",
	"SHA256-RSAPSS",
	"SHA384-RSAPSS",
	"SHA512-RSAPSS",
	"Ed25519",
}

func (s SignatureAlgorithm) MarshalText() (text []byte, err error) {
	i := int(s)
	if i < 0 || i >= len(signatureNames) {
		i = 0
	}
	return []byte(signatureNames[i]), nil
}

func (s *SignatureAlgorithm) UnmarshalText(text []byte) error {
	sz := string(text)
	for i, n := range signatureNames {
		if !strings.EqualFold(n, sz) {
			continue
		}
		(*s) = SignatureAlgorithm(i)
		return nil
	}
	return fmt.Errorf("%s is an unknown signature algorithm", sz)
}
