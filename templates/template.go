package templates

import (
	"crypto/sha256"
	"encoding"
	"fmt"
	"strings"
)


type Template interface {
	String() string
	Location() string
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

func NewTemplate(p string, tt string) (Template, error) {
	switch strings.ToUpper(tt) {
	case "CERTIFICATE":
		return &CertificateTemplate{FilePath: p}, nil
	case "CERTIFICATE REQUEST":
		return &CSRTemplate{FilePath: p}, nil
	case "PUBLIC KEY":
		return &PublicKeyTemplate{FilePath: p}, nil
	case "PRIVATE KEY":
		return &PrivateKeyTemplate{FilePath: p}, nil
	case "X509 CRL":
		return &CRLTemplate{FilePath: p}, nil
	default:
		return nil, fmt.Errorf("%s is an unknown template type\n", tt)
	}
}

func TemplateType(t Template) string {
	switch t.(type) {
	case *CertificateTemplate:
		return "CERTIFICATE"
	case *CSRTemplate:
		return "CERTIFICATE REQUEST"
	case *PublicKeyTemplate:
		return "PUBLIC KEY"
	case *PrivateKeyTemplate:
		return "PRIVATE KEY"
	case *CRLTemplate:
		return "X509 CRL"
	case *PKCS12Template:
		return "PKCS12"
	default:
		return ""
	}
}

func fingerprint(by []byte) string {
	h := sha256.New()
	_, _ = h.Write(by)
	return string(h.Sum(nil))
}
