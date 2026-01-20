package model

import (
	"crypto/x509"
	"fmt"
	"strings"
)

type PublicKeyAlgorithm x509.PublicKeyAlgorithm

var supportedAlgos = [...]x509.PublicKeyAlgorithm{
	x509.UnknownPublicKeyAlgorithm,
	x509.RSA,
	x509.DSA,
	x509.ECDSA,
	x509.Ed25519,
}
var algosName = [...]string{
	"UnknownPublicKeyAlgorithm",
	"RSA",
	"DSA",
	"ECDSA",
	"Ed25519",
}

func (pka PublicKeyAlgorithm) String() string {
	i := int(pka)
	if i < 0 || i >= len(algosName) {
		return ""
	}
	return algosName[i]
}

func (pka PublicKeyAlgorithm) MarshalText() (text []byte, err error) {
	return []byte(pka.String()), nil
}

func (pka *PublicKeyAlgorithm) UnmarshalText(text []byte) error {
	algo, err := ParsePublicKeyAlgorithm(string(text))
	if err != nil {
		return err
	}
	*pka = algo
	return nil
}

func ParsePublicKeyAlgorithm(s string) (PublicKeyAlgorithm, error) {
	for i, algo := range algosName {
		if !strings.EqualFold(algo, s) {
			continue
		}
		return PublicKeyAlgorithm(i), nil
	}
	return PublicKeyAlgorithm(0), fmt.Errorf("unknown public key algorithm: %q", s)
}
