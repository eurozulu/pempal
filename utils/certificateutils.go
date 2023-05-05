package utils

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/x509"
	"encoding/hex"
	"strings"
)

func CertificateFingerprint(c *x509.Certificate) string {
	buf := bytes.NewBuffer(nil)
	for i, b := range md5.Sum(c.Raw) {
		if i > 0 {
			buf.WriteRune(':')
		}
		buf.WriteString(hex.EncodeToString([]byte{b}))
	}
	return buf.String()
}

func IsRootCertificate(c *x509.Certificate) bool {
	return c != nil && c.IsCA && (c.Subject.String() == c.Issuer.String())
}

func ParseSignatureAlgorithm(s string) x509.SignatureAlgorithm {
	for _, sigDetail := range signatureAlgorithmDetails {
		if strings.EqualFold(s, sigDetail.name) {
			return sigDetail.algo
		}
	}
	return x509.UnknownSignatureAlgorithm
}

// copied from x509 package
var signatureAlgorithmDetails = []struct {
	algo       x509.SignatureAlgorithm
	name       string
	pubKeyAlgo x509.PublicKeyAlgorithm
	hash       crypto.Hash
}{
	{x509.MD2WithRSA, "MD2-RSA", x509.RSA, crypto.Hash(0) /* no value for MD2 */},
	{x509.MD5WithRSA, "MD5-RSA", x509.RSA, crypto.MD5},
	{x509.SHA1WithRSA, "SHA1-RSA", x509.RSA, crypto.SHA1},
	{x509.SHA1WithRSA, "SHA1-RSA", x509.RSA, crypto.SHA1},
	{x509.SHA256WithRSA, "SHA256-RSA", x509.RSA, crypto.SHA256},
	{x509.SHA384WithRSA, "SHA384-RSA", x509.RSA, crypto.SHA384},
	{x509.SHA512WithRSA, "SHA512-RSA", x509.RSA, crypto.SHA512},
	{x509.SHA256WithRSAPSS, "SHA256-RSAPSS", x509.RSA, crypto.SHA256},
	{x509.SHA384WithRSAPSS, "SHA384-RSAPSS", x509.RSA, crypto.SHA384},
	{x509.SHA512WithRSAPSS, "SHA512-RSAPSS", x509.RSA, crypto.SHA512},
	{x509.DSAWithSHA1, "DSA-SHA1", x509.DSA, crypto.SHA1},
	{x509.DSAWithSHA256, "DSA-SHA256", x509.DSA, crypto.SHA256},
	{x509.ECDSAWithSHA1, "ECDSA-SHA1", x509.ECDSA, crypto.SHA1},
	{x509.ECDSAWithSHA256, "ECDSA-SHA256", x509.ECDSA, crypto.SHA256},
	{x509.ECDSAWithSHA384, "ECDSA-SHA384", x509.ECDSA, crypto.SHA384},
	{x509.ECDSAWithSHA512, "ECDSA-SHA512", x509.ECDSA, crypto.SHA512},
	{x509.PureEd25519, "Ed25519", x509.Ed25519, crypto.Hash(0) /* no pre-hashing */},
}
