package resources

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/utils"
	"log"
)

type Keys interface {
	PrivateKeyFromID(id model.KeyId) (crypto.PrivateKey, error)
	PublicKeyFromID(id model.KeyId) (crypto.PublicKey, error)
	PrivateKeyIDs() ([]model.KeyId, error)
}

type keys struct {
	keypath []string
}

func (k keys) PrivateKeyFromID(id model.KeyId) (crypto.PrivateKey, error) {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	scan := NewPemScan(model.PrivateKey)
	for keyPems := range scan.ScanPath(ctx, k.keypath...) {
		prk, err := findPrivateKeyByID(id, keyPems.Content)
		if err != nil {
			log.Println(err)
		}
		if prk != nil {
			return prk, nil
		}
	}
	return nil, fmt.Errorf("key %s not found", id)

}

// PublicKeyFromID will scan the keypath, in order of paths, to locate the first instance of a public key with the given ID.
// It searches Private and Public Keys.  If a Private key is found, its public key of that private key is returned.
// Otherwise, the public key is returned.
func (k keys) PublicKeyFromID(id model.KeyId) (crypto.PublicKey, error) {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	scan := NewPemScan(model.PublicKey, model.PrivateKey)
	for keyPems := range scan.ScanPath(ctx, k.keypath...) {
		puk, err := findPublicKeyByID(id, keyPems.Content)
		if err != nil {
			return nil, err
		}
		if puk != nil {
			return puk, nil
		}
	}
	return nil, fmt.Errorf("key %s not found", id)
}

func (k keys) PrivateKeyIDs() ([]model.KeyId, error) {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	scan := NewPemScan(model.PrivateKey)
	var found []model.KeyId
	for keyPems := range scan.ScanPath(ctx, k.keypath...) {
		ids, err := privateKeyIDs(keyPems.Content)
		if err != nil {
			return nil, err
		}
		if len(ids) > 0 {
			found = append(found, ids...)
		}
	}
	return found, nil
}

func findPublicKeyByID(id model.KeyId, pems []*pem.Block) (crypto.PublicKey, error) {
	for _, pm := range pems {
		puk, err := utils.PublicKeyFromPem(pm)
		if err != nil {
			return nil, err
		}
		pid, err := model.NewKeyIdFromKey(puk)
		if err != nil {
			return nil, err
		}
		if id.Equals(pid) {
			return puk, nil
		}
	}
	return nil, nil
}

func findPrivateKeyByID(id model.KeyId, pems []*pem.Block) (crypto.PrivateKey, error) {
	for _, pm := range pems {
		puk, err := utils.PublicKeyFromPem(pm)
		if err != nil {
			return nil, err
		}
		pid, err := model.NewKeyIdFromKey(puk)
		if err != nil {
			return nil, err
		}
		if id.Equals(pid) {
			return x509.ParsePKCS8PrivateKey(pm.Bytes)
		}
	}
	return nil, nil
}

func privateKeyIDs(pems []*pem.Block) ([]model.KeyId, error) {
	var found []model.KeyId
	for _, pm := range pems {
		puk, err := utils.PublicKeyFromPem(pm)
		if err != nil {
			return nil, err
		}
		pid, err := model.NewKeyIdFromKey(puk)
		if err != nil {
			return nil, err
		}
		found = append(found, pid)
	}
	return found, nil
}

func NewKeys(keypath []string) Keys {
	return &keys{keypath: keypath}
}
