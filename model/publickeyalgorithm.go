package model

import (
	"crypto/x509"
	"fmt"
	"strings"
)

type PublicKeyAlgorithm x509.PublicKeyAlgorithm

var publicKeyAlgoName = []string{
	"UnknownPublicKeyAlgorithm",
	"RSA",
	"DSA",
	"ECDSA",
	"Ed25519",
}

func (p PublicKeyAlgorithm) String() string {
	s, _ := p.MarshalText()
	return string(s)
}

func (p PublicKeyAlgorithm) MarshalText() (text []byte, err error) {
	i := int(p)
	if i < 0 || i >= len(publicKeyAlgoName) {
		i = 0
	}
	return []byte(publicKeyAlgoName[i]), nil
}

func (p *PublicKeyAlgorithm) UnmarshalText(text []byte) error {
	sz := string(text)
	for i, s := range publicKeyAlgoName {
		if !strings.EqualFold(s, sz) {
			continue
		}
		*p = PublicKeyAlgorithm(i)
		return nil
	}
	return fmt.Errorf("%s is an unknown public key algorithm", sz)
}
