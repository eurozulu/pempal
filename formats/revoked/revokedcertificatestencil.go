package revoked

import (
	"crypto/x509/pkix"
	"pempal/templates"
)

type revokedCertificateStencil struct {
}

func (st revokedCertificateStencil) copyToTemplate(t *templates.RevokedCertificateTemplate, cert pkix.RevokedCertificate) {
}
