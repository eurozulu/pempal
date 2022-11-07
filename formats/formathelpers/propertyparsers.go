package formathelpers

import (
	"crypto"
	"crypto/elliptic"
	"crypto/x509"
	"net"
	"net/url"
	"strings"
	"time"
)

const TimeFormat = time.RFC850

var keyUsageNames = [...]string{
	"",
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

var ExtKeyUsageNames = [...]string{
	"ExtKeyUsageAny",
	"ExtKeyUsageServerAuth",
	"ExtKeyUsageClientAuth",
	"ExtKeyUsageCodeSigning",
	"ExtKeyUsageEmailProtection",
	"ExtKeyUsageIPSECEndSystem",
	"ExtKeyUsageIPSECTunnel",
	"ExtKeyUsageIPSECUser",
	"ExtKeyUsageTimeStamping",
	"ExtKeyUsageOCSPSigning",
	"ExtKeyUsageMicrosoftServerGatedCrypto",
	"ExtKeyUsageNetscapeServerGatedCrypto",
	"ExtKeyUsageMicrosoftCommercialCodeSigning",
	"ExtKeyUsageMicrosoftKernelCodeSigning",
}

func KeyUsageString(k x509.KeyUsage) string {
	if k < 1 || int(k) >= len(keyUsageNames) {
		k = 0
	}
	return keyUsageNames[k]
}

func ParseURIs(ss []string) []*url.URL {
	var uris []*url.URL
	for _, s := range ss {
		url, _ := url.Parse(s)
		if url != nil {
			uris = append(uris, url)
		}
	}
	return uris
}

func ParseIPAddresses(ss []string) []net.IP {
	var ips []net.IP
	for _, s := range ss {
		ip := net.ParseIP(s)
		if ip != nil {
			ips = append(ips, ip)
		}
	}
	return ips
}

func ParseKeyUsage(s string) x509.KeyUsage {
	sl := strings.ToLower(s)
	for i, kun := range keyUsageNames {
		if strings.ToLower(kun) == sl {
			return x509.KeyUsage(i)
		}
	}
	return 0
}

func ParseExtKeyUsage(ss []string) []x509.ExtKeyUsage {
	var ekus []x509.ExtKeyUsage
	// outer loop over each string in the given slice
	for _, s := range ss {
		sl := strings.ToLower(s)
		// inner loop through known usage names to match
		for i, kun := range ExtKeyUsageNames {
			if strings.ToLower(kun) != sl {
				continue
			}
			ekus = append(ekus, x509.ExtKeyUsage(i))
			break
		}
	}
	return ekus
}

func ParseTime(s string) (time.Time, error) {
	return time.Parse(TimeFormat, s)
}

func ParseSignatureAlgorithm(s string) x509.SignatureAlgorithm {
	switch s {
	case x509.MD5WithRSA.String():
		return x509.MD5WithRSA
	case x509.SHA1WithRSA.String():
		return x509.SHA1WithRSA
	case x509.SHA256WithRSA.String():
		return x509.SHA256WithRSA
	case x509.SHA384WithRSA.String():
		return x509.SHA384WithRSA
	case x509.SHA512WithRSA.String():
		return x509.SHA512WithRSA
	case x509.DSAWithSHA1.String():
		return x509.DSAWithSHA1
	case x509.DSAWithSHA256.String():
		return x509.DSAWithSHA256
	case x509.ECDSAWithSHA1.String():
		return x509.ECDSAWithSHA1
	case x509.ECDSAWithSHA256.String():
		return x509.ECDSAWithSHA256
	case x509.ECDSAWithSHA384.String():
		return x509.ECDSAWithSHA384
	case x509.ECDSAWithSHA512.String():
		return x509.ECDSAWithSHA512
	case x509.SHA256WithRSAPSS.String():
		return x509.SHA256WithRSAPSS
	case x509.SHA384WithRSAPSS.String():
		return x509.SHA384WithRSAPSS
	case x509.SHA512WithRSAPSS.String():
		return x509.SHA512WithRSAPSS
	case x509.PureEd25519.String():
		return x509.PureEd25519
	default:
		return x509.UnknownSignatureAlgorithm
	}
}

func ParseCurve(s string) elliptic.Curve {
	switch strings.ToLower(s) {
	case "p224":
		return elliptic.P224()
	case "p256":
		return elliptic.P256()
	case "p384":
		return elliptic.P384()
	case "p521":
		return elliptic.P521()
	default:
		return nil
	}
}

func ParsePublicKey(s string) (crypto.PublicKey, error) {
	by, err := ParseHexBytes(s)
	if err != nil {
		return "", err
	}
	return x509.ParsePKIXPublicKey(by)
}

func ParsePublicKeyAlgorithm(s string) x509.PublicKeyAlgorithm {
	switch s {
	case x509.RSA.String():
		return x509.RSA
	case x509.ECDSA.String():
		return x509.ECDSA
	case x509.Ed25519.String():
		return x509.Ed25519
	default:
		return x509.UnknownPublicKeyAlgorithm
	}
}
