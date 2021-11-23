package keytools

const (
	PEM_CERTIFICATE         = "CERTIFICATE"
	PEM_X509_CERTIFICATE    = "X509 CERTIFICATE"
	PEM_TRUSTED_CERTIFICATE = "TRUSTED CERTIFICATE"
)

const (
	PEM_PRIVATE_KEY           = "PRIVATE KEY"
	PEM_ANY_PRIVATE_KEY       = "ANY PRIVATE KEY"
	PEM_ENCRYPTED_PRIVATE_KEY = "ENCRYPTED PRIVATE KEY"
	PEM_RSA_PRIVATE_KEY       = "RSA PRIVATE KEY"
	PEM_EC_PRIVATE_KEY        = "EC PRIVATE KEY"
	PEM_DSA_PRIVATE_KEY       = "DSA PRIVATE KEY"
)

const (
	PEM_PUBLIC_KEY       = "PUBLIC KEY"
	PEM_RSA_PUBLIC_KEY   = "RSA PUBLIC KEY"
	PEM_DSA_PUBLIC_KEY   = "DSA PUBLIC KEY"
	PEM_ECDSA_PUBLIC_KEY = "ECDSA PUBLIC KEY"
	PEM_EC_PUBLIC_KEY    = "EC PUBLIC KEY"
	PEM_ANY_PUBLIC_KEY   = "ANY_PUBLIC_KEY"
)

var CertificateTypes = map[string]bool{
	PEM_X509_CERTIFICATE:    true,
	PEM_CERTIFICATE:         true,
	PEM_TRUSTED_CERTIFICATE: true,
}
var PublicKeyTypes = map[string]bool{
	PEM_PUBLIC_KEY:       true,
	PEM_RSA_PUBLIC_KEY:   true,
	PEM_DSA_PUBLIC_KEY:   true,
	PEM_ECDSA_PUBLIC_KEY: true,
	PEM_EC_PUBLIC_KEY:    true,
	PEM_ANY_PUBLIC_KEY:   true,
}
var PrivateKeyTypes = map[string]bool{
	PEM_PRIVATE_KEY:           true,
	PEM_ANY_PRIVATE_KEY:       true,
	PEM_ENCRYPTED_PRIVATE_KEY: true,
	PEM_RSA_PRIVATE_KEY:       true,
	PEM_EC_PRIVATE_KEY:        true,
	PEM_DSA_PRIVATE_KEY:       true,
}

var CSRTypes = map[string]bool{
	"NEW CERTIFICATE REQUEST": true,
	"CERTIFICATE REQUEST":     true,
}
var CRLTypes = map[string]bool{
	"X509 CRL": true,
}

var PKCS7Types = map[string]bool{
	"PKCS7":               true,
	"PKCS #7 SIGNED DATA": true,
}
var ParamTypes = map[string]bool{
	"DH PARAMETERS":          true,
	"X9.42 DH PARAMETERS":    true,
	"SSL SESSION PARAMETERS": true,
	"DSA PARAMETERS":         true,
	"EC PARAMETERS":          true,
	"PARAMETERS":             true,
}
var CMSTypes = map[string]bool{
	"CMS": true,
}

func CombineMaps(ms ...map[string]bool) map[string]bool {
	m := map[string]bool{}
	for _, mm := range ms {
		for k, v := range mm {
			m[k] = v
		}
	}
	return m
}
