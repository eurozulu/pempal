package resources

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/model"
)

type Issuers interface {
	Issuers() []*x509.Certificate
	IssuersWithKeys() []*x509.Certificate
}

type issuers struct {
	certpath []string
}

func (isu issuers) Issuers() []*x509.Certificate {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	scan := NewPemScan(model.Certificate)
	var found []*x509.Certificate
	for certPems := range scan.ScanPath(ctx, isu.certpath...) {
		certs := isu.filterCACertificates(certPems.Content)
		if len(certs) == 0 {
			continue
		}
		found = append(found, certs...)
	}
	found = FilterExpiredCertificates(found)

	if len(found) == 0 {
		return nil
	}
	return found
}

func (isu issuers) IssuersWithKeys() []*x509.Certificate {
	issuerCerts := isu.Issuers()
	if len(issuerCerts) == 0 {
		return nil
	}

	keys := NewKeys(isu.certpath)
	knownKeyIDs, err := keys.PrivateKeyIDs()
	if err != nil {
		logging.Error("issuers", "Failed to read private key IDs %v", err)
		return nil
	}

	var found []*x509.Certificate
	for _, cert := range issuerCerts {
		id, err := model.NewKeyIdFromKey(cert.PublicKey)
		if err != nil {
			logging.Error("issuers", "Failed to read certificate key id %v", err)
			continue
		}
		if !ContainsID(knownKeyIDs, id) {
			continue
		}
		found = append(found, cert)
	}
	found = FilterExpiredCertificates(found)

	if len(found) == 0 {
		return nil
	}
	return found
}

func (isu issuers) filterCACertificates(pems []*pem.Block) []*x509.Certificate {
	var found []*x509.Certificate
	for _, pem := range pems {
		cert, err := x509.ParseCertificate(pem.Bytes)
		if err != nil {
			logging.Error("certificates", "failed to parse as certificate. %v", err)
			continue
		}
		if !cert.IsCA {
			continue
		}
		found = append(found, cert)
	}
	return found
}

func ContainsID(ids []model.KeyId, id model.KeyId) bool {
	for _, i := range ids {
		if i.Equals(id) {
			return true
		}
	}
	return false
}

func NewIssuers(certpath []string) Issuers {
	return &issuers{certpath}
}
