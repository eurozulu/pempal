package utils

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/x509"
	"encoding/hex"
	"fmt"
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

func SignatureAlgorithmNames() []string {
	names := make([]string, len(signatureAlgorithmDetails))
	for i, sigDetail := range signatureAlgorithmDetails {
		names[i] = sigDetail.name
	}
	return names
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

var keyusages = []x509.KeyUsage{
	x509.KeyUsageDigitalSignature,
	x509.KeyUsageContentCommitment,
	x509.KeyUsageKeyEncipherment,
	x509.KeyUsageDataEncipherment,
	x509.KeyUsageKeyAgreement,
	x509.KeyUsageCertSign,
	x509.KeyUsageCRLSign,
	x509.KeyUsageEncipherOnly,
	x509.KeyUsageDecipherOnly,
}
var keyusageNames = []string{
	"KeyUsageDigitalSignature",
	"KeyUsageContentCommitment",
	"KeyUsageKeyEncipherment",
	"KeyUsageDataEncipherment",
	"KeyUsageKeyAgreement",
	"KeyUsageCertSign",
	"KeyUsageCRLSign",
	"KeyUsageEncipherOnly",
	"KeyUsageDecipherOnly",
}

func ParseKeyUsage(s []string) (x509.KeyUsage, error) {
	var k x509.KeyUsage
	for _, kun := range s {
		ku := lookupKeyUsage(kun)
		if ku == 0 {
			return 0, fmt.Errorf("%s is not a known key usage", kun)
		}
		k |= ku
	}
	return k, nil
}

func lookupKeyUsage(s string) x509.KeyUsage {
	for i, kun := range keyusageNames {
		if strings.EqualFold(s, kun) {
			return keyusages[i]
		}
	}
	return 0
}

func KeyUsageToStrings(k x509.KeyUsage) []string {
	var names []string
	for _, ku := range keyusages {
		if k&ku != ku {
			continue
		}
		names = append(names)
	}
	return names
}

var ExtKeyUsageName = []string{
	"Any",
	"ServerAuth",
	"ClientAuth",
	"CodeSigning",
	"EmailProtection",
	"IPSECEndSystem",
	"IPSECTunnel",
	"IPSECUser",
	"TimeStamping",
	"OCSPSigning",
	"MicrosoftServerGatedCrypto",
	"NetscapeServerGatedCrypto",
	"MicrosoftCommercialCodeSigning",
	"MicrosoftKernelCodeSigning",
}

func ParseExtKeyUsage(eks []string) ([]x509.ExtKeyUsage, error) {
	var found []x509.ExtKeyUsage
	for _, s := range eks {
		ku := lookupExtKeyUsage(s)
		if ku == 0 {
			return nil, fmt.Errorf("%s is not a known extended key usage", s)
		}
		found = append(found, ku)
	}
	return found, nil
}

func ExtKeyUsageToStrings(ek []x509.ExtKeyUsage) []string {
	names := make([]string, len(ek))
	for i, e := range ek {
		if e < 0 || int(e) >= len(ExtKeyUsageName) {
			e = 0
		}
		names[i] = ExtKeyUsageName[e]
	}
	return names
}

func lookupExtKeyUsage(s string) x509.ExtKeyUsage {
	for i, kun := range ExtKeyUsageName {
		if strings.EqualFold(s, kun) {
			return x509.ExtKeyUsage(i)
		}
	}
	return 0
}
