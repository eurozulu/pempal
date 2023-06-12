package keys

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/resourceio"
	"strings"
)

type Keys interface {
	AllKeys(ctx context.Context) <-chan resourceio.ResourceLocation
	AllCertificates(ctx context.Context) <-chan resourceio.ResourceLocation
	KeyForIdentity(id model.Identity) (Key, error)
	CertificatesById(id model.Identity) []*x509.Certificate
	CertificateByName(dn pkix.Name) (*x509.Certificate, error)
	UserByName(dn pkix.Name) (User, error)
}

type keys struct {
	NonRecursive bool
	Keypath      []string
	Certpath     []string
}

func (k keys) AllKeys(ctx context.Context) <-chan resourceio.ResourceLocation {
	return resourceio.NewResourceScanner(!k.NonRecursive).Scan(ctx, k.Keypath...)
}
func (k keys) AllCertificates(ctx context.Context) <-chan resourceio.ResourceLocation {
	return resourceio.NewResourceScanner(!k.NonRecursive).Scan(ctx, k.Certpath...)
}

func (k keys) KeyForIdentity(id model.Identity) (Key, error) {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	for loc := range k.AllKeys(ctx) {
		prks := ParseKeyLocation(loc)
		for _, prk := range prks {
			if id.String() == prk.Identity().String() {
				return prk, nil
			}
		}
	}
	return nil, fmt.Errorf("no key found with id %s", id.String())
}

func (k keys) UserByName(dn pkix.Name) (User, error) {
	c, err := k.CertificateByName(dn)
	if err != nil {
		return nil, err
	}
	id, err := model.NewIdentity(c.PublicKey)
	if err != nil {
		return nil, err
	}
	prk, err := k.KeyForIdentity(id)
	if err != nil {
		return nil, err
	}
	return user{
		cert: c,
		key:  prk,
	}, nil
}

func (k keys) CertificateByName(dn pkix.Name) (*x509.Certificate, error) {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	for loc := range k.AllCertificates(ctx) {
		certs := filterCertificates(loc.Resources(model.Certificate), dn)
		if len(certs) == 0 {
			continue
		}
		return certs[0], nil
	}
	return nil, fmt.Errorf("no certificate found with name %s", dn.String())
}

func (k keys) CertificatesById(id model.Identity) []*x509.Certificate {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()

	var found []*x509.Certificate
	var emptyName pkix.Name
	for loc := range k.AllCertificates(ctx) {
		certs := filterCertificates(loc.Resources(model.Certificate), emptyName)
		if len(certs) == 0 {
			continue
		}
		for _, c := range certs {
			if c.PublicKey == nil {
				continue
			}
			cid, err := model.NewIdentity(c.PublicKey)
			if err != nil {
				logger.Error("%v", err)
				continue
			}
			if id.String() != cid.String() {
				continue
			}
			found = append(found, c)
		}
	}
	return found
}

func filterCertificates(res []model.Resource, dn pkix.Name) []*x509.Certificate {
	var certs []*x509.Certificate
	dns := dn.String()

	for _, r := range res {
		b, _ := pem.Decode([]byte(r.String()))
		c, err := x509.ParseCertificate(b.Bytes)
		if err != nil {
			logger.Error("Failed to parse certificate %v", err)
			continue
		}
		if dns != "" && !strings.Contains(c.Subject.String(), dns) {
			continue
		}
		certs = append(certs, c)
	}
	return certs
}

func ParseKeyLocation(l resourceio.ResourceLocation) []Key {
	var found []Key
	prks := l.Resources(model.PrivateKey)

	for _, prk := range prks {
		b, _ := pem.Decode([]byte(prk.String()))
		var k Key
		var err error
		if !x509.IsEncryptedPEMBlock(b) {
			k, err = NewKey(b)
		} else {
			if len(prks) > 1 {
				logger.Warning("ignoring encrypted key as multiple private keys found in %s", l.Location())
				continue
			}
			// Single key, look for corrisponding public key in same location
			puks := l.Resources(model.PublicKey)
			if len(puks) == 0 {
				logger.Warning("ignoring encrypted key as no public key found in %s", l.Location())
				continue
			}
			if len(puks) > 1 {
				logger.Warning("ignoring encrypted key as multiple public keys found in %s", l.Location())
				continue
			}
			k, err = NewKeyPair(b, puks[0])
		}
		if err != nil {
			logger.Error("ignoring key in %s  %v", l.Location(), err)
			continue
		}
		found = append(found, k)
	}
	return found
}

func NewKeys(keypath, certpath []string) Keys {
	return &keys{
		Keypath:  keypath,
		Certpath: certpath,
	}
}
