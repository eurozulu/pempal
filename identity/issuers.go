package identity

import (
	"context"
	"fmt"
	"github.com/eurozulu/pempal/logger"
)

type Issuers interface {
	AllIssuers(ctx context.Context) <-chan Issuer
	IssuersByIdentity(id Identity) []Issuer
	IssuerByName(name string) (Issuer, error)
	Keys() Keys
	Certificates() Certificates
}

type issuers struct {
	keyz  Keys
	certz Certificates
}

func (is issuers) Certificates() Certificates {
	return is.certz
}

func (is issuers) Keys() Keys {
	return is.keyz
}

func (is issuers) IssuersByIdentity(id Identity) []Issuer {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	var found []Issuer
	for is := range is.AllIssuers(ctx) {
		if is.Identity().String() != id.String() {
			continue
		}
		found = append(found, is)
	}
	return found
}

func (is issuers) IssuerByName(name string) (Issuer, error) {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()

	for is := range is.AllIssuers(ctx) {
		if compareDN(is.Certificate().Certificate().Subject, name) {
			return is, nil
		}
	}
	return nil, fmt.Errorf("%s unknown issuer", name)
}

func (is issuers) AllIssuers(ctx context.Context) <-chan Issuer {
	ch := make(chan Issuer)
	go func() {
		for key := range is.keyz.AllKeys(ctx) {
			certs, err := is.certz.CertificatesByIdentity(key.Identity())
			if err != nil {
				logger.Error("failed to read certificates for key '%s' %v", key.Identity().String(), err)
				continue
			}
			for _, c := range certs {
				select {
				case <-ctx.Done():
					return
				case ch <- &issuer{
					key:  key,
					cert: c,
				}:
				}
			}
		}
	}()
	return ch
}

func NewIssuers(keypath, certpath []string) Issuers {
	keyz := NewKeys(keypath)
	certz := NewCertificates(certpath)
	return &issuers{
		keyz:  keyz,
		certz: certz,
	}
}
