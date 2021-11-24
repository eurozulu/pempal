package keytracker

import (
	"encoding/pem"
	"path"
	"pempal/keytools"
	"pempal/pemreader"
	"strings"
)

type collector struct {
	anons map[string]*key
	puks  map[string]*pem.Block
}

// AddBlock adds the given pem block to he collector.
func (c *collector) AddBlock(blk *pem.Block) Key {
	// Check if new block is a private key
	if keytools.PrivateKeyTypes[blk.Type] {
		return c.addPrivateKey(blk)
	}
	// public keys are match, by location to encypted keys
	if keytools.PublicKeyTypes[blk.Type] {
		return c.addPublicKey(blk)
	}
	// unknown block type
	return nil
}

// collects the unmatched private keys, unknown publickey.
func (c *collector) UnknownKeys() []Key {
	ks := make([]Key, len(c.anons))
	var i int
	for _, v := range c.anons {
		ks[i] = v
		i++
	}
	return ks
}

func (c *collector) addPrivateKey(blk *pem.Block) *key {
	// If private key encrypted, attempt to match to known puk.
	k := &key{pemBlock: blk}
	if k.IsEncrypted() {
		// New encrypted key.  attempt to match with known publics by a shared location.
		l := trimLocation(k.Location())
		pk, ok := c.puks[l]
		if !ok {
			// not know, place in waiting room.
			c.anons[l] = k
			return nil
		}
		// matched to a puk by location, move puk intp key
		k.puk = pk
		delete(c.puks, l)
	}
	return k
}

func (c *collector) addPublicKey(blk *pem.Block) *key {
	l := trimLocation(blk.Headers[pemreader.LocationHeaderKey])
	ek, ok := c.anons[l]
	if !ok {
		// move to the waiting room
		c.puks[l] = blk
		return nil
	}
	// move from anon and return
	ek.puk = blk
	delete(c.anons, l)
	return ek
}

func trimLocation(l string) string {
	// strip any index
	lp := strings.TrimRight(l,"0123456789:")
	// strip any extension
	return lp[:len(lp) - len(path.Ext(l))]
}
