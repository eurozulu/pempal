package keytracker

import (
	"encoding/pem"
	"path"
	"pempal/keytools"
	"pempal/pemreader"
	"strings"
)

// keyMatcher collects both private and public key pems and attempts to match them using their sha1 hash or location.
type keyMatcher struct {
	// anons holds the private (Encrypted) keys, which have yet to be matched to any public key
	anons map[string]*key
	// puks holds the public keys which have yet to be matched to any private key
	puks map[string]*pem.Block
	// idPuks holds the public keys containing an 'encryptedHash' header, which have yet to be matched to any private key
	idPuks map[string]*pem.Block
}

// AddBlock adds the given pem block to the keyMatcher.
// If the block results in a paired key, the newly paired key is returned.
// otherwise the block is stored until it is matched and nil is returned.
func (c *keyMatcher) AddBlock(blk *pem.Block) Key {
	// Check if new block is a private key
	if keytools.PrivateKeyTypes[blk.Type] {
		return c.addPrivateKey(blk)
	}
	// check for public key
	if keytools.PublicKeyTypes[blk.Type] {
		return c.addPublicKey(blk)
	}
	// unknown block type
	return nil
}

// collects the unmatched private keys, unknown publickey.
func (c *keyMatcher) UnknownKeys() []Key {
	ks := make([]Key, len(c.anons))
	var i int
	for _, v := range c.anons {
		ks[i] = v
		i++
	}
	return ks
}

func (c *keyMatcher) addPrivateKey(blk *pem.Block) *key {
	k := &key{pemBlock: blk}

	// If private key encrypted, attempt to match to known puk.
	// check key has its public available. (Not encrypted)
	if k.PublicKey() == nil {
		// no PUK on key, attempt to match it with a known PUK
		// look for an id puk
		id := k.String()
		loc := trimLocation(k.Location())
		puk, ok := c.idPuks[id]
		if ok {
			delete(c.idPuks, id)
		} else {
			// no ID match, see if it shares location with a PUK
			puk, ok = c.puks[loc]
			if ok {
				delete(c.puks, loc)
			}
		}
		// if no PUK found, att this to anon
		if puk == nil {
			c.anons[loc] = k
			return nil
		}
		// have a paired match, link and return new key
		k.puk = puk
	}
	return k
}

func (c *keyMatcher) addPublicKey(blk *pem.Block) *key {
	// check if PUK is an ID key (contains the encryptedHash header)
	_, ok := readBlockHeader(encryptedHash, blk)
	if ok {
		return c.addIDPublicKey(blk)
	}
	// see if any anons share the same location
	l := trimLocation(blk.Headers[pemreader.LocationHeaderKey])
	ak, ok := c.anons[l]
	if ok {
		// pair anon with puk and send them on their way
		ak.puk = blk
		delete(c.anons, l)
		return ak
	}
	// unknown PUK, place in the puks pile, to be claimed
	c.puks[l] = blk
	return nil
}

func (c *keyMatcher) addIDPublicKey(blk *pem.Block) *key {
	id, _ := readBlockHeader(encryptedHash, blk)
	// search anons for a matching id
	var anonKey string
	// find key as we dont want to delete from anons in the loop
	for k, v := range c.anons {
		if v.String() != id {
			continue
		}
		anonKey = k
		break
	}
	if anonKey == "" {
		// no matching anon, place on the waiting pile, keyed by ID
		c.idPuks[id] = blk
		return nil
	}
	// found matching anon.  Link the keys and send them on the way
	ak := c.anons[anonKey]
	delete(c.anons, anonKey)
	ak.puk = blk
	return ak
}

func readBlockHeader(key string, b *pem.Block) (string, bool) {
	if b == nil || b.Headers == nil {
		return "", false
	}
	s, ok := b.Headers[key]
	return s, ok
}

func trimLocation(l string) string {
	// strip any index
	lp := strings.TrimRight(l, "0123456789:")
	// strip any extension
	return lp[:len(lp)-len(path.Ext(l))]
}
