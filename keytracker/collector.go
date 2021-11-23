package keytracker

import (
	"crypto"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"pempal/keytools"
	"pempal/pemreader"
)

var certsAndKeys = keytools.CombineMaps(keytools.CertificateTypes, keytools.PrivateKeyTypes, keytools.PublicKeyTypes, keytools.CSRTypes)

// collector collects the pem blocks from the routines scanning the various directories, into a single collection.
// It maintains a map of Identites, PRivate keys bound to zero or more publickey bearing pems (puks, certs, csrs)
// Encrypted keys, where the PUK is not available without decrypting it, are matched to a public key using any other pem containing a location header
// pointing to the private key.
type collector struct {
	ids       map[string]*identity
	anons     map[string]*pem.Block
	unclaimed map[string][]*pem.Block
}

// Identity gets all the identites collected
func (c collector) Identities() []Identity {
	var ids []Identity
	for _, id := range c.ids {
		ids = append(ids, id)
	}
	return ids
}

// AddBlock adds the given pem block to he collector.
func (c *collector) AddBlock(blk *pem.Block) error {
	// Check if new block is a private key
	if keytools.PrivateKeyTypes[blk.Type] {
		return c.addPrivateKey(blk)
	}
	// not a private key, attempt to match to existing keys with its public key
	puk, err := c.publicKeyFromResource(blk)
	if err != nil {
		return err
	}
	hash, err := publicKeyHash(puk)
	if err != nil {
		return err
	}

	// Check if this is a known key
	id := c.ids[hash]
	if id != nil {
		// a known puk, add new block to its collection
		id.blks = append(id.blks, blk)
		return nil
	}

	// check if a known location in anons
	if len(blk.Headers) > 0 {
		k, ok := c.anons[blk.Headers[pemreader.LocationHeaderKey]]
		if ok {
			// matched location of anon key, to this public key, create new ID using both
			pukBlk, err := keytools.MarshalPublicKey(puk)
			if err != nil {
				return err
			}
			c.ids[hash] = &identity{
				prk: &key{
					pemBlock: k,
					puk:      pukBlk,
				},
				blks: []*pem.Block{blk},
			}
			return nil
		}
	}
	// Unmatched resource that's not a private key.  Add to unclaimed
	c.unclaimed[hash] = append(c.unclaimed[hash], blk)
	return nil
}

func (c *collector) addPrivateKey(blk *pem.Block) error {
	k := &key{pemBlock: blk}
	puk := k.PublicKey()
	if puk == nil {
		// no public key == encrypted key.  Map under anon by its location
		if len(blk.Headers) == 0 || blk.Headers[pemreader.LocationHeaderKey] == "" {
			return fmt.Errorf("key %s has no location or public key", k)
		}
		c.anons[blk.Headers[pemreader.LocationHeaderKey]] = blk
		return nil
	}
	hash, err := publicKeyHash(puk)
	if err != nil {
		return err
	}
	if _, ok := c.ids[hash]; ok {
		return fmt.Errorf("duplicate private key found at %s", blk.Headers[pemreader.LocationHeaderKey])
	}
	c.ids[hash] = &identity{
		prk:  k,
		blks: c.unclaimed[hash],
	}
	// remove unclaim resources
	delete(c.unclaimed, hash)
	return nil
}

func (c collector) publicKeyFromResource(blk *pem.Block) (crypto.PublicKey, error) {
	if keytools.PublicKeyTypes[blk.Type] {
		puk, err := keytools.ParsePublicKey(blk)
		if err != nil {
			return "", err
		}
		return puk, nil
	}
	if keytools.CertificateTypes[blk.Type] {
		c, err := x509.ParseCertificate(blk.Bytes)
		if err != nil {
			return nil, err
		}
		return c.PublicKey, nil
	}
	if keytools.CSRTypes[blk.Type] {
		c, err := x509.ParseCertificateRequest(blk.Bytes)
		if err != nil {
			return nil, err
		}
		return c.PublicKey, nil
	}
	return nil, fmt.Errorf("%s is not a known identity pem type", blk.Type)
}

func publicKeyHash(puk crypto.PublicKey) (string, error) {
	if puk == nil {
		return "", nil
	}
	by, err := keytools.MarshalPublicKey(puk)
	if err != nil {
		return "", err
	}
	return stringHash(by.Bytes), nil
}

func stringHash(by []byte) string {
	hash := sha1.New()
	_, _ = hash.Write(by)
	return hex.EncodeToString(hash.Sum(nil))
}

func newCollector() *collector {
	return &collector{
		ids:       map[string]*identity{},
		anons:     map[string]*pem.Block{},
		unclaimed: map[string][]*pem.Block{},
	}
}
