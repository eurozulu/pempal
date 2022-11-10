package pemtypes

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"math/big"
	"pempal/pemproperties"
	"pempal/templates"
	"strings"
)

type certificateType struct {
	cert x509.Certificate
}

func (ct certificateType) String() string {
	return fmt.Sprintf("%s\t%s", Certificate.String(), ct.cert.Subject.String())
}

// MarshalBinary marshals the certificate as ASN.1 DER bytes
func (ct certificateType) MarshalBinary() (data []byte, err error) {
	return ct.cert.Raw, nil
}

// UnmarshalBinary attempts to parse the given data as an ASN.1 DER encoded certificate
func (ct *certificateType) UnmarshalBinary(data []byte) error {
	c, err := x509.ParseCertificate(data)
	if err != nil {
		return err
	}
	ct.cert = *c
	return nil
}

// MarshalText marshals the certificate as a PEM encoded block
func (ct certificateType) MarshalText() (text []byte, err error) {
	der, err := ct.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  Certificate.String(),
		Bytes: der,
	}), nil
}

// UnmarshalText attemps to parse the given text as a PEM encoded certificate
// If the data contains ore than one PEM blocks, the first certificate found is parsed
func (ct *certificateType) UnmarshalText(text []byte) error {
	blocks := ReadPEMBlocks(text, Certificate)
	if len(blocks) == 0 {
		return fmt.Errorf("no certificate pem found")
	}
	return ct.UnmarshalBinary(blocks[0].Bytes)
}

// MarshalYAML marshals the certificate as a YAML document
func (ct certificateType) MarshalYAML() (interface{}, error) {
	t := templates.CertificateTemplate{}
	ct.applyToTemplate(&t)
	return yaml.Marshal(&t)
}

// UnmarshalYAML attempts to read the given value as a YAML encoded certificate
func (ct *certificateType) UnmarshalYAML(value *yaml.Node) error {
	t := templates.CertificateTemplate{}
	if err := value.Decode(&t); err != nil {
		return err
	}
	ct.applyTemplate(t)
	return nil
}

// applyToTemplate applies the current state of the certificate into the given template
func (ct certificateType) applyToTemplate(t *templates.CertificateTemplate) {
	t.Signature = string(ct.cert.Signature)
	t.SignatureAlgorithm = ct.cert.SignatureAlgorithm.String()

	t.SignatureAlgorithm = ct.cert.SignatureAlgorithm.String()
	t.PublicKeyAlgorithm = ct.cert.PublicKeyAlgorithm.String()
	if ct.cert.PublicKey != nil {
		t.PublicKey = &templates.PublicKeyTemplate{}
		put := publicKeyType{puk: ct.cert.PublicKey}
		put.applyToTemplate(t.PublicKey)
	}

	t.Version = ct.cert.Version
	t.SerialNumber = ct.cert.SerialNumber.Int64()

	t.Issuer = &templates.NameTemplate{}
	nt := dnameType{dname: ct.cert.Issuer}
	nt.applyToTemplate(t.Issuer)
	t.Subject = &templates.NameTemplate{}
	nt.dname = ct.cert.Subject
	nt.applyToTemplate(t.Subject)

	tt := pemproperties.TimeProperty{}
	t.NotBefore = tt.String(ct.cert.NotBefore)
	t.NotAfter = tt.String(ct.cert.NotAfter)

	t.KeyUsage = pemproperties.KeyUsageProperty{}.String(ct.cert.KeyUsage)
	t.ExtKeyUsage = pemproperties.ExtKeyUsageProperty{}.Strings(ct.cert.ExtKeyUsage)
	//UnknownExtKeyUsage []asn1.ObjectIdentifier // Encountered extended key usages unknown to this package.

	// TODO
	//Extensions []pkix.Extension
	//ExtraExtensions []pkix.Extension
	//UnhandledCriticalExtensions []asn1.ObjectIdentifier

	t.BasicConstraintsValid = ct.cert.BasicConstraintsValid
	t.IsCA = ct.cert.IsCA

	t.MaxPathLen = ct.cert.MaxPathLen
	t.MaxPathLenZero = ct.cert.MaxPathLenZero

	t.SubjectKeyId = string(ct.cert.SubjectKeyId)
	t.AuthorityKeyId = string(ct.cert.AuthorityKeyId)

	// RFC 5280, 4.2.2.1 (Authority Information Access)
	t.OCSPServer = ct.cert.OCSPServer
	t.IssuingCertificateURL = ct.cert.IssuingCertificateURL

	if len(ct.cert.DNSNames) > 0 {
		t.DNSNames = ct.cert.DNSNames
	}
	if len(ct.cert.EmailAddresses) > 0 {
		t.EmailAddresses = ct.cert.EmailAddresses
	}
	if len(ct.cert.IPAddresses) > 0 {
		t.IPAddresses = strings.Split(pemproperties.IPAddressListProperty{}.String(ct.cert.IPAddresses), ",")
	}
	if len(ct.cert.URIs) > 0 {
		t.URIs = strings.Split(pemproperties.URIListProperty{}.String(ct.cert.URIs), ",")
	}

	t.PermittedDNSDomainsCritical = ct.cert.PermittedDNSDomainsCritical
	t.PermittedDNSDomains = ct.cert.PermittedDNSDomains
	t.ExcludedDNSDomains = ct.cert.ExcludedDNSDomains
	//t.PermittedIPRanges = MarshalIPAddresses(ct.cert.PermittedIPRanges)
	//t.ExcludedIPRanges            []*net.IPNet
	t.PermittedEmailAddresses = ct.cert.PermittedEmailAddresses
	t.ExcludedEmailAddresses = ct.cert.ExcludedEmailAddresses
	t.PermittedURIDomains = ct.cert.PermittedURIDomains
	t.ExcludedURIDomains = ct.cert.ExcludedURIDomains
	t.CRLDistributionPoints = ct.cert.CRLDistributionPoints

}

// applyTemplate applies any non empty properties from the given template, into the certificate.
func (ct *certificateType) applyTemplate(t templates.CertificateTemplate) {
	if t.Signature != "" {
		ct.cert.Signature = []byte(t.Signature)
	}
	if t.SignatureAlgorithm != "" {
		ct.cert.SignatureAlgorithm = pemproperties.SignatureAlgorithmProperty{}.Parse(t.SignatureAlgorithm)
	}
	if t.PublicKeyAlgorithm != "" {
		ct.cert.PublicKeyAlgorithm = pemproperties.PublicKeyAlgorithmProperty{}.Parse(t.PublicKeyAlgorithm)
	}
	if t.PublicKey != nil && t.PublicKey.PublicKey != "" {
		put := publicKeyType{}
		if err := put.UnmarshalBinary([]byte(t.PublicKey.PublicKey)); err != nil {
			log.Println(err)
		} else {
			ct.cert.PublicKey = put.puk
		}
	}
	if t.Version != 0 {
		ct.cert.Version = t.Version
	}
	if t.SerialNumber != 0 {
		ct.cert.SerialNumber = big.NewInt(t.SerialNumber)
	}

	if t.Issuer != nil {
		nt := dnameType{}
		nt.applyTemplate(*t.Issuer)
		ct.cert.Issuer = nt.dname
	}
	if t.Subject != nil {
		nt := dnameType{}
		nt.applyTemplate(*t.Subject)
		ct.cert.Subject = nt.dname
	}

	if t.NotBefore != "" {
		ct.cert.NotBefore = pemproperties.TimeProperty{}.Parse(t.NotBefore)
	}
	if t.NotAfter != "" {
		ct.cert.NotAfter = pemproperties.TimeProperty{}.Parse(t.NotAfter)
	}

	if t.KeyUsage != "" {
		ct.cert.KeyUsage = pemproperties.KeyUsageProperty{}.Parse(t.KeyUsage)
	}
	if len(t.ExtKeyUsage) > 0 {
		ct.cert.ExtKeyUsage = pemproperties.ExtKeyUsageProperty{}.ParseList(t.ExtKeyUsage)
	}
	//UnknownExtKeyUsage[] asn1.ObjectIdentifier // Encountered extended key usages unknown to this package.

	// Extensions contains raw X.509 extensions. When parsing certificates,
	// this can be used to extract non-critical extensions that are not
	// parsed by this package. When marshaling certificates, the Extensions
	// field is ignored, see ExtraExtensions.
	//Extensions[] pkix.Extension

	// ExtraExtensions contains extensions to be copied, raw, into any
	// marshaled certificates. Values override any extensions that would
	// otherwise be produced based on the other fields. The ExtraExtensions
	// field is not populated when parsing certificates, see Extensions.
	// ExtraExtensions[] pkix.Extension

	//UnhandledCriticalExtensions[] asn1.ObjectIdentifier

	ct.cert.BasicConstraintsValid = t.BasicConstraintsValid
	ct.cert.IsCA = t.IsCA
	ct.cert.MaxPathLen = t.MaxPathLen
	ct.cert.MaxPathLenZero = t.MaxPathLenZero

	if t.SubjectKeyId != "" {
		ct.cert.SubjectKeyId = []byte(t.SubjectKeyId)
	}
	if t.AuthorityKeyId != "" {
		ct.cert.AuthorityKeyId = []byte(t.AuthorityKeyId)
	}

	// RFC 5280, 4.2.2.1 (Authority Information Access)
	if len(t.OCSPServer) > 0 {
		ct.cert.OCSPServer = t.OCSPServer
	}
	if len(t.IssuingCertificateURL) > 0 {
		ct.cert.IssuingCertificateURL = t.IssuingCertificateURL
	}
	if len(t.DNSNames) > 0 {
		ct.cert.DNSNames = t.DNSNames
	}
	if len(t.EmailAddresses) > 0 {
		ct.cert.EmailAddresses = t.EmailAddresses
	}
	if len(t.IPAddresses) > 0 {
		ct.cert.IPAddresses = pemproperties.IPAddressListProperty{}.Parse(strings.Join(t.IPAddresses, ","))
	}
	if len(t.URIs) > 0 {
		ct.cert.URIs = pemproperties.URIListProperty{}.Parse(strings.Join(t.URIs, ","))
	}

	// Name constraints
	ct.cert.PermittedDNSDomainsCritical = t.PermittedDNSDomainsCritical
	if len(t.PermittedDNSDomains) > 0 {
		ct.cert.PermittedDNSDomains = t.PermittedDNSDomains
	}
	if len(t.ExcludedDNSDomains) > 0 {
		ct.cert.ExcludedDNSDomains = t.ExcludedDNSDomains
	}

	ct.cert.PermittedIPRanges = pemproperties.IPNetProperty{}.ParseList(t.PermittedIPRanges)
	ct.cert.ExcludedIPRanges = pemproperties.IPNetProperty{}.ParseList(t.ExcludedIPRanges)
	if len(t.PermittedEmailAddresses) > 0 {
		ct.cert.PermittedEmailAddresses = t.PermittedEmailAddresses
	}
	if len(t.ExcludedEmailAddresses) > 0 {
		ct.cert.ExcludedEmailAddresses = t.ExcludedEmailAddresses
	}
	if len(t.PermittedURIDomains) > 0 {
		ct.cert.PermittedURIDomains = t.PermittedURIDomains
	}
	if len(t.ExcludedURIDomains) > 0 {
		ct.cert.ExcludedURIDomains = t.ExcludedURIDomains
	}
	if len(t.CRLDistributionPoints) > 0 {
		ct.cert.CRLDistributionPoints = t.CRLDistributionPoints
	}
}
