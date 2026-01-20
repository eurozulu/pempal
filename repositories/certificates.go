package repositories

import (
	"context"
	"fmt"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/resourcefiles"
	"math/big"
)

type Certificates string

type CertificateFilter func(*model.Certificate) bool

func (certs Certificates) ByName(dn model.DistinguishedName) (*model.Certificate, error) {
	return certs.FindFirst(func(cert *model.Certificate) bool {
		return model.DistinguishedName(cert.Subject).Equals(dn)
	})
}

func (certs Certificates) MatchByName(name string) ([]*model.Certificate, error) {
	dn, err := model.ParseName(name)
	if err != nil {
		return nil, err
	}
	return certs.FindAll(func(cert *model.Certificate) bool {
		return model.DistinguishedName(cert.Subject).Matches(*dn)
	}), nil
}

func (certs Certificates) ByPublicKey(puk model.PublicKey) []*model.Certificate {
	pukS := puk.String()
	return certs.FindAll(func(cert *model.Certificate) bool {
		key := model.NewPublicKey(cert.PublicKey)
		return key.String() == pukS
	})
}

func (certs Certificates) ByIssuer(dn model.DistinguishedName) []*model.Certificate {
	return certs.FindAll(func(cert *model.Certificate) bool {
		return model.DistinguishedName(cert.Issuer).Equals(dn)
	})
}

func (certs Certificates) MatchByIssuer(name string) ([]*model.Certificate, error) {
	dn, err := model.ParseName(name)
	if err != nil {
		return nil, err
	}
	return certs.FindAll(func(cert *model.Certificate) bool {
		return model.DistinguishedName(cert.Issuer).Matches(*dn)
	}), nil
}

func (certs Certificates) ByCA() []*model.Certificate {
	return certs.FindAll(func(cert *model.Certificate) bool {
		return cert.IsCA
	})
}

func (certs Certificates) BySerialNumber(n *model.SerialNumber) (*model.Certificate, error) {
	ni := (*big.Int)(n)
	return certs.FindFirst(func(certificate *model.Certificate) bool {
		return ni.Cmp(certificate.SerialNumber) == 0
	})
}

func (certs Certificates) ByFingerPrint(fingerPrint model.Fingerprint) (*model.Certificate, error) {
	fp := fingerPrint.String()
	return certs.FindFirst(func(cert *model.Certificate) bool {
		return cert.Fingerprint().Equals(fp)
	})
}

func (certs Certificates) MatchByFingerPrint(fingerPrint string) []*model.Certificate {
	return certs.FindAll(func(cert *model.Certificate) bool {
		return cert.Fingerprint().Match(fingerPrint)
	})
}

func (certs Certificates) FindFirst(filter CertificateFilter) (*model.Certificate, error) {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	found := certs.Find(ctx, filter)
	cert, ok := <-found
	if !ok {
		return nil, fmt.Errorf("no certificates found")
	}
	return cert, nil
}

func (certs Certificates) FindAll(filter CertificateFilter) []*model.Certificate {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	var found []*model.Certificate
	for cert := range certs.Find(ctx, filter) {
		found = append(found, cert)
	}
	return found
}

func (certs Certificates) Find(ctx context.Context, filter CertificateFilter) <-chan *model.Certificate {
	ch := make(chan *model.Certificate)
	go func() {
		defer close(ch)
		certFiles := resourcefiles.PemFiles(string(certs)).FindByType(ctx, model.ResourceTypeCertificate)
		for pf := range certFiles {

			for _, res := range pf.Resources() {
				c, ok := res.(*model.Certificate)
				if !ok {
					continue
				}
				if filter != nil && !filter(c) {
					continue
				}
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
