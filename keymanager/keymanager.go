package keymanager

import (
	"context"
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"pempal/logger"
	"pempal/model"
	"pempal/resourceio"
)

type KeyManager interface {
	PublicKeys() []Identity
	PrivateKeys() map[Identity]crypto.PrivateKey

	PrivateKey(id Identity) (crypto.PrivateKey, error)
	CertificatesForIdentity(id Identity) []*x509.Certificate
	CertificateForDN(dn pkix.Name) *x509.Certificate
	User(dn pkix.Name) (User, error)
	Issuers() []User
}

type keyManager struct {
	keypath  string
	certpath string
}

func (km keyManager) PublicKeys() []Identity {
	keys := km.PrivateKeys()
	ids := make([]Identity, len(keys))
	var index int
	for id := range keys {
		ids[index] = id
		index++
	}
	return ids
}

func (km keyManager) PrivateKeys() map[Identity]crypto.PrivateKey {
	keyLocs := resourceio.NewResourceScanner(false,
		resourceio.NewTypeResourceLocationFilter(model.PrivateKey)).Scan(context.Background(), km.keypath)

	keys := map[Identity]crypto.PrivateKey{}
	for keyLoc := range keyLocs {
		keyDTOs, err := privateKeysFromLocation(keyLoc)
		if err != nil {
			logger.Log(logger.Error, "failed to read private keys from %s  %v", keyLoc.Path, err)
			continue
		}

		for _, dto := range keyDTOs {
			prk, _ := dto.ToPrivateKey()
			puk, _ := dto.ToPublicKey()
			id, err := NewIdentity(puk)
			if err != nil {
				logger.Log(logger.Error, "Failed to parse public key into ID  %v", err)
				continue
			}
			keys[id] = prk
		}
	}
	return keys
}

func (km keyManager) PrivateKey(id Identity) (crypto.PrivateKey, error) {
	prk, ok := km.PrivateKeys()[id]
	if !ok {
		return nil, fmt.Errorf("Failed to locate private key for id %s", id.String())
	}
	return prk, nil
}

func (km keyManager) CertificatesForIdentity(id Identity) []*x509.Certificate {
	certLocs := resourceio.NewResourceScanner(false,
		resourceio.NewTypeResourceLocationFilter(model.Certificate)).Scan(context.Background(), km.certpath)
	var certs []*x509.Certificate
	for loc := range certLocs {
		certs = append(certs, certificatesByIdFromLocation(id, loc)...)
	}
	return certs
}

func (km keyManager) CertificateForDN(dn pkix.Name) *x509.Certificate {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()

	dns := dn.String()
	certLocs := resourceio.NewResourceScanner(false,
		resourceio.NewTypeResourceLocationFilter(model.Certificate)).Scan(ctx, km.certpath)
	for loc := range certLocs {
		if cert := findCertificateWithDN(dns, loc); cert != nil {
			return cert
		}
	}
	return nil
}

func (km keyManager) User(dn pkix.Name) (User, error) {
	logger.Log(logger.Debug, "searching for user %s", dn.String())
	cert := km.CertificateForDN(dn)
	if cert == nil {
		return nil, fmt.Errorf("user %s is not known.  A certificate matching that DN could not be found", dn.String())
	}
	logger.Log(logger.Debug, "certificate, serial number %v for user %s found", cert.SerialNumber, dn.String())

	id, err := NewIdentity(cert.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the user %s public key  %v", dn.String(), err)
	}
	logger.Log(logger.Debug, "user identity: %s", id.String())

	prk, err := km.PrivateKey(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find private key for user %s  %v", dn.String(), err)
	}
	return &user{
		id:   id,
		cert: cert,
		key:  prk,
	}, nil
}

func (km keyManager) Issuers() []User {
	keys := km.PrivateKeys()
	var issuers []User
	for id, key := range keys {
		certs := caCertificates(km.CertificatesForIdentity(id))
		for _, cert := range certs {
			issuers = append(issuers, user{
				id:   id,
				cert: cert,
				key:  key,
			})
		}
	}
	return issuers
}

func privateKeysFromLocation(loc *resourceio.ResourceLocation) ([]*model.PrivateKeyDTO, error) {
	var keys []*model.PrivateKeyDTO
	for _, r := range loc.Resources {
		if r.ResourceType() != model.PrivateKey {
			continue
		}
		keydto, err := keyFromResource(r)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse private key at location %s  %v", loc.Path, err)

		}
		keys = append(keys, keydto)
	}
	return keys, nil

}

func certificatesByIdFromLocation(id Identity, loc *resourceio.ResourceLocation) []*x509.Certificate {
	var certs []*x509.Certificate
	for _, r := range loc.Resources {
		if r.ResourceType() != model.Certificate {
			continue
		}
		cert, err := certificateFromResource(r)
		if err != nil {
			logger.Log(logger.Error, "Failed to parse certificate at location %s  %v", loc.Path, err)
			continue
		}
		cid, err := NewIdentity(cert.PublicKey)
		if err != nil {
			logger.Log(logger.Error, "Failed to parse public key for certificate %s  %v", cert.Subject.String(), err)
			continue
		}
		if id != cid {
			continue
		}
		certs = append(certs, cert)
	}
	return certs
}

func keyFromResource(r model.PEMResource) (*model.PrivateKeyDTO, error) {
	dto, err := model.DTOForResource(r)
	if err != nil {
		return nil, err
	}
	return dto.(*model.PrivateKeyDTO), nil
}

func certificateFromResource(r model.PEMResource) (*x509.Certificate, error) {
	dto, err := model.DTOForResource(r)
	if err != nil {
		return nil, err
	}
	certDto := dto.(*model.CertificateDTO)
	return certDto.ToCertificate()
}

func findCertificateWithDN(dn string, loc *resourceio.ResourceLocation) *x509.Certificate {
	for _, r := range loc.Resources {
		if r.ResourceType() != model.Certificate {
			continue
		}
		cert, err := certificateFromResource(r)
		if err != nil {
			logger.Log(logger.Error, "Failed to parse certificate at location %s  %v", loc.Path, err)
			continue
		}
		if cert.Subject.String() != dn {
			continue
		}
		return cert
	}
	return nil
}

func caCertificates(certs []*x509.Certificate) []*x509.Certificate {
	var cas []*x509.Certificate
	for _, c := range certs {
		if !c.IsCA {
			continue
		}
		cas = append(cas, c)
	}
	return cas
}
