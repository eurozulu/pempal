package templates

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"pempal/keytracker"
	"pempal/templates/parsers"
)

// MakeCertificate copies the given template values into a new Certificate.
// The resulting certificate is not signed by any issuer.
func MakeCertificate(t Template) *x509.Certificate {
	c := &x509.Certificate{
		Signature:          stringToBytes(t.Value(parsers.X509Signature)),
		SignatureAlgorithm: stringToSignatureAlgorithm(t.Value(parsers.X509SignatureAlgorithm)),
		Version:            stringToInt(t.Value(parsers.X509Version)),
		SerialNumber:       stringToBigInt(t.Value(parsers.X509SerialNumber)),

		SubjectKeyId:   stringToBytes(t.Value(parsers.X509SubjectKeyId)),
		AuthorityKeyId: stringToBytes(t.Value(parsers.X509AuthorityKeyId)),
		Issuer:         readNameProperties(t.ValueMap(parsers.X509Issuer)),
		Subject:        readNameProperties(t.ValueMap(parsers.X509Subject)),

		NotBefore:             stringToTime(t.Value(parsers.X509NotBefore)),
		NotAfter:              stringToTime(t.Value(parsers.X509NotAfter)),
		KeyUsage:              stringToKeyUsage(t.Value(parsers.X509KeyUsage)),
		Extensions:            stringToExtensions(t.Value(parsers.X509Extensions)),
		ExtraExtensions:       stringToExtensions(t.Value(parsers.X509ExtraExtensions)),
		ExtKeyUsage:           stringToExtKeyUsage(t.Value(parsers.X509ExtKeyUsage)),
		IsCA:                  stringToBool(t.Value(parsers.X509IsCA)),
		MaxPathLen:            stringToInt(t.Value(parsers.X509MaxPathLen)),
		MaxPathLenZero:        stringToBool(t.Value(parsers.X509MaxPathLenZero)),
		IssuingCertificateURL: stringToStringArray(t.Value(parsers.X509IssuingCertificateURL)),
		DNSNames:              stringToStringArray(t.Value(parsers.X509DNSNames)),
		EmailAddresses:        stringToStringArray(t.Value(parsers.X509EmailAddresses)),
		IPAddresses:           stringToIPs(t.Value(parsers.X509IPAddresses)),
		URIs:                  stringToURLs(t.Value(parsers.X509URIs)),
		CRLDistributionPoints: stringToStringArray(t.Value(parsers.X509CRLDistributionPoints)),
	}

	puk, pka := t.PublicKey()
	if puk != nil && !reflect.ValueOf(puk).IsNil() {
		c.PublicKey = puk
		c.PublicKeyAlgorithm = pka
	}
	return c
}

func IssueCertificate(issuer keytracker.Identity, keypass string, t Template) ([]byte, error) {
	c := MakeCertificate(t)
	puk := c.PublicKey
	if puk == nil || reflect.ValueOf(puk).IsNil() {
		return nil, fmt.Errorf("no public key in template")
	}
	prk, err := issuer.Key().PrivateKeyDecrypted(keypass)
	if err != nil {
		return nil, err
	}

	// Ensure identity has a certificate, otherwise use template cert = self signed.
	issCert := issuer.Certificate()
	if issCert == nil {
		issCert = c
	}
	return x509.CreateCertificate(rand.Reader, c, issuer.Certificate(), puk, prk)
}

func stringToBytes(s string) []byte {
	if s == "" {
		return nil
	}
	by, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Println(err)
		return nil
	}
	return by
}

func stringToTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", s)
	if err != nil {
		log.Println(err)
		return time.Time{}
	}
	return t
}

func stringToInt(s string) int {
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func stringToBigInt(s string) *big.Int {
	return big.NewInt(int64(stringToInt(s)))
}

func stringToBool(s string) bool {
	if s == "" {
		return false
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return b
}

func stringToStringArray(s string) []string {
	if s == "" {
		return nil
	}
	var ss []string
	for _, sz := range strings.Split(s, ",") {
		ss = append(ss, strings.TrimSpace(sz))
	}
	return ss

}

func stringToIPs(s string) []net.IP {
	if s == "" {
		return nil
	}
	ss := strings.Split(s, ",")
	ips := make([]net.IP, len(ss))
	for i, sz := range ss {
		ips[i] = net.ParseIP(sz)
	}
	return ips
}

func stringToURLs(s string) []*url.URL {
	if s == "" {
		return nil
	}
	ss := strings.Split(s, ",")
	urls := make([]*url.URL, len(ss))
	for i, sz := range ss {
		u, err := url.Parse(sz)
		if err != nil {
			log.Println(err)
			continue
		}
		urls[i] = u
	}
	return urls
}

func readNameProperties(t Template) pkix.Name {
	var n pkix.Name
	n.SerialNumber = t.Value(parsers.DNSerialNumber)
	n.CommonName = t.Value(parsers.DNCommonName)
	n.Organization = strings.Split(t.Value(parsers.DNOrganization), ",")
	n.OrganizationalUnit = strings.Split(t.Value(parsers.DNOrganizationalUnit), ",")
	n.StreetAddress = strings.Split(t.Value(parsers.DNStreetAddress), ",")
	n.Locality = strings.Split(t.Value(parsers.DNLocality), ",")
	n.Province = strings.Split(t.Value(parsers.DNProvince), ",")
	n.Country = strings.Split(t.Value(parsers.DNCountry), ",")
	n.PostalCode = strings.Split(t.Value(parsers.DNPostalCode), ",")
	return n
}

func stringToKeyUsage(s string) x509.KeyUsage {
	if s == "" {
		return 0
	}
	var ku x509.KeyUsage
	ss := strings.Split(s, ",")
	for _, sz := range ss {
		sk := parseKeyUsage(sz)
		if sk == 0 {
			continue
		}
		ku = ku | (1 << sk)
	}
	return ku
}

func stringToExtKeyUsage(s string) []x509.ExtKeyUsage {
	if s == "" {
		return nil
	}
	var ekus []x509.ExtKeyUsage
	for _, ss := range strings.Split(s, ", ") {
		eku := parseExtKeyUsage(ss)
		if eku == 0 {
			continue
		}
		ekus = append(ekus, eku)
	}
	return ekus
}
func stringToExtensions(s string) []pkix.Extension {
	if s == "" {
		return nil
	}
	var exts []pkix.Extension
	for _, ss := range strings.Split(s, ",") {
		// expect extensions as "id:!:value"  where != is critical, :: when not critical
		se := strings.Split(ss, ":")
		if len(se) < 3 {
			continue
		}
		exts = append(exts, pkix.Extension{
			Id:       parseObjectIdentifier(se[0]),
			Critical: se[1] == "!",
			Value:    []byte(se[2]),
		})
	}
	return exts
}

func parseObjectIdentifier(s string) asn1.ObjectIdentifier {
	iv := strings.Split(s, ".")
	oi := make(asn1.ObjectIdentifier, len(iv))
	for i, v := range iv {
		pv, err := strconv.Atoi(v)
		if err != nil {
			log.Println(err)
			return nil
		}
		oi[i] = pv
	}
	return oi
}

func parseKeyUsage(s string) x509.KeyUsage {
	for i, ku := range parsers.KeyUsageNames {
		if strings.EqualFold(ku, s) {
			return x509.KeyUsage(i)
		}
	}
	return 0
}
func parseExtKeyUsage(s string) x509.ExtKeyUsage {
	for i, ku := range parsers.ExtKeyUsageNames {
		if strings.EqualFold(ku, s) {
			return x509.ExtKeyUsage(i)
		}
	}
	return 0
}

func stringToSignatureAlgorithm(s string) x509.SignatureAlgorithm {
	if s == "" {
		return x509.UnknownSignatureAlgorithm
	}
	for i, sa := range signatureAlgorithmNames {
		if strings.EqualFold(sa, s) {
			return x509.SignatureAlgorithm(i)
		}
	}
	return x509.UnknownSignatureAlgorithm
}

var signatureAlgorithmNames = []string{
	"UnknownSignatureAlgorithm",
	"MD2WithRSA",
	"MD5WithRSA",
	"SHA1WithRSA",
	"SHA256WithRSA",
	"SHA384WithRSA",
	"SHA512WithRSA",
	"DSAWithSHA1 ",
	"DSAWithSHA256",
	"ECDSAWithSHA1",
	"ECDSAWithSHA256",
	"ECDSAWithSHA384",
	"ECDSAWithSHA512",
	"SHA256WithRSAPSS",
	"SHA384WithRSAPSS",
	"SHA512WithRSAPSS",
	"PureEd25519",
}
