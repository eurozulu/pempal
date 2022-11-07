package certificate

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
	"math/big"
	"net/http"
	"pempal/formats"
	"pempal/formats/dname"
	"pempal/formats/formathelpers"
	"pempal/resources"
	"pempal/templates"
	"time"
)

// certificateBuilder formats a new, signed x509.Certificate based on the given template data
type certificateBuilder struct {
	certTemp templates.CertificateTemplate
	location string
}

func (fm certificateBuilder) Template() templates.Template {
	return fm.certTemp
}

func (fm *certificateBuilder) SetLocation(l string) {
	fm.location = l
}

func (fm certificateBuilder) AddTemplate(ts ...templates.Template) error {
	for _, t := range ts {
		by, err := yaml.Marshal(t)
		if err != nil {
			return err
		}
		if err = yaml.Unmarshal(by, &fm.certTemp); err != nil {
			return err
		}
	}
	return nil
}

func (fm certificateBuilder) Validate() []error {
	var names []string
	if fm.certTemp.PublicKey == "" {
		names = append(names, "public-key")
	}
	if fm.certTemp.Signature == "" {
		names = append(names, "signature")
	}
	if fm.certTemp.Subject == nil {
		names = append(names, "subject")
	}
	if fm.certTemp.Subject.CommonName == "" {
		names = append(names, "subject.common-name")
	}
	if fm.certTemp.Issuer == nil {
		names = append(names, "issuer")
	}
	if fm.certTemp.NotAfter == "" {
		names = append(names, "not-after")
	}
	if len(names) == 0 {
		return nil
	}
	errs := make([]error, len(names))
	for i, n := range names {
		errs[i] = fmt.Errorf("missing %s", n)
	}
	return errs
}

func (fm certificateBuilder) Build() (resources.Resources, error) {
	if errs := fm.Validate(); len(errs) > 0 {
		return nil, fmt.Errorf("%v", errs)
	}

	// create the template certificate
	var cert x509.Certificate
	fm.ApplyTemplate(&cert)

	// locate the public key the request was signed by
	cert.PublicKey, err = formathelpers.ResolvePublicKey(ct.PublicKey)
	if err != nil {
		return nil, err
	}

	// Establish issuer.
	var issuer *x509.Certificate
	if ct.Issuer.IsEmpty() {
		// empty issuer == self-signed / root
		issuer = &cert
	} else {
		issuer, err = formats.resolveIssuer(*ct.Issuer, ct.IssuingCertificateURL)
		if err != nil {
			return nil, err
		}
	}

	// locate issuer private key
	prk, err := formats.matchPrivateKeyFromPublic(issuer.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("missing issuer private key for %s", issuer.Subject.CommonName)
	}

	by, err := x509.CreateCertificate(rand.Reader, &cert, issuer, cert.PublicKey, prk)
	if err != nil {
		return nil, err
	}
	return resources.NewResource("", &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: by,
	}), nil
}

func (fm certificateBuilder) ApplyTemplate(cert *x509.Certificate) error {
	if fm.certTemp.Version != 0 {
		cert.Version = fm.certTemp.Version
	}

	if fm.certTemp.PublicKeyAlgorithm != "" {
		cert.PublicKeyAlgorithm = formathelpers.ParsePublicKeyAlgorithm(fm.certTemp.PublicKeyAlgorithm)
	}

	if fm.certTemp.SignatureAlgorithm != "" {
		cert.SignatureAlgorithm = formathelpers.ParseSignatureAlgorithm(fm.certTemp.SignatureAlgorithm)
	}
	var err error
	if fm.certTemp.Signature != "" {
		cert.Signature, err = formathelpers.ParseHexBytes(fm.certTemp.Signature)
		if err != nil {
			return err``
		}
	}
	if fm.certTemp.Subject != nil {
		dname.NewNameBuilder(fm.certTemp.Subject).ApplyTemplate(&cert.Subject)
	}
	if fm.certTemp.Issuer != nil {
		dname.NewNameBuilder(fm.certTemp.Issuer).ApplyTemplate(&cert.Issuer)
	}
	if fm.certTemp.SerialNumber != 0 {
		cert.SerialNumber = big.NewInt(fm.certTemp.SerialNumber)
	}
	//Extensions            []pkix.Extension
	//ExtraExtensions       []pkix.Extension
	cert.NotBefore = time.Now()
	if fm.certTemp.NotBefore != "" {
		tm, err := http.ParseTime(fm.certTemp.NotBefore)
		if err != nil {
			return fmt.Errorf("failed to parse not-before time.")
		}
		if tm.Before(cert.NotBefore) {
			return fmt.Errorf("not-before time is in the past")
		}
		cert.NotBefore = tm
	}
	cert.NotAfter = time.Now()
	tm, err := http.ParseTime(fm.certTemp.NotAfter)
	if err != nil {
		return fmt.Errorf("failed to parse not-after time.")
	}
	if !tm.After(cert.NotBefore) {
		return fmt.Errorf("not-after time must be after the not-before time")
	}
	cert.NotAfter = tm
	if fm.certTemp.KeyUsage != "" {
		cert.KeyUsage = formathelpers.ParseKeyUsage(fm.certTemp.KeyUsage)
	}
	if len(fm.certTemp.ExtKeyUsage) > 0 {
		cert.ExtKeyUsage = formathelpers.ParseExtKeyUsage(fm.certTemp.ExtKeyUsage)
	}
	cert.BasicConstraintsValid = fm.certTemp.BasicConstraintsValid
	cert.IsCA = fm.certTemp.IsCA
	cert.MaxPathLen = fm.certTemp.MaxPathLen
	cert.MaxPathLen = fm.certTemp.MaxPathLen
	if fm.certTemp.SubjectKeyId != "" {
		by, err := formathelpers.ParseHexBytes(fm.certTemp.SubjectKeyId)
		if err != nil {
			cert.SubjectKeyId = by
		}
	}
	if fm.certTemp.AuthorityKeyId != "" {
		by, err := formathelpers.ParseHexBytes(fm.certTemp.AuthorityKeyId)
		if err != nil {
			cert.AuthorityKeyId = by
		}
	}
	cert.OCSPServer = fm.certTemp.OCSPServer
	cert.IssuingCertificateURL = fm.certTemp.IssuingCertificateURL
	cert.DNSNames = fm.certTemp.DNSNames
	cert.EmailAddresses = fm.certTemp.EmailAddresses
	cert.IPAddresses = formathelpers.ParseIPAddresses(fm.certTemp.IPAddresses)
	cert.URIs = formathelpers.ParseURIs(fm.certTemp.URIs)
	cert.PermittedDNSDomainsCritical = fm.certTemp.PermittedDNSDomainsCritical
	cert.PermittedDNSDomains = fm.certTemp.PermittedDNSDomains
	cert.ExcludedDNSDomains = fm.certTemp.ExcludedDNSDomains
	//cert.PermittedIPRanges           = parseIPAddresses(fm.certTemp.PermittedIPRanges)
	//cert.ExcludedIPRanges            = parseIPAddresses(fm.certTemp.ExcludedIPRanges)
	cert.PermittedEmailAddresses = fm.certTemp.PermittedEmailAddresses
	cert.ExcludedEmailAddresses = fm.certTemp.ExcludedEmailAddresses
	cert.PermittedURIDomains = fm.certTemp.PermittedURIDomains
	cert.ExcludedURIDomains = fm.certTemp.ExcludedURIDomains
	cert.CRLDistributionPoints = fm.certTemp.CRLDistributionPoints
	return nil
}
