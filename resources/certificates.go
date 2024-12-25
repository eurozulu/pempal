package resources

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/model"
	"time"
)

type Certificates interface {
	CertificatesBySubject(name pkix.Name) []*x509.Certificate
	CertificatesById(id model.KeyId) []*x509.Certificate
}

type certificates struct {
	ExposeExpiredCertificates bool
	certpath                  []string
}

func (c certificates) CertificatesBySubject(name pkix.Name) []*x509.Certificate {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	scan := NewPemScan(model.Certificate)
	var found []*x509.Certificate
	for certPems := range scan.ScanPath(ctx, c.certpath...) {
		certs := c.filterCertificatesByName(name, certPems.Content)
		if len(certs) == 0 {
			continue
		}
		found = append(found, certs...)
	}
	if !c.ExposeExpiredCertificates {
		found = FilterExpiredCertificates(found)
	}
	if len(found) == 0 {
		return nil
	}
	return found
}

func (c certificates) CertificatesById(id model.KeyId) []*x509.Certificate {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	scan := NewPemScan(model.Certificate)
	var found []*x509.Certificate
	for certPems := range scan.ScanPath(ctx, c.certpath...) {
		certs := c.filterCertificatesByID(id, certPems.Content)
		if len(certs) == 0 {
			continue
		}
		found = append(found, certs...)
	}
	if !c.ExposeExpiredCertificates {
		found = FilterExpiredCertificates(found)
	}
	if len(found) == 0 {
		return nil
	}
	return found
}

func (c certificates) filterCertificatesByName(name pkix.Name, pems []*pem.Block) []*x509.Certificate {
	var found []*x509.Certificate
	dn := model.DistinguishedName(name)

	for _, pem := range pems {
		cert, err := x509.ParseCertificate(pem.Bytes)
		if err != nil {
			logging.Error("certificates", "failed to parse as certificate. %v", err)
			continue
		}
		if !dn.Equals(model.DistinguishedName(cert.Subject)) {
			continue
		}
		found = append(found, cert)
	}
	return found
}

func (c certificates) filterCertificatesByID(id model.KeyId, pems []*pem.Block) []*x509.Certificate {
	var found []*x509.Certificate
	for _, pem := range pems {
		cert, err := x509.ParseCertificate(pem.Bytes)
		if err != nil {
			logging.Error("certificates", "failed to parse as certificate. %v", err)
			continue
		}
		cid, err := model.NewKeyIdFromKey(cert.PublicKey)
		if err != nil {
			logging.Error("certificates", "failed to parse as key. %v", err)
			continue
		}
		if !id.Equals(cid) {
			continue
		}
		found = append(found, cert)
	}
	return found
}

func FilterExpiredCertificates(certs []*x509.Certificate) []*x509.Certificate {
	now := time.Now()
	var found []*x509.Certificate
	for _, cert := range certs {
		if cert.NotBefore.After(now) {
			continue
		}
		if cert.NotAfter.Before(now) {
			continue
		}
		found = append(found, cert)
	}
	return found
}

func NewCertificates(certpath []string) Certificates {
	return &certificates{certpath: certpath}
}
