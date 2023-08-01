package identity

import (
	"context"
	"fmt"
	"github.com/eurozulu/pempal/logger"
)

type Users interface {
	AllUsers(ctx context.Context) <-chan User
	UsersByIdentity(id Identity) []User
	UserByName(name string) (User, error)
	Keys() Keys
	Certificates() Certificates
}

type users struct {
	keyz  Keys
	certz Certificates
}

func (is users) Certificates() Certificates {
	return is.certz
}

func (is users) Keys() Keys {
	return is.keyz
}

func (is users) UsersByIdentity(id Identity) []User {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	var found []User
	for is := range is.AllUsers(ctx) {
		if is.Identity().String() != id.String() {
			continue
		}
		found = append(found, is)
	}
	return found
}

func (is users) UserByName(name string) (User, error) {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()

	for is := range is.AllUsers(ctx) {
		if compareDN(is.Certificate().Certificate().Subject, name) {
			return is, nil
		}
	}
	return nil, fmt.Errorf("%s unknown user", name)
}

func (is users) AllUsers(ctx context.Context) <-chan User {
	ch := make(chan User)
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
				case ch <- &user{
					key:  key,
					cert: c,
				}:
				}
			}
		}
	}()
	return ch
}

func NewIssuers(keypath, certpath []string) Users {
	keyz := NewKeys(keypath)
	certz := NewCertificates(certpath)
	return &users{
		keyz:  keyz,
		certz: certz,
	}
}
