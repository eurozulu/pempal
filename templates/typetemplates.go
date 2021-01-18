package templates

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"gopkg.in/yaml.v3"
	"math/big"
	"strconv"
	"strings"
	"time"
)

// The type templates template the standard enumerated properties
// so they can be interchanged with their string equivelents.

type KeyUsage x509.KeyUsage

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

func (k KeyUsage) String() string {
	if k < 1 {
		return ""
	}
	buf := bytes.NewBuffer(nil)
	for i := 1; i < len(KeyUsages); i++ {
		if (k & KeyUsage(i)) > 0 {
			if buf.Len() > 0 {
				buf.WriteString(" | ")
			}
			buf.WriteString(KeyUsages[i])
		}
	}
	return buf.String()
}

func (k KeyUsage) MarshalYAML() (interface{}, error) {
	return k.String(), nil
}

func (k *KeyUsage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	if err := unmarshal(&v); err != nil {
		return nil
	}
	vs := strings.Split(v, "|")
	var ku KeyUsage
	for _, kus := range vs {
		if kus == "" {
			continue
		}
		i := indexOf(kus, KeyUsages)
		if i < 1 {
			return fmt.Errorf("%s is not a known key usage", kus)
		}
		ku += KeyUsage(i)
	}
	*k = ku
	return nil
}

type ExtKeyUsage x509.ExtKeyUsage

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

func ExtKeyUsagesReslice(rku []ExtKeyUsage) []x509.ExtKeyUsage {
	kus := make([]x509.ExtKeyUsage, len(rku))
	for i, ku := range rku {
		kus[i] = x509.ExtKeyUsage(ku)
	}
	return kus
}
func ExtKeyUsagesSlice(rku []x509.ExtKeyUsage) []ExtKeyUsage {
	kus := make([]ExtKeyUsage, len(rku))
	for i, ku := range rku {
		kus[i] = ExtKeyUsage(ku)
	}
	return kus
}

func (k ExtKeyUsage) String() string {
	if k < 1 || int(k) > len(ExtKeyUsages) {
		return ""
	}
	return ExtKeyUsages[k]
}

func (k ExtKeyUsage) MarshalYAML() (interface{}, error) {
	return k.String(), nil
}

func (k *ExtKeyUsage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	if err := unmarshal(&v); err != nil {
		return nil
	}
	if v == "" {
		return nil
	}
	i := indexOf(v, ExtKeyUsages)
	if i < 1 {
		return fmt.Errorf("%s is not a known extended key usage", v)
	}
	*k = ExtKeyUsage(i)
	return nil
}

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

type SignatureAlgorithm x509.SignatureAlgorithm

func (s SignatureAlgorithm) String() string {
	if s < 1 || int(s) > len(SignatureAlgorithms) {
		return ""
	}
	return SignatureAlgorithms[s]
}

func (s SignatureAlgorithm) MarshalYAML() (interface{}, error) {
	return s.String(), nil
}

func (s *SignatureAlgorithm) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	if err := unmarshal(&v); err != nil {
		return err
	}
	if v == "" {
		return nil
	}
	i := indexOf(v, SignatureAlgorithms)
	if i < 1 {
		return fmt.Errorf("%s is not a known SignatureAlgorithm", v)
	}
	*s = SignatureAlgorithm(i)
	return nil
}

func indexOf(s string, ss []string) int {
	for i, sz := range ss {
		if strings.EqualFold(sz, s) {
			return i
		}
	}
	return -1
}

var PublicKeyAlgorithms = []string{
	"",
	"RSA",
	"DSA",
	"ECDSA",
	"Ed25519",
}

type PublicKeyAlgorithm x509.PublicKeyAlgorithm

func (p PublicKeyAlgorithm) String() string {
	if p < 1 || int(p) >= len(PublicKeyAlgorithms) {
		return ""
	}
	return PublicKeyAlgorithms[p]
}

func (p PublicKeyAlgorithm) MarshalYAML() (interface{}, error) {
	return p.String(), nil
}

func (p *PublicKeyAlgorithm) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	if err := unmarshal(&v); err != nil {
		return err
	}
	if v == "" {
		return nil
	}
	pka, err := ParsePublicKeyAlgorithm(v)
	if err != nil {
		return err
	}
	*p = pka
	return nil
}

func ParsePublicKeyAlgorithm(s string) (PublicKeyAlgorithm, error) {
	i := indexOf(s, PublicKeyAlgorithms)
	if i < 1 {
		return 0, fmt.Errorf("%s is not a known PublicKey Algorithm", s)
	}
	return PublicKeyAlgorithm(i), nil
}

var PEMCiphers = []string{
	"",
	"DES",
	"3DES",
	"AES128",
	"AES192",
	"AES256",
}

type PEMCipher x509.PEMCipher

func (pc PEMCipher) String() string {
	if pc < 1 || int(pc) >= len(PEMCiphers) {
		return ""
	}
	return PEMCiphers[pc]
}

func (pc PEMCipher) MarshalYAML() (interface{}, error) {
	return pc.String(), nil
}

func (pc *PEMCipher) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	if err := unmarshal(&v); err != nil {
		return err
	}
	if v == "" {
		return nil
	}
	pka, err := ParsePEMCipher(v)
	if err != nil {
		return err
	}
	*pc = pka
	return nil
}

func ParsePEMCipher(s string) (PEMCipher, error) {
	i := indexOf(s, PEMCiphers)
	if i < 1 {
		return 0, fmt.Errorf("%s is not a known PEMCipher", s)
	}
	return PEMCipher(i), nil
}

type RevokedCertificate struct {
	SerialNumber   *big.Int    `yaml:"SerialNumber"`
	RevocationTime time.Time   `yaml:"RevocationTime"`
	Extensions     []Extension `yaml:"optional"`
}

func NewRevokedCertificate(serialNumber *big.Int, revocationTime time.Time, extensions []pkix.Extension) *RevokedCertificate {
	return &RevokedCertificate{SerialNumber: serialNumber, RevocationTime: revocationTime, Extensions: ExtensionSlice(extensions)}
}
func RevokedCertificatesSlice(rcl []pkix.RevokedCertificate) []RevokedCertificate {
	rcs := make([]RevokedCertificate, len(rcl))
	for i, ex := range rcl {
		rcs[i] = RevokedCertificate{
			SerialNumber:   ex.SerialNumber,
			RevocationTime: ex.RevocationTime,
			Extensions:     ExtensionSlice(ex.Extensions),
		}
	}
	return rcs
}
func RevokedCertificatesReslice(rcl []RevokedCertificate) []pkix.RevokedCertificate {
	rcs := make([]pkix.RevokedCertificate, len(rcl))
	for i, ex := range rcl {
		rcs[i] = pkix.RevokedCertificate{
			SerialNumber:   ex.SerialNumber,
			RevocationTime: ex.RevocationTime,
			Extensions:     ExtensionReslice(ex.Extensions),
		}
	}
	return rcs
}

type Extension struct {
	Id       ObjectIdentifier `yaml:"id"`
	Critical bool             `yaml:"optional"`
	Value    []byte           `yaml:"value"`
}

func (o *Extension) UnmarshalYAML(value *yaml.Node) error {
	s := &struct {
		Id       ObjectIdentifier `yaml:"id"`
		Critical bool             `yaml:"optional"`
		Value    string           `yaml:"value"`
	}{}
	if err := yaml.Unmarshal([]byte(value.Value), s); err != nil {
		return err
	}
	o.Id = s.Id
	o.Critical = s.Critical

	vs := strings.Split(s.Value, " ")
	o.Value = make([]byte, len(vs))
	for i, bs := range vs {
		v, err := strconv.Atoi(bs)
		if err != nil {
			return err
		}
		o.Value[i] = byte(v)
	}
	return nil
}

func (o Extension) MarshalYAML() (interface{}, error) {
	s := &struct {
		Id       ObjectIdentifier `yaml:"Id"`
		Critical bool             `yaml:"Critical"`
		Value    string           `yaml:"Value"`
	}{
		Id:       o.Id,
		Critical: o.Critical,
		Value:    fmt.Sprintf("%X", o.Value),
	}
	return s, nil
}

func ExtensionSlice(e []pkix.Extension) []Extension {
	exts := make([]Extension, len(e))
	for i, ex := range e {
		exts[i] = Extension{
			Id:       ObjectIdentifier(ex.Id),
			Critical: ex.Critical,
			Value:    ex.Value,
		}
	}
	return exts
}
func ExtensionReslice(e []Extension) []pkix.Extension {
	exts := make([]pkix.Extension, len(e))
	for i, ex := range e {
		exts[i] = pkix.Extension{
			Id:       asn1.ObjectIdentifier(ex.Id),
			Critical: ex.Critical,
			Value:    ex.Value,
		}
	}
	return exts
}

type ObjectIdentifier asn1.ObjectIdentifier

func (o *ObjectIdentifier) UnmarshalYAML(value *yaml.Node) error {
	sbs := strings.Split(value.Value, " ")
	for _, sb := range sbs {
		if sb == "" {
			continue
		}
		i, err := strconv.ParseInt(sb, 10, 32)
		if err != nil {
			return err
		}
		*o = append(*o, int(i))
	}
	return nil
}

func (o ObjectIdentifier) MarshalYAML() (interface{}, error) {
	return fmt.Sprintf("%x", o), nil
}

type AttributeTypeAndValue struct {
	Type  ObjectIdentifier `yaml:"ObjectIdentifier,omitempty"`
	Value string           `yaml:"Value,omitempty"`
}

func AttributeTypeAndValueSlice(e []pkix.AttributeTypeAndValue) []AttributeTypeAndValue {
	attrs := make([]AttributeTypeAndValue, len(e))
	for i, atv := range e {
		attrs[i] = AttributeTypeAndValue{
			Type:  ObjectIdentifier(atv.Type),
			Value: fmt.Sprintf("%v", atv.Value),
		}
	}
	return attrs
}

func AttributeTypeAndValueReslice(e []AttributeTypeAndValue) []pkix.AttributeTypeAndValue {
	atvs := make([]pkix.AttributeTypeAndValue, len(e))
	for i, atv := range e {
		atvs[i] = pkix.AttributeTypeAndValue{
			Type:  asn1.ObjectIdentifier(atv.Type),
			Value: atv.Value,
		}
	}
	return atvs
}
