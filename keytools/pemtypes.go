package keytools

const (
	PEM_CERTIFICATE         = "CERTIFICATE"
	PEM_X509_CERTIFICATE    = "X509 CERTIFICATE"
	PEM_TRUSTED_CERTIFICATE = "TRUSTED CERTIFICATE"
)

const (
	PEM_NEW_CERTIFICATE_REQUEST = "NEW CERTIFICATE REQUEST"
	PEM_CERTIFICATE_REQUEST     = "CERTIFICATE REQUEST"
)

const PEM_X509_CRL = "X509 CRL"

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

// CertificateTypes is a map of all the pem types for a certificate.
var CertificateTypes = map[string]bool{
	PEM_X509_CERTIFICATE:    true,
	PEM_CERTIFICATE:         true,
	PEM_TRUSTED_CERTIFICATE: true,
}

// PublicKeyTypes is a map of all the pem types for a Public Key.
var PublicKeyTypes = map[string]bool{
	PEM_PUBLIC_KEY:       true,
	PEM_RSA_PUBLIC_KEY:   true,
	PEM_DSA_PUBLIC_KEY:   true,
	PEM_ECDSA_PUBLIC_KEY: true,
	PEM_EC_PUBLIC_KEY:    true,
	PEM_ANY_PUBLIC_KEY:   true,
}

// PrivateKeyTypes is a map of all the pem types for a Private Key.
var PrivateKeyTypes = map[string]bool{
	PEM_PRIVATE_KEY:           true,
	PEM_ANY_PRIVATE_KEY:       true,
	PEM_ENCRYPTED_PRIVATE_KEY: true,
	PEM_RSA_PRIVATE_KEY:       true,
	PEM_EC_PRIVATE_KEY:        true,
	PEM_DSA_PRIVATE_KEY:       true,
}

// CSRTypes is a map of all the pem types for a Certificate Request.
var CSRTypes = map[string]bool{
	PEM_NEW_CERTIFICATE_REQUEST: true,
	PEM_CERTIFICATE_REQUEST:     true,
}

// CRLTypes is a map of all the pem types for a Certificate Revokation List.
var CRLTypes = map[string]bool{
	PEM_X509_CRL: true,
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

// CombineMaps takes one or more string/bool maps and merges them into one.
// Used to combine types of PEMS into a single map
// e.g. CombineMaps(PrivateKeyTypes, PublicKeyTypes) results in a map which will find all key pems
func CombineMaps(ms ...map[string]bool) map[string]bool {
	m := map[string]bool{}
	for _, mm := range ms {
		for k, v := range mm {
			m[k] = v
		}
	}
	return m
}
