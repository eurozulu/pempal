package certificate

import (
	"crypto/x509"
	"fmt"
	"pempal/formats/formathelpers"
	"pempal/resources"
	"pempal/stencils"
	"pempal/templates"
)

// certificateStencil copies a x509.Certificate into a new Template
type certificateStencil struct {
}

func (st certificateStencil) MakeTemplate(r resources.Resource) (templates.Template, error) {
	blk := r.Pem()
	if blk == nil {
		return nil, nil
	}
	cert, err := x509.ParseCertificate(blk.Bytes)
	if err != nil {
		return nil, err
	}
	t := &templates.CertificateTemplate{}
	err = st.copyToTemplate(t, cert)
	return t, err
}

func (st certificateStencil) copyToTemplate(t *templates.CertificateTemplate, cert *x509.Certificate) error {
	if cert.PublicKey != nil {
		puk, err := stencils.marshalPublicKey(cert.PublicKey)
		if err != nil {
			return fmt.Errorf("failed to read public key for %s  %v", cert.Subject.CommonName, err)
		}
		t.PublicKey = puk
	}

	t.Version = cert.Version
	t.SerialNumber = cert.SerialNumber.Int64()
	t.Signature = formathelpers.HexByteFormat(cert.Signature).String()
	t.SignatureAlgorithm = cert.SignatureAlgorithm.String()
	t.PublicKeyAlgorithm = cert.PublicKeyAlgorithm.String()
	ns := stencils.nameStencil{}
	t.Subject = &templates.NameTemplate{}
	ns.copyToTemplate(t.Subject, cert.Subject)
	t.Issuer = &templates.NameTemplate{}
	ns.copyToTemplate(t.Issuer, cert.Issuer)

	//t.Extensions            []pkix.Extension
	//t.ExtraExtensions       []pkix.Extension
	t.NotBefore = stencils.marshalTime(cert.NotBefore)
	t.NotAfter = stencils.marshalTime(cert.NotAfter)
	t.KeyUsage = formathelpers.KeyUsageString(cert.KeyUsage)
	t.ExtKeyUsage = stencils.marshalExtKeyUsage(cert.ExtKeyUsage)
	t.BasicConstraintsValid = cert.BasicConstraintsValid
	t.IsCA = cert.IsCA
	t.MaxPathLen = cert.MaxPathLen
	t.MaxPathLenZero = cert.MaxPathLenZero
	t.SubjectKeyId = formathelpers.HexByteFormat(cert.SubjectKeyId).String()
	t.AuthorityKeyId = formathelpers.HexByteFormat(cert.AuthorityKeyId).String()
	t.OCSPServer = cert.OCSPServer
	t.IssuingCertificateURL = cert.IssuingCertificateURL
	t.DNSNames = cert.DNSNames
	t.EmailAddresses = cert.EmailAddresses
	t.IPAddresses = stencils.marshalIPAddresses(cert.IPAddresses)
	t.URIs = stencils.marshalURIs(cert.URIs)
	t.PermittedDNSDomainsCritical = cert.PermittedDNSDomainsCritical
	t.PermittedDNSDomains = cert.PermittedDNSDomains
	t.ExcludedDNSDomains = cert.ExcludedDNSDomains
	//t.PermittedIPRanges = cert.PermittedIPRanges
	//t.ExcludedIPRanges = cert.ExcludedIPRanges
	t.PermittedEmailAddresses = cert.PermittedEmailAddresses
	t.ExcludedEmailAddresses = cert.ExcludedEmailAddresses
	t.PermittedURIDomains = cert.PermittedURIDomains
	t.ExcludedURIDomains = cert.ExcludedURIDomains
	t.CRLDistributionPoints = cert.CRLDistributionPoints
	return nil
}
