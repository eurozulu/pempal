package pemproperties

import (
	"crypto/x509"
	"net"
	"net/url"
	"strings"
	"time"
)

const TimeFormat = time.RFC850

type Property[T any] interface {
	String(T) string
	Parse(s string) T
}

type SelectProperty[T any] interface {
	Property[T]
	Values() []string
}

var publicKeyAlgoName = [...]string{
	x509.RSA:     "RSA",
	x509.DSA:     "DSA",
	x509.ECDSA:   "ECDSA",
	x509.Ed25519: "ED25519",
}

var keyUsageNames = [...]string{
	x509.KeyUsageDigitalSignature:  "KeyUsageDigitalSignature",
	x509.KeyUsageContentCommitment: "KeyUsageContentCommitment",
	x509.KeyUsageKeyEncipherment:   "KeyUsageKeyEncipherment",
	x509.KeyUsageDataEncipherment:  "KeyUsageDataEncipherment",
	x509.KeyUsageKeyAgreement:      "KeyUsageKeyAgreement",
	x509.KeyUsageCertSign:          "KeyUsageCertSign",
	x509.KeyUsageCRLSign:           "KeyUsageCRLSign",
	x509.KeyUsageEncipherOnly:      "KeyUsageEncipherOnly",
	x509.KeyUsageDecipherOnly:      "KeyUsageDecipherOnly",
}

var ExtKeyUsageNames = [...]string{
	x509.ExtKeyUsageAny:                            "ExtKeyUsageAny",
	x509.ExtKeyUsageServerAuth:                     "ExtKeyUsageServerAuth",
	x509.ExtKeyUsageClientAuth:                     "ExtKeyUsageClientAuth",
	x509.ExtKeyUsageCodeSigning:                    "ExtKeyUsageCodeSigning",
	x509.ExtKeyUsageEmailProtection:                "ExtKeyUsageEmailProtection",
	x509.ExtKeyUsageIPSECEndSystem:                 "ExtKeyUsageIPSECEndSystem",
	x509.ExtKeyUsageIPSECTunnel:                    "ExtKeyUsageIPSECTunnel",
	x509.ExtKeyUsageIPSECUser:                      "ExtKeyUsageIPSECUser",
	x509.ExtKeyUsageTimeStamping:                   "ExtKeyUsageTimeStamping",
	x509.ExtKeyUsageOCSPSigning:                    "ExtKeyUsageOCSPSigning",
	x509.ExtKeyUsageMicrosoftServerGatedCrypto:     "ExtKeyUsageMicrosoftServerGatedCrypto",
	x509.ExtKeyUsageNetscapeServerGatedCrypto:      "ExtKeyUsageNetscapeServerGatedCrypto",
	x509.ExtKeyUsageMicrosoftCommercialCodeSigning: "ExtKeyUsageMicrosoftCommercialCodeSigning",
	x509.ExtKeyUsageMicrosoftKernelCodeSigning:     "ExtKeyUsageMicrosoftKernelCodeSigning",
}

var signatureAlgoNames = [...]string{
	x509.MD5WithRSA:       "MD5WithRSA",
	x509.SHA1WithRSA:      "SHA1WithRSA",
	x509.SHA256WithRSA:    "SHA256WithRSA",
	x509.SHA384WithRSA:    "SHA384WithRSA",
	x509.SHA512WithRSA:    "SHA512WithRSA",
	x509.DSAWithSHA1:      "DSAWithSHA1",
	x509.DSAWithSHA256:    "DSAWithSHA256",
	x509.ECDSAWithSHA1:    "ECDSAWithSHA1",
	x509.ECDSAWithSHA256:  "ECDSAWithSHA256",
	x509.ECDSAWithSHA384:  "ECDSAWithSHA384",
	x509.ECDSAWithSHA512:  "ECDSAWithSHA512",
	x509.SHA256WithRSAPSS: "SHA256WithRSAPSS",
	x509.SHA384WithRSAPSS: "SHA384WithRSAPSS",
	x509.SHA512WithRSAPSS: "SHA512WithRSAPSS",
	x509.PureEd25519:      "PureEd25519",
}

type PublicKeyAlgorithmProperty struct{}

var _ SelectProperty[x509.PublicKeyAlgorithm] = &PublicKeyAlgorithmProperty{}

func (p PublicKeyAlgorithmProperty) String(t x509.PublicKeyAlgorithm) string {
	return t.String()
}

func (p PublicKeyAlgorithmProperty) Parse(s string) x509.PublicKeyAlgorithm {
	su := strings.ToUpper(s)
	for i, ss := range publicKeyAlgoName {
		if ss == su {
			return x509.PublicKeyAlgorithm(i)
		}
	}
	return x509.UnknownPublicKeyAlgorithm
}

func (p PublicKeyAlgorithmProperty) Values() []string {
	return publicKeyAlgoName[:]
}

type SignatureAlgorithmProperty struct {
}

func (sa SignatureAlgorithmProperty) String(t x509.SignatureAlgorithm) string {
	return t.String()
}

func (sa SignatureAlgorithmProperty) Parse(s string) x509.SignatureAlgorithm {
	su := strings.ToUpper(s)
	for i, ss := range signatureAlgoNames {
		if strings.ToUpper(ss) == su {
			return x509.SignatureAlgorithm(i)
		}
	}
	return x509.UnknownSignatureAlgorithm
}

func (sa SignatureAlgorithmProperty) Values() []string {
	return signatureAlgoNames[:]
}

type IPAddressListProperty struct{}

func (I IPAddressListProperty) String(t []net.IP) string {
	ss := make([]string, len(t))
	for i, ip := range t {
		ss[i] = ip.String()
	}
	return strings.Join(ss, ",")
}

func (I IPAddressListProperty) Parse(s string) []net.IP {
	ss := strings.Split(s, ",")
	ips := make([]net.IP, len(ss))
	for i, sz := range ss {
		ips[i] = net.ParseIP(strings.TrimSpace(sz))
	}
	return ips
}

type IPNetProperty struct{}

func (I IPNetProperty) String(t net.IPNet) string {
	return t.String()
}
func (I IPNetProperty) Strings(t []net.IPNet) []string {
	ipnets := make([]string, len(t))
	for i, ip := range t {
		ipnets[i] = I.String(ip)
	}
	return ipnets
}

func (I IPNetProperty) Parse(s string) *net.IPNet {
	return &net.IPNet{}
}
func (I IPNetProperty) ParseList(ss []string) []*net.IPNet {
	ipnets := make([]*net.IPNet, len(ss))
	for i, s := range ss {
		ipnets[i] = I.Parse(s)
	}
	return ipnets
}

type URIListProperty struct{}

func (I URIListProperty) String(t []*url.URL) string {
	ss := make([]string, len(t))
	for i, u := range t {
		ss[i] = u.String()
	}
	return strings.Join(ss, ",")
}

func (I URIListProperty) Parse(s string) []*url.URL {
	ss := strings.Split(s, ",")
	urls := make([]*url.URL, len(ss))
	for i, sz := range ss {
		urls[i], _ = url.Parse(sz)
	}
	return urls
}

type TimeProperty struct{}

func (tp TimeProperty) String(t time.Time) string {
	return t.Format(TimeFormat)
}

func (tp TimeProperty) Parse(s string) time.Time {
	t, err := time.Parse(TimeFormat, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

type KeyUsageProperty struct{}

func (k KeyUsageProperty) String(t x509.KeyUsage) string {
	return keyUsageNames[t]
}

func (k KeyUsageProperty) Parse(s string) x509.KeyUsage {
	s = strings.ToUpper(s)
	for i, n := range keyUsageNames {
		if strings.ToUpper(n) == s {
			return x509.KeyUsage(i)
		}
	}
	return 0
}
func (k KeyUsageProperty) Values() []string {
	return keyUsageNames[:]
}

type ExtKeyUsageProperty struct{}

func (k ExtKeyUsageProperty) String(t x509.ExtKeyUsage) string {
	return ExtKeyUsageNames[t]
}

func (k ExtKeyUsageProperty) Strings(t []x509.ExtKeyUsage) []string {
	ss := make([]string, len(t))
	for i, ku := range t {
		ss[i] = k.String(ku)
	}
	return ss
}

func (k ExtKeyUsageProperty) ParseList(ss []string) []x509.ExtKeyUsage {
	kus := make([]x509.ExtKeyUsage, len(ss))
	for i, s := range ss {
		kus[i] = k.Parse(s)
	}
	return kus
}

func (k ExtKeyUsageProperty) Parse(s string) x509.ExtKeyUsage {
	s = strings.ToUpper(s)
	for i, n := range ExtKeyUsageNames {
		if strings.ToUpper(n) == s {
			return x509.ExtKeyUsage(i)
		}
	}
	return 0
}
func (k ExtKeyUsageProperty) Values() []string {
	return ExtKeyUsageNames[:]
}
