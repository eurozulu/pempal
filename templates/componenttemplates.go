package templates

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"crypto/x509"
	"crypto/x509/pkix"
	"net/url"
)

type SubjectTemplate pkix.Name

func (s SubjectTemplate) String() string {
	return pkix.Name(s).String()
}

func (s SubjectTemplate) MarshalYAML() (interface{}, error) {
	return struct {
		CommonName         string                          `yaml:"CommonName"`
		SerialNumber       string                          `yaml:"SerialNumber,omitempty"`
		Organization       string                          `yaml:"Organization,omitempty"`
		OrganizationalUnit string                          `yaml:"OrganizationalUnit,omitempty"`
		Locality           string                          `yaml:"Locality,omitempty"`
		Province           string                          `yaml:"Province,omitempty"`
		Country            string                          `yaml:"Country,omitempty"`
		StreetAddress      string                          `yaml:"StreetAddress,omitempty"`
		PostalCode         string                          `yaml:"PostalCode,omitempty"`
		Names              []AttributeTypeAndValueTemplate `yaml:"Names,omitempty"`
		ExtraNames         []AttributeTypeAndValueTemplate `yaml:"ExtraNames,omitempty"`
	}{
		SerialNumber:       s.SerialNumber,
		CommonName:         s.CommonName,
		Organization:       strings.Join(s.Organization, ", "),
		OrganizationalUnit: strings.Join(s.OrganizationalUnit, ", "),
		Locality:           strings.Join(s.Locality, ", "),
		Province:           strings.Join(s.Province, ", "),
		Country:            strings.Join(s.Country, ", "),
		StreetAddress:      strings.Join(s.StreetAddress, ", "),
		PostalCode:         strings.Join(s.PostalCode, ", "),
		Names:              newAttributeTypeAndValueTemplateSlice(s.Names),
		ExtraNames:         newAttributeTypeAndValueTemplateSlice(s.ExtraNames),
	}, nil
}

type AttributeTypeAndValueTemplate struct {
	atv pkix.AttributeTypeAndValue
}

func newAttributeTypeAndValueTemplateSlice(avs []pkix.AttributeTypeAndValue) []AttributeTypeAndValueTemplate {
	var yav []AttributeTypeAndValueTemplate
	for _, av := range avs {
		yav = append(yav, AttributeTypeAndValueTemplate{atv: av})
	}
	return yav
}

func (crt AttributeTypeAndValueTemplate) MarshalYAML() (interface{}, error) {
	return fmt.Sprintf("%s : %v", crt.atv.Type, crt.atv.Value), nil
}

type ExtensionsTemplate pkix.Extension

func (et ExtensionsTemplate) MarshalYAML() (interface{}, error) {
	return &struct {
		Id       string `yaml:"Id,omitempty"`
		Critical bool   `yaml:"Critical,omitempty"`
		Value    string `yaml:"Value,omitempty"`
	}{
		Id:       fmt.Sprintf("%v", et.Id),
		Critical: et.Critical,
		Value:    fmt.Sprintf("%X", et.Value),
	}, nil
}

func ExtensionsTemplateSlice(exs []pkix.Extension) []ExtensionsTemplate {
	exts := make([]ExtensionsTemplate, len(exs))
	for i, ex := range exs {
		exts[i] = ExtensionsTemplate(ex)
	}
	return exts
}

func ExtensionsTemplateReslice(exs []ExtensionsTemplate) []pkix.Extension {
	exts := make([]pkix.Extension, len(exs))
	for i, ex := range exs {
		exts[i] = pkix.Extension(ex)
	}
	return exts
}

func ExtensionsSlice(exs []ExtensionsTemplate) []pkix.Extension {
	exts := make([]pkix.Extension, len(exs))
	for i, ex := range exs {
		exts[i] = pkix.Extension(ex)
	}
	return exts
}

type PublicKeyAlgorithmTemplate x509.PublicKeyAlgorithm

func NewPublicKeyAlgorithmTemplate(pka string) PublicKeyAlgorithmTemplate {
	for i, a := range PublicKeyAlgorithms {
		if strings.EqualFold(pka, a) {
			return PublicKeyAlgorithmTemplate(i)
		}
	}
	return PublicKeyAlgorithmTemplate(0)
}

var PublicKeyAlgorithms = []string{
	"",
	"RSA",
	"DSA",
	"ECDSA",
	"Ed25519",
}

func (pka PublicKeyAlgorithmTemplate) String() string {
	if pka < 0 || int(pka) >= len(PublicKeyAlgorithms) {
		return ""
	}
	return PublicKeyAlgorithms[int(pka)]
}
func (pka PublicKeyAlgorithmTemplate) MarshalYAML() (interface{}, error) {
	return pka.String(), nil
}

type SignatureAlgorithmTemplate x509.SignatureAlgorithm

var SignatureAlgorithms = []string{
	"",
	"MD2WithRSA",
	"MD5WithRSA",
	"SHA1WithRSA",
	"SHA256WithRSA",
	"SHA384WithRSA",
	"SHA512WithRSA",
	"DSAWithSHA1",
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

func (sa SignatureAlgorithmTemplate) String() string {
	if sa < 0 || int(sa) >= len(SignatureAlgorithms) {
		return ""
	}
	return SignatureAlgorithms[int(sa)]
}
func (sa SignatureAlgorithmTemplate) MarshalYAML() (interface{}, error) {
	return sa.String(), nil
}

var KeyUsages = []string{
	"",
	"DigitalSignature",
	"ContentCommitment",
	"KeyEncipherment",
	"DataEncipherment",
	"KeyAgreement",
	"CertSign",
	"CRLSign",
	"EncipherOnly",
	"DecipherOnly",
}

type KeyUsageTemplate x509.KeyUsage

func (ku KeyUsageTemplate) MarshalYAML() (interface{}, error) {
	if ku < 1 {
		return nil, nil
	}
	buf := bytes.NewBuffer(nil)
	for i := 1; i < len(KeyUsages); i++ {
		if (int(ku) & i) > 0 {
			if buf.Len() > 0 {
				buf.WriteString(" | ")
			}
			buf.WriteString(KeyUsages[i])
		}
	}
	return buf.String(), nil
}

var ExtKeyUsages = []string{
	"",
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

type ExtKeyUsagesTemplate x509.ExtKeyUsage

func (eku ExtKeyUsagesTemplate) MarshalYAML() (interface{}, error) {
	if eku < 0 || int(eku) >= len(ExtKeyUsages) {
		return "", nil
	}
	return ExtKeyUsages[int(eku)], nil
}
func ExtKeyUsagesTemplateSlice(eku []x509.ExtKeyUsage) []ExtKeyUsagesTemplate {
	if len(eku) == 0 {
		return nil
	}
	ekus := make([]ExtKeyUsagesTemplate, len(eku))
	for i, e := range eku {
		ekus[i] = ExtKeyUsagesTemplate(e)
	}
	return ekus
}

func ExtKeyUsagesTemplateReslice(eku []ExtKeyUsagesTemplate) []x509.ExtKeyUsage {
	if len(eku) == 0 {
		return nil
	}
	ekus := make([]x509.ExtKeyUsage, len(eku))
	for i, e := range eku {
		ekus[i] = x509.ExtKeyUsage(e)
	}
	return ekus
}

type IPAddressTemplate []net.IP

func (ips IPAddressTemplate) MarshalYAML() (interface{}, error) {
	if len(ips) == 0 {
		return nil, nil
	}
	ss := make([]string, len(ips))
	for i, s := range ips {
		ss[i] = s.String()
	}
	return strings.Join(ss, ", "), nil
}

type URIsTemplate []*url.URL

func (urit URIsTemplate) MarshalYAML() (interface{}, error) {
	if len(urit) == 0 {
		return nil, nil
	}
	ss := make([]string, len(urit))
	for i, s := range urit {
		ss[i] = s.String()
	}
	return strings.Join(ss, ", "), nil
}

type PEMCipherTemplate x509.PEMCipher

func NewPEMCipher(s string) PEMCipherTemplate {
	for i, pc := range PEMCiphers {
		if strings.EqualFold(s, pc) {
			return PEMCipherTemplate(i)
		}
	}
	return 0
}

func (pc PEMCipherTemplate) String() string {
	if pc < 0 || int(pc) >= len(PEMCiphers) {
		return ""
	}
	return PEMCiphers[pc]
}

var PEMCiphers = []string{
	"",
	"PEMCipherDES",
	"PEMCipher3DES",
	"PEMCipherAES128",
	"PEMCipherAES192",
	"PEMCipherAES256",
}
