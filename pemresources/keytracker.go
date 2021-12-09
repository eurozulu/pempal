package pemresources

import (
	"encoding/pem"
	"path"
	"pempal/fileformats"
	"strings"
)

// keyTracker collects public and private keys and attempts to match them together.
// private, encrypted keys (where public is not available until decrypted) are matched to public keys in two way:
// linked keys, public keys with a header identifying the private key it belongs to
// both private and public keys share the same location (ignoring any index and file extensoion)
type keyTracker struct {
	anonKeys   map[string]*PrivateKey
	publicKeys map[string]*PublicKey
	linkedKeys map[string]*PublicKey
}

func (kc keyTracker) AnonymousKeys() []*PrivateKey {
	keys := make([]*PrivateKey, len(kc.anonKeys))
	var i int
	for _, v := range kc.anonKeys {
		keys[i] = v
		i++
	}
	return keys
}

func (kc keyTracker) AddBlock(blk *pem.Block) (*PrivateKey, error) {
	if fileformats.PemTypesPrivateKey[blk.Type] {
		return kc.addPrivateKey(blk)
	}

	// try to make a public key out of block
	puk := &PublicKey{}
	if err := puk.UnmarshalPem(blk); err != nil {
		return nil, err
	}
	return kc.addPublicKey(puk)
}

func (kc keyTracker) addPrivateKey(blk *pem.Block) (*PrivateKey, error) {
	k := &PrivateKey{}
	if err := k.UnmarshalPem(blk); err != nil {
		return nil, err
	}
	// first check if it doesn't need matching (unencrypted)
	if k.PublicKey != "" {
		return k, nil
	}

	// unknown public key, attempt to match

	// look for a linked public key
	id := k.PublicKeyHash
	lk, ok := kc.linkedKeys[id]
	if ok {
		if err := k.LinkPublicKey(lk); err == nil {
			delete(kc.linkedKeys, id)
			return k, nil
		}
	}

	// Look for a puk sharing the same location
	loc := TrimLocation(k.Location)
	if puk, ok := kc.publicKeys[loc]; ok {
		if err := k.LinkPublicKey(puk); err == nil {
			delete(kc.publicKeys, loc)
			return k, nil
		}
	}
	// no matching public key, add to anon under location
	kc.anonKeys[loc] = k
	return nil, nil
}

func (kc keyTracker) addPublicKey(puk *PublicKey) (*PrivateKey, error) {
	// can we match by location
	loc := TrimLocation(puk.Location)
	if k, ok := kc.anonKeys[loc]; ok {
		if err := k.LinkPublicKey(puk); err != nil {
			delete(kc.anonKeys, loc)
			return k, nil
		}
	}

	// attempt to match a linked key
	if puk.LinkedId != "" {
		if kk := kc.findLinkedKey(puk.LinkedId); kk != "" {
			k := kc.anonKeys[kk]
			if err := k.LinkPublicKey(puk); err == nil {
				delete(kc.anonKeys, kk)
				return k, nil
			}
		}
		// no match to private, for linked key, map by its id
		kc.linkedKeys[puk.LinkedId] = puk
		return nil, nil
	}

	// no match to any key
	kc.publicKeys[loc] = puk
	return nil, nil
}

func (kc keyTracker) findLinkedKey(id string) string {
	for k, v := range kc.anonKeys {
		if v.PublicKeyHash == id {
			return k
		}
	}
	return ""
}

func TrimLocation(l string) string {
	l = strings.TrimRight(l, "0123456789:") // remove any index
	return l[:len(l)-len(path.Ext(l))]      // remove the file extension
}

func newKeyTracker() *keyTracker {
	return &keyTracker{
		anonKeys:   map[string]*PrivateKey{},
		publicKeys: map[string]*PublicKey{},
		linkedKeys: map[string]*PublicKey{},
	}
}
