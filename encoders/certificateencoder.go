package encoders

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"pempal/pemtypes"
	"pempal/templates"
)

type certificateEncoder struct {
}

func (ce certificateEncoder) Encode(p *pem.Block) (templates.Template, error) {
	pt := pemtypes.ParsePEMType(p.Type)
	if pt != pemtypes.Certificate {
		return nil, fmt.Errorf("%s cannot be encoded into a certificate", p.Type)
	}
	cert, err := x509.ParseCertificate(p.Bytes)
	if err != nil {
		return nil, err
	}
	var t templates.CertificateTemplate
	ce.ApplyPem(cert, &t)
	return &t, nil
}

func (ce certificateEncoder) ApplyPem(cert *x509.Certificate, t *templates.CertificateTemplate) {
	t.Signature = hex.EncodeToString(cert.Signature)
	t.SignatureAlgorithm = cert.SignatureAlgorithm.String()

	t.SignatureAlgorithm = cert.SignatureAlgorithm.String()
	t.PublicKeyAlgorithm = cert.PublicKeyAlgorithm.String()

	if cert.PublicKey != nil {
		t.PublicKey = MarshalPublicKey(cert.PublicKey)
	}

	t.Version = cert.Version
	t.SerialNumber = cert.SerialNumber.Int64()

	nc := nameEncoder{}
	t.Issuer = &templates.NameTemplate{}
	nc.ApplyPem(&cert.Issuer, t.Issuer)
	t.Subject = &templates.NameTemplate{}
	nc.ApplyPem(&cert.Subject, t.Subject)

	t.NotBefore = MarshalTime(cert.NotBefore)
	t.NotAfter = MarshalTime(cert.NotAfter)
	t.KeyUsage = MarshalKeyUsage(cert.KeyUsage)

	// TODO
	//Extensions []pkix.Extension
	//ExtraExtensions []pkix.Extension
	//UnhandledCriticalExtensions []asn1.ObjectIdentifier

	t.ExtKeyUsage = MarshalExtKeyUsage(cert.ExtKeyUsage)
	//UnknownExtKeyUsage []asn1.ObjectIdentifier // Encountered extended key usages unknown to this package.

	t.BasicConstraintsValid = cert.BasicConstraintsValid
	t.IsCA = cert.IsCA

	t.MaxPathLen = cert.MaxPathLen
	t.MaxPathLenZero = cert.MaxPathLenZero

	t.SubjectKeyId = hex.EncodeToString(cert.SubjectKeyId)
	t.AuthorityKeyId = hex.EncodeToString(cert.AuthorityKeyId)

	// RFC 5280, 4.2.2.1 (Authority Information Access)
	t.OCSPServer = cert.OCSPServer
	t.IssuingCertificateURL = cert.IssuingCertificateURL

	if len(cert.DNSNames) > 0 {
		t.DNSNames = cert.DNSNames
	}
	if len(cert.EmailAddresses) > 0 {
		t.EmailAddresses = cert.EmailAddresses
	}
	if len(cert.IPAddresses) > 0 {
		t.IPAddresses = MarshalIPAddresses(cert.IPAddresses)
	}
	if len(cert.URIs) > 0 {
		t.URIs = MarshalURIs(cert.URIs)
	}

	t.PermittedDNSDomainsCritical = cert.PermittedDNSDomainsCritical
	t.PermittedDNSDomains = cert.PermittedDNSDomains
	t.ExcludedDNSDomains = cert.ExcludedDNSDomains
	//t.PermittedIPRanges = MarshalIPAddresses(cert.PermittedIPRanges)
	//t.ExcludedIPRanges            []*net.IPNet
	t.PermittedEmailAddresses = cert.PermittedEmailAddresses
	t.ExcludedEmailAddresses = cert.ExcludedEmailAddresses
	t.PermittedURIDomains = cert.PermittedURIDomains
	t.ExcludedURIDomains = cert.ExcludedURIDomains
	t.CRLDistributionPoints = cert.CRLDistributionPoints

}
