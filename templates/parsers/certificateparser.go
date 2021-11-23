package parsers

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
)

const (
	X509Signature             = "Signature"
	X509SignatureAlgorithm    = "SignatureAlgorithm"
	X509PublicKey             = "PublicKey"
	X509PublicKeyAlgorithm    = "PublicKeyAlgorithm"
	X509Version               = "Version"
	X509SerialNumber          = "SerialNumber"
	X509SubjectKeyId          = "SubjectKeyId"
	X509AuthorityKeyId        = "AuthorityKeyId"
	X509Issuer                = "Issuer"
	X509IssuerDN              = "IssuerDN"
	X509Subject               = "Subject"
	X509SubjectDN             = "SubjectDN"
	X509NotBefore             = "NotBefore"
	X509NotAfter              = "NotAfter"
	X509KeyUsage              = "Certificates"
	X509ExtKeyUsage           = "ExtKeyUsage"
	X509Extensions            = "Extensions"
	X509ExtraExtensions       = "ExtExtensions"
	X509IsCA                  = "IsCA"
	X509MaxPathLen            = "MaxPathLen"
	X509MaxPathLenZero        = "MaxPathLenZero"
	X509DNSNames              = "DNSNames"
	X509EmailAddresses        = "EmailAddresses"
	X509IPAddresses           = "IPAddresses"
	X509URIs                  = "URIs"
	X509IssuingCertificateURL = "IssuingCertificateURL"
	X509CRLDistributionPoints = "CRLDistributionPoints"
)
const (
	DNSerialNumber       = "SerialNumber"
	DNCommonName         = "CommonName"
	DNOrganizationalUnit = "OrganizationalUnit"
	DNOrganization       = "Organization"
	DNStreetAddress      = "StreetAddress"
	DNLocality           = "Locality"
	DNProvince           = "Province"
	DNCountry            = "Country"
	DNPostalCode         = "PostalCode"
)

var AllDNNames = []string{
	DNSerialNumber,
	DNCommonName,
	DNOrganizationalUnit,
	DNOrganization,
	DNStreetAddress,
	DNLocality,
	DNProvince,
	DNCountry,
	DNPostalCode,
}

var AllCertificateNames = []string{
	X509Signature,
	X509SignatureAlgorithm,
	X509PublicKey,
	X509PublicKeyAlgorithm,
	X509Version,
	X509SerialNumber,
	X509SubjectKeyId,
	X509Issuer,
	X509Subject,
	X509NotBefore,
	X509NotAfter,
	X509KeyUsage,
	X509ExtKeyUsage,
	X509Extensions,
	X509ExtraExtensions,
	X509IsCA,
	X509MaxPathLen,
	X509MaxPathLenZero,
	X509DNSNames,
	X509EmailAddresses,
	X509IPAddresses,
	X509URIs,
	X509IssuingCertificateURL,
	X509CRLDistributionPoints,
}

var KeyUsageNames = []string{
	"Unknown",
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

var ExtKeyUsageNames = []string{
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

type CertificateParser struct{}

func (cp CertificateParser) KnownNames() []string {
	return AllCertificateNames
}

func (cp CertificateParser) Parse(b *pem.Block) (map[string]interface{}, error) {
	c, err := x509.ParseCertificate(b.Bytes)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{
		PEM_TYPE: b.Type,
	}
	m[X509Signature] = bytesToString(c.Signature)
	m[X509SignatureAlgorithm] = c.SignatureAlgorithm.String()

	m[X509PublicKey] = publicKeyToString(c.PublicKey)
	m[X509PublicKeyAlgorithm] = c.PublicKeyAlgorithm.String()

	m[X509Version] = strconv.Itoa(c.Version)
	m[X509SerialNumber] = c.SerialNumber.String()

	m[X509SubjectKeyId] = bytesToString(c.SubjectKeyId)
	m[X509AuthorityKeyId] = bytesToString(c.AuthorityKeyId)

	m[X509Issuer] = nameToStringMap(c.Issuer)
	m[X509IssuerDN] = c.Issuer.String()

	m[X509Subject] = nameToStringMap(c.Subject)
	m[X509SubjectDN] = c.Subject.String()

	m[X509NotBefore] = c.NotBefore.String()
	m[X509NotAfter] = c.NotAfter.String()

	m[X509Extensions] = extensionsToString(c.Extensions)
	m[X509ExtraExtensions] = extensionsToString(c.ExtraExtensions)

	m[X509KeyUsage] = KeyUsageToString(c.KeyUsage)
	m[X509ExtKeyUsage] = ExtKeyUsageToString(c.ExtKeyUsage)

	m[X509IsCA] = strconv.FormatBool(c.IsCA)
	m[X509MaxPathLen] = strconv.Itoa(c.MaxPathLen)
	m[X509MaxPathLenZero] = strconv.FormatBool(c.MaxPathLenZero)

	m[X509IssuingCertificateURL] = strings.Join(c.IssuingCertificateURL, ", ")

	m[X509DNSNames] = strings.Join(c.DNSNames, ", ")
	m[X509EmailAddresses] = strings.Join(c.EmailAddresses, ", ")
	m[X509IPAddresses] = iPsToString(c.IPAddresses)
	m[X509URIs] = uRLsToString(c.URIs)

	m[X509CRLDistributionPoints] = strings.Join(c.CRLDistributionPoints, ", ")

	return m, nil
}

func nameToStringMap(n pkix.Name) map[string]interface{} {
	return map[string]interface{}{
		DNSerialNumber:       n.SerialNumber,
		DNCommonName:         n.CommonName,
		DNOrganization:       n.Organization,
		DNOrganizationalUnit: n.OrganizationalUnit,
		DNStreetAddress:      n.StreetAddress,
		DNLocality:           n.Locality,
		DNProvince:           n.Province,
		DNCountry:            n.Country,
		DNPostalCode:         n.PostalCode,
	}
}

func addDNProperty(m map[string]string, s, title, name string) {
	if s == "" {
		return
	}
	m[strings.Join([]string{title, name}, ".")] = s
}

func addDNProperties(m map[string]string, s []string, title, name string) {
	if len(s) == 0 {
		return
	}
	m[strings.Join([]string{title, name}, ".")] = strings.Join(s, ", ")
}

func bytesToString(by []byte) string {
	if len(by) == 0 {
		return ""
	}
	return base64.StdEncoding.EncodeToString(by)
}

func iPsToString(ips []net.IP) string {
	var s []string
	for _, ip := range ips {
		s = append(s, ip.String())
	}
	return strings.Join(s, ", ")
}
func uRLsToString(urls []*url.URL) string {
	var s []string
	for _, u := range urls {
		s = append(s, u.String())
	}
	return strings.Join(s, ", ")
}

func KeyUsageToString(ku x509.KeyUsage) string {
	kus := bytes.NewBufferString("")
	for i, k := range KeyUsageNames {
		if i == 0 { // skip 'unknown'
			continue
		}
		// Logic AND bitmask from index with ku
		if ku&(1<<i) == 0 {
			continue
		}
		if kus.Len() > 0 {
			kus.WriteString(", ")
		}
		kus.WriteString(k)
	}
	return kus.String()
}

func ExtKeyUsageToString(kus []x509.ExtKeyUsage) string {
	ekus := bytes.NewBufferString("")
	for _, k := range kus {
		if k < 0 || int(k) > len(ExtKeyUsageNames) {
			continue
		}
		if ekus.Len() > 0 {
			ekus.WriteString(", ")
		}
		ekus.WriteString(ExtKeyUsageNames[k])
	}
	return ekus.String()
}

func extensionsToString(exts []pkix.Extension) string {
	if len(exts) == 0 {
		return ""
	}
	es := make([]string, len(exts))
	for i, ex := range exts {
		buf := bytes.NewBufferString(ex.Id.String())
		buf.WriteRune(':')
		if ex.Critical {
			buf.WriteRune('!')
		}
		buf.WriteRune(':')
		buf.Write(ex.Value)
		es[i] = buf.String()
	}
	return strings.Join(es, ", ")
}

func publicKeyToString(puk crypto.PublicKey) string {
	by, err := x509.MarshalPKIXPublicKey(puk)
	if err != nil {
		log.Println(err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(by)
}
