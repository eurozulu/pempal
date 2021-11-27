package keytracker

import (
	"crypto/x509"
	"encoding/pem"
	"pempal/keytools"
	"pempal/pemreader"
	"reflect"
)

// idMatcher collects both Keys and certificates and attempts to match them using their public key or location
type idMatcher struct {
	keys    map[string]Key
	keyLocs map[string]Key
	certs   map[string][]*pem.Block

	keyMatcher *keyMatcher
}

func (idc *idMatcher) AddPem(blk *pem.Block) ([]Identity, error) {
	if keytools.CertificateTypes[blk.Type] {
		return idc.addCertificate(blk)
	}
	// Assume its a key block
	k := idc.keyMatcher.AddPem(blk)
	if k == nil || reflect.ValueOf(k).IsNil() {
		return nil, nil
	}
	return idc.addKey(k)
}

func (idc *idMatcher) addCertificate(blk *pem.Block) ([]Identity, error) {
	c, err := x509.ParseCertificate(blk.Bytes)
	if err != nil {
		return nil, err
	}
	id := keytools.PublicKeySha1Hash(c.PublicKey)

	// Check if it matches known key
	k, ok := idc.keys[id]
	if ok {
		nid, err := NewIdentity(k, blk)
		if err != nil {
			return nil, err
		}
		return []Identity{nid}, nil
	}
	// check if can match location
	loc, _ := readBlockHeader(pemreader.LocationHeaderKey, blk)
	k, ok = idc.keyLocs[trimLocation(loc)]
	if ok {
		nid, err := NewIdentity(k, blk)
		if err != nil {
			return nil, err
		}
		return []Identity{nid}, nil
	}

	// no matching key, place in the pile under its id
	idc.certs[id] = append(idc.certs[id], blk)
	return nil, nil
}

func (idc *idMatcher) addKey(k Key) ([]Identity, error) {
	// Add new key to both id and location maps
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
	// match with any existing certs
	certs, ok := idc.certs[id]
	if !ok {
		return nil, nil
	}
	ids := make([]Identity, len(certs))
	var err error
	for i, c := range certs {
		ids[i], err = NewIdentity(k, c)
		if err != nil {
			return nil, err
		}
	}
	return ids, nil
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
func (idc *idMatcher) RemainingCerts() []*pem.Block {
	var certs []*pem.Block
	for _, c := range idc.certs {
		certs = append(certs, c...)
	}
	return certs
}

func newIDMatcher() *idMatcher {
	return &idMatcher{
		keys:       map[string]Key{},
		keyLocs:    map[string]Key{},
		certs:      map[string][]*pem.Block{},
		keyMatcher: newKeyMatcher(),
	}
}
