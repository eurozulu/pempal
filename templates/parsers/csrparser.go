package parsers

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"pempal/keytools"
	"strconv"
	"strings"
)

var AllCSRNames = []string{
	X509Signature,
	X509SignatureAlgorithm,
	X509PublicKey,
	X509PublicKeyAlgorithm,
	X509Version,
	X509Subject,
	X509Extensions,
	X509ExtraExtensions,
	X509DNSNames,
	X509EmailAddresses,
	X509IPAddresses,
	X509URIs,
}

type CsrParser struct{}

func (cp CsrParser) KnownNames() []string {
	return AllCSRNames
}

func (cp CsrParser) Parse(b *pem.Block) (map[string]interface{}, error) {
	csr, err := x509.ParseCertificateRequest(b.Bytes)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{
		PEM_TYPE: b.Type,
	}
	m[X509Signature] = base64.StdEncoding.EncodeToString(csr.Signature)
	m[X509SignatureAlgorithm] = csr.SignatureAlgorithm.String()

	m[X509PublicKey] = publicKeyToString(csr.PublicKey)
	m[X509PublicKeyAlgorithm] = csr.PublicKeyAlgorithm.String()
	m[X509PublicKeyHash] = keytools.PublicKeySha1Hash(csr.PublicKey)

	m[X509Version] = strconv.Itoa(csr.Version)

	m[X509Subject] = nameToStringMap(csr.Subject)

	m[X509Extensions] = extensionsToString(csr.Extensions)
	m[X509ExtraExtensions] = extensionsToString(csr.ExtraExtensions)

	m[X509DNSNames] = strings.Join(csr.DNSNames, ", ")
	m[X509EmailAddresses] = strings.Join(csr.EmailAddresses, ", ")
	m[X509IPAddresses] = iPsToString(csr.IPAddresses)
	m[X509URIs] = uRLsToString(csr.URIs)

	return m, nil
}
