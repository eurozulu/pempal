package repositories

import (
	"context"
	"fmt"
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/model"
)

type Issuers string

func (u Issuers) ByName(dn model.DistinguishedName) (*model.Issuer, error) {
	issuers := u.FindAll(func(certificate *model.Certificate) bool {
		return model.DistinguishedName(certificate.Subject).Equals(dn)
	})
	if len(issuers) == 0 {
		return nil, fmt.Errorf("no issue named %q found", dn)
	}
	if len(issuers) > 1 {
		return nil, fmt.Errorf("multiple issuers match the name %q", dn)
	}
	return issuers[0], nil
}

func (u Issuers) MatchByName(dn model.DistinguishedName) []*model.Issuer {
	return u.FindAll(func(certificate *model.Certificate) bool {
		return model.DistinguishedName(certificate.Subject).Matches(dn)
	})
}
func (u Issuers) FindAll(filter CertificateFilter) []*model.Issuer {
	var issuers []*model.Issuer
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	for issuer := range u.Find(ctx, filter) {
		issuers = append(issuers, issuer)
	}
	return issuers
}

func (u Issuers) Find(ctx context.Context, filter CertificateFilter) <-chan *model.Issuer {
	ch := make(chan *model.Issuer)
	go func() {
		defer close(ch)
		keyz := Keys(u)
		caFilter := func(certificate *model.Certificate) bool {
			if !certificate.IsCA {
				return false
			}
			if filter == nil {
				return true
			}
			return filter(certificate)
		}

		for certificate := range Certificates(u).Find(ctx, caFilter) {
			prk, err := keyz.ByPublicKey(model.NewPublicKey(certificate.PublicKey))
			if err != nil {
				logging.Debug("failed to lookup key for %s %v", certificate.Subject, err)
				continue
			}
			select {
			case <-ctx.Done():
			case ch <- model.NewIssuer(certificate, prk):
			}
		}
	}()
	return ch
}
