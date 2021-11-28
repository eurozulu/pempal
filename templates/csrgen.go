package templates

import (
	"crypto/rand"
	"crypto/x509"
	"pempal/keytracker"
	"pempal/templates/parsers"
)

// MakeCSR copies the given template values into a new Certificate Request.
// The given signers public key and PublicKeyAlgorithm are copied into the csr. the CSR is NOT signed
func MakeCSR(signer keytracker.Key, t Template) *x509.CertificateRequest {
	csr := &x509.CertificateRequest{
		Signature:          stringToBytes(t.Value(parsers.X509Signature)),
		SignatureAlgorithm: stringToSignatureAlgorithm(t.Value(parsers.X509SignatureAlgorithm)),
		Version:            stringToInt(t.Value(parsers.X509Version)),
		Subject:            readNameProperties(t.ValueMap(parsers.X509Subject)),
		Extensions:         stringToExtensions(t.Value(parsers.X509Extensions)),
		ExtraExtensions:    stringToExtensions(t.Value(parsers.X509ExtraExtensions)),
		DNSNames:           stringToStringArray(t.Value(parsers.X509DNSNames)),
		EmailAddresses:     stringToStringArray(t.Value(parsers.X509EmailAddresses)),
		IPAddresses:        stringToIPs(t.Value(parsers.X509IPAddresses)),
		URIs:               stringToURLs(t.Value(parsers.X509URIs)),
	}
	csr.PublicKey = signer.PublicKey()
	csr.PublicKeyAlgorithm = signer.PublicKeyAlgorithm()
	return csr
}

func SignRequest(signer keytracker.Key, keypass string, t Template) ([]byte, error) {
	csr := MakeCSR(signer, t)
	prk, err := signer.PrivateKeyDecrypted(keypass)
	if err != nil {
		return nil, err
	}
	return x509.CreateCertificateRequest(rand.Reader, csr, prk)
}
