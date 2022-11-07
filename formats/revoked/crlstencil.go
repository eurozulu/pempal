package revoked

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"pempal/formats"
	"pempal/formats/dname"
	"pempal/resources"
	"pempal/stencils"
	"pempal/templates"
)

type crlStencil struct {
}

func (st crlStencil) MakeTemplate(r resources.Resource) (templates.Template, error) {
	blk := r.Pem()
	if blk == nil {
		return nil, nil
	}
	crl, err := x509.ParseRevocationList(blk.Bytes)
	if err != nil {
		return nil, err
	}

	t := &templates.CRLTemplate{}
	st.copyToTemplate(t, crl)
	return t, nil
}

func (st crlStencil) copyToTemplate(t *templates.CRLTemplate, crl *x509.RevocationList) {
	t.Issuer = &templates.NameTemplate{}
	dname.nameStencil{}.copyToTemplate(t.Issuer, crl.Issuer)
	t.Number = crl.Number.Int64()
	t.ThisUpdate = formats.marshalTime(crl.ThisUpdate)
	t.NextUpdate = formats.marshalTime(crl.NextUpdate)
	t.SignatureAlgorithm = crl.SignatureAlgorithm.String()
	t.Signature = formats.hexByteFormat(crl.Signature).String()
	t.AuthorityKeyId = formats.hexByteFormat(crl.AuthorityKeyId).String()
	t.RevokedCertificates = marshalRevokedCertificates(crl.RevokedCertificates)
}

func marshalRevokedCertificates(certs []pkix.RevokedCertificate) []*templates.RevokedCertificateTemplate {
	var temps []*templates.RevokedCertificateTemplate
	rcs := stencils.revokedCertificateStencil{}
	for _, c := range certs {
		ct := &templates.RevokedCertificateTemplate{}
		rcs.copyToTemplate(ct, c)
		temps = append(temps, ct)
	}
	return temps
}
