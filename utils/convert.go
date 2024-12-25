package utils

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/model"
	"strings"
)

func DERToPem(der []byte, resourceType model.ResourceType) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  resourceType.PEMString(),
		Bytes: der,
	})
}

func PemToCertificate(pembytes []byte) (*x509.Certificate, error) {
	for {
		blk, rest := pem.Decode(pembytes)
		if blk == nil {
			break
		}
		if blk.Type ==
			"CERTIFICATE" {
			return x509.ParseCertificate(blk.Bytes)
		}
		pembytes = rest
	}
	return nil, fmt.Errorf("no certificate found")
}

func PemToCSR(pembytes []byte) (*x509.CertificateRequest, error) {
	for {
		blk, rest := pem.Decode(pembytes)
		if blk == nil {
			break
		}
		if blk.Type ==
			"CERTIFICATE REQUEST" {
			return x509.ParseCertificateRequest(blk.Bytes)
		}
		pembytes = rest
	}
	return nil, fmt.Errorf("no CSR found")
}

func PemToCRL(pembytes []byte) (*x509.RevocationList, error) {
	for {
		blk, rest := pem.Decode(pembytes)
		if blk == nil {
			break
		}
		if blk.Type ==
			"X509 CRL" {
			return x509.ParseRevocationList(blk.Bytes)
		}
		pembytes = rest
	}
	return nil, fmt.Errorf("no CRL found")
}

func X509ToTls(cert *x509.Certificate) (*tls.Certificate, error) {
	return &tls.Certificate{Certificate: [][]byte{cert.Raw}}, nil
}

func CapitaliseString(s string) string {
	if s == "" {
		return s
	}
	c := strings.ToUpper(s[0:1])
	if len(s) > 1 {
		c = strings.Join([]string{c, strings.ToLower(s[1:])}, "")
	}
	return c
}
