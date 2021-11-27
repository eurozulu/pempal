package keytracker

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"pempal/keytools"
	"pempal/pemreader"
)

// idMatcher collects both Keys and certificates and attempts to match them using their public key or location
type idMatcher struct {
	keys    map[string]Key
	keyLocs map[string]Key
	certs   map[string][]*x509.Certificate
}

func (idc *idMatcher) AddCertificate(blk *pem.Block) (Identity, error) {
	if !keytools.CertificateTypes[blk.Type] {
		return nil, fmt.Errorf("%s is not a known certificate type", blk.Type)
	}
	c, err := x509.ParseCertificate(blk.Bytes)
	if err != nil {
		return nil, err
	}
	id := keytools.PublicKeySha1Hash(c.PublicKey)

	// Check if it matches known key
	k, ok := idc.keys[id]
	if ok {
		return NewIdentity(k, c), nil
	}
	// check if can match location
	loc, ok := readBlockHeader(pemreader.LocationHeaderKey, blk)
	if ok {
		k, ok := idc.keyLocs[trimLocation(loc)]
		if ok {
			return NewIdentity(k, c), nil
		}
	}
	// no matching key, place in the pile under its id
	idc.certs[id] = append(idc.certs[id], c)
	return nil, nil
}

func (idc *idMatcher) AddKey(k Key) []Identity {
	id := k.String()
	_, ok := idc.keys[id]
	if !ok {
		idc.keys[id] = k
	}
	loc := trimLocation(k.Location())
	_, ok = idc.keyLocs[loc]
	if !ok {
		idc.keyLocs[loc] = k
	}

	certs, ok := idc.certs[id]
	if !ok {
		return nil
	}
	ids := make([]Identity, len(certs))
	for i, c := range certs {
		ids[i] = NewIdentity(k, c)
	}
	return ids
}

func (idc *idMatcher) RemainingKeys() []Key {
	keys := make([]Key, len(idc.keys))
	var i int
	for _, k := range idc.keys {
		keys[i] = k
		i++
	}
	return keys
}
func (idc *idMatcher) RemainingCerts() []*x509.Certificate {
	var certs []*x509.Certificate
	for _, c := range idc.certs {
		certs = append(certs, c...)
	}
	return certs
}

func newIDCollector() *idMatcher {
	return &idMatcher{
		keys:  map[string]Key{},
		certs: map[string][]*x509.Certificate{},
	}
}
