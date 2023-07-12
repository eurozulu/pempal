package identity

import (
	"context"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/resourceio"
	"github.com/eurozulu/pempal/resources"
	"strings"
)

type Certificates interface {
	AllCertificates(ctx context.Context) <-chan Certificate
	CertificatesByIdentity(id Identity) ([]Certificate, error)
	CertificatesByName(name string) ([]Certificate, error)
}

type certificates struct {
	NonRecursive bool
	CertPath     []string
}

func (cs certificates) AllCertificates(ctx context.Context) <-chan Certificate {
	ch := make(chan Certificate)
	go func() {
		for loc := range resourceio.NewResourceScanner(!cs.NonRecursive).Scan(ctx, cs.CertPath...) {
			certs, err := locationCertificates(loc)
			if err != nil {
				logger.Error("failed to read certificates  %v", err)
				continue
			}
			for _, c := range certs {
				select {
				case <-ctx.Done():
					return
				case ch <- c:
				}
			}
		}
	}()
	return ch
}

func (cs certificates) CertificatesByIdentity(id Identity) ([]Certificate, error) {
	var found []Certificate
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	ids := id.String()
	for c := range cs.AllCertificates(ctx) {
		if c.Identity().String() != ids {
			continue
		}
		found = append(found, c)
	}
	if len(found) == 0 {
		return nil, fmt.Errorf("no certificates found with id %s", ids)
	}
	return found, nil

}

func (cs certificates) CertificatesByName(name string) ([]Certificate, error) {
	var found []Certificate
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	for c := range cs.AllCertificates(ctx) {
		cert := c.Certificate()
		if cert == nil {
			continue
		}
		if !compareDN(cert.Subject, name) {
			continue
		}
		found = append(found, c)
	}
	if len(found) == 0 {
		return nil, fmt.Errorf("no certificates found with name %s", name)
	}
	return found, nil
}

func compareDN(dn pkix.Name, s string) bool {
	return strings.Contains(dn.String(), s)
}

func locationCertificates(loc resourceio.ResourceLocation) ([]Certificate, error) {
	var found []Certificate
	certRes := loc.Resources(resources.Certificate)
	for _, cr := range certRes {
		blk, _ := pem.Decode([]byte(cr.String()))
		if blk == nil {
			continue
		}
		found = append(found, &certificate{
			loc:  loc.Location(),
			cert: blk,
		})
	}
	if len(found) == 0 {
		return nil, fmt.Errorf("no certificates found in location %s", loc.Location())
	}
	return found, nil
}

func NewCertificates(certpath []string) Certificates {
	return &certificates{
		CertPath: certpath,
	}
}
