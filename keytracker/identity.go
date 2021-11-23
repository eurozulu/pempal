package keytracker

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"pempal/keytools"
	"pempal/pemreader"
)

// Identity represents a public key in the form of a public key or signed resource containing a public key
// Identities may also include a Private key, when once has been matched to the identities PublicKey.
type Identity interface {
	fmt.Stringer
	// PublicKey is the Public key of this identity
	PublicKey() crypto.PublicKey

	// PrivateKey gets the key for this Identity, if it is present and available.
	// When no provate key is linked to the identity, this returns nil
	// The resulting key may or may not be encrypted
	Key() Key

	// Certificates will get all the certificates bound to this ID, filtered by any Keyusage/ExtKeyusage set.
	Certificates(ku x509.KeyUsage, eku ...x509.ExtKeyUsage) []*x509.Certificate
}

type identity struct {
	prk  Key
	blks []*pem.Block
}

func (id identity) String() string {
	puk := id.PublicKey()
	if puk != nil {
		return keytools.PublicKeySha1Hash(puk)
	}
	prk := id.Key()
	if prk == nil {
		return ""
	}
	return prk.String()
}

func (id identity) Key() Key {
	return id.prk
}

// PublicKey gets the single, unique public key this identity
// If the privzte key is present and unencrypted, this is used to generate the public key.
// When not present, the key uses an associcated public key, if set. Otherwise nil is returned
func (id identity) PublicKey() crypto.PublicKey {
	var puk crypto.PublicKey

	// use private key if available
	puk = id.prk.PublicKey()
	if puk != nil {
		return puk
	}

	// First see if we have a public key block
	blks := id.filteredBlocks(keytools.PublicKeyTypes)
	if len(blks) > 0 {
		return blks[0]
	}

	// look for a certificate
	certs := id.Certificates(0, 0)
	for _, c := range certs {
		if c.PublicKey == nil {
			continue
		}
		return c.PublicKey
	}
	// no public key available, possibly unmatched encrypted key
	return nil
}

func (id identity) Locations() []string {
	var locs []string
	for _, blk := range id.blks {
		if len(blk.Headers) == 0 {
			continue
		}
		l := blk.Headers[pemreader.LocationHeaderKey]
		if l == "" {
			continue
		}
		locs = append(locs, l)
	}
	return locs
}

func (id identity) Certificates(ku x509.KeyUsage, eku ...x509.ExtKeyUsage) []*x509.Certificate {
	blks := id.filteredBlocks(keytools.CertificateTypes)
	var certs []*x509.Certificate
	for _, puk := range blks {
		c, err := x509.ParseCertificate(puk.Bytes)
		if err != nil {
			log.Println(err)
			continue
		}
		if ku != 0 && c.KeyUsage&ku != ku {
			continue
		}
		if len(eku) > 0 && !containsExtKeyUsages(c, eku) {
			continue
		}
		certs = append(certs, c)
	}
	return certs
}

func (id identity) filteredBlocks(types map[string]bool) []*pem.Block {
	var blocks []*pem.Block
	for _, blk := range id.blks {
		if !types[blk.Type] {
			continue
		}
		blocks = append(blocks, blk)
	}
	return blocks
}

func containsExtKeyUsages(c *x509.Certificate, ekus []x509.ExtKeyUsage) bool {
	for _, eku := range ekus {
		if !containsExtKeyUsage(c, eku) {
			return false
		}
	}
	return true
}
func containsExtKeyUsage(c *x509.Certificate, eku x509.ExtKeyUsage) bool {
	for _, cku := range c.ExtKeyUsage {
		if cku == eku {
			return true
		}
	}
	return false
}

func NewIdentity(k Key, blks []*pem.Block) Identity {
	return &identity{
		prk:  k,
		blks: blks,
	}
}
