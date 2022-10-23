package main

import (
	"crypto/x509"
	"io"
	"pempal/keys"
)

type genKeyCommand struct {
	KeyAlgorithm string `flag:"keyalgorithm"`
	Ka           string `flag:"ka"`
	Size         string `flag:"size"`
}

func (g genKeyCommand) Run(out io.Writer, args ...string) error {
	k, err := keys.MakeKey(g.keyAlgorithm(), g.Size)
	if err != nil {
		return err
	}
	k.WriteKey()
}

func newGenKeyCommand() Command {
	return &genKeyCommand{}
}

func (g genKeyCommand) keyAlgorithm() x509.PublicKeyAlgorithm {
	s := g.KeyAlgorithm
	if s == "" {
		s = g.Ka
	}
	switch s {
	case "RSA":
		return x509.RSA
	case "ECDSA":
		return x509.ECDSA
	case "Ed25519":
		return x509.Ed25519
	case "DSA":
		return x509.DSA
	default:
		return x509.UnknownPublicKeyAlgorithm
	}
}
