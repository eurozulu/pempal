package resourcefiles

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/model"
)

type Derfileformat struct{}

// TODO: Check if parsing correctly (See DerFormat)
func (d Derfileformat) Format(data []byte) ([]*pem.Block, error) {
	if certs, err := x509.ParseCertificates(data); err == nil && len(certs) > 0 {
		blks := make([]*pem.Block, len(certs))
		for i, cert := range certs {
			blks[i] = &pem.Block{
				Type:  model.ResourceTypeCertificate.String(),
				Bytes: cert.Raw,
			}
		}
		return blks, nil
	}

	if _, err := x509.ParseCertificate(data); err == nil {
		return []*pem.Block{{
			Type:  model.ResourceTypeCertificate.String(),
			Bytes: data,
		}}, nil
	}

	if _, err := x509.ParsePKIXPublicKey(data); err == nil {
		return []*pem.Block{{
			Type:  model.ResourceTypePublicKey.String(),
			Bytes: data,
		}}, nil
	}
	if _, err := x509.ParsePKCS8PrivateKey(data); err == nil {
		return []*pem.Block{{
			Type:  model.ResourceTypePrivateKey.String(),
			Bytes: data,
		}}, nil

	}
	if _, err := x509.ParseCertificateRequest(data); err == nil {
		return []*pem.Block{{
			Type:  model.ResourceTypeCertificateRequest.String(),
			Bytes: data,
		}}, nil

	}
	if _, err := x509.ParseRevocationList(data); err == nil {
		return []*pem.Block{{
			Type:  model.ResourceTypeRevokationList.String(),
			Bytes: data,
		}}, nil
	}
	return nil, fmt.Errorf("unknown DER format")
}
