package templates

import (
	"crypto/x509"
)

var (
	oidOrganization = []int{2, 5, 4, 10}
	oidCountry      = []int{2, 5, 4, 6}
)

func TemplateCA() Template {
	return &CertificateTemplate{
		KeyUsage:              KeyUsage(x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign),
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
		MaxPathLenZero:        false,
		SignatureAlgorithm:    SignatureAlgorithm(x509.SHA512WithRSA),
	}
}

func TemplateUser() Template {
	return &CertificateTemplate{
		KeyUsage:              KeyUsage(x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement),
		ExtKeyUsage:           []ExtKeyUsage{ExtKeyUsage(x509.ExtKeyUsageClientAuth)},
		BasicConstraintsValid: true,
		IsCA:                  false,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
		SignatureAlgorithm:    SignatureAlgorithm(x509.SHA512WithRSA),
	}
}

func TemplateServer() Template {
	return &CertificateTemplate{
		KeyUsage:              KeyUsage(x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement),
		ExtKeyUsage:           []ExtKeyUsage{ExtKeyUsage(x509.ExtKeyUsageServerAuth)},
		BasicConstraintsValid: true,
		IsCA:                  false,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
		SignatureAlgorithm:    SignatureAlgorithm(x509.SHA512WithRSA),
	}
}

func DefaultTemplate(name string) Template {
	switch name {
	case "ca" :
		return TemplateCA()
	case "server" :
		return TemplateServer()
	case "user":
		return TemplateUser()
	default:
		return nil
	}
}