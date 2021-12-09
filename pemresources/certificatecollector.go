package pemresources

import (
	"encoding/pem"
	"pempal/fileformats"
)

type certificateCollector struct {
	keyCollection *keyTracker
	keys          map[string]*PrivateKey
	certs         map[string][]*Certificate
}

func (cc *certificateCollector) AddBlocks(blocks []*pem.Block) ([]*Identity, error) {
	var ids []*Identity
	for _, block := range blocks {
		if k, _ := cc.keyCollection.AddBlock(block); k != nil {
			ts, err := cc.addKey(k)
			if err != nil {
				return nil, err
			}
			ids = append(ids, ts...)
		}
		if !fileformats.PemTypesCertificate[block.Type] {
			return nil, nil
		}
		var cert Certificate
		if err := cert.UnmarshalPem(block); err != nil {
			return nil, err
		}
		k, ok := cc.keys[cert.PublicKeyHash]
		if ok {
			ids = append(ids, &Identity{
				Key:  k,
				Cert: &cert,
			})
		} else {
			cc.certs[cert.PublicKeyHash] = append(cc.certs[cert.PublicKeyHash], &cert)
		}
	}
	return ids, nil
}

func (cc *certificateCollector) addKey(k *PrivateKey) ([]*Identity, error) {
	cc.keys[k.PublicKeyHash] = k
	certs, ok := cc.certs[k.PublicKeyHash]
	if !ok {
		return nil, nil
	}
	ids := make([]*Identity, len(certs))
	for i, c := range certs {
		ids[i] = &Identity{
			Key:  k,
			Cert: c,
		}
	}
	delete(cc.certs, k.PublicKeyHash)
	return ids, nil
}

func newCertificateCollector() *certificateCollector {
	return &certificateCollector{
		keyCollection: newKeyTracker(),
		keys:          map[string]*PrivateKey{},
		certs:         map[string][]*Certificate{},
	}
}
