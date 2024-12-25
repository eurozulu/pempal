package factories

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/utils"
	"github.com/eurozulu/pempal/validation"
	"os"
	"path/filepath"
	"strings"
)

const RSADefaultKeySize = 2048

type keyFactory struct{}

func (kf keyFactory) Build(t templates.Template) ([]byte, error) {
	vld := validation.ValidatorForTemplate("privatekey")
	if err := vld.Validate(t); err != nil {
		return nil, err
	}
	keyTemplate := t.(*templates.KeyTemplate)

	var prk crypto.PrivateKey
	var err error
	switch x509.PublicKeyAlgorithm(keyTemplate.PublicKeyAlgorithm) {
	case x509.RSA:
		prk, err = kf.createRSAKey(keyTemplate.KeySize)
	case x509.Ed25519:
		prk, err = kf.createEd25519Key()
	default:
		err = fmt.Errorf("%s is an unknown public key algorithm", keyTemplate.PublicKeyAlgorithm)
	}
	if err != nil {
		return nil, err
	}
	return utils.PrivateKeyToPem(prk)
}

func (kf keyFactory) createRSAKey(keysize int) (crypto.PrivateKey, error) {
	if keysize == 0 {
		keysize = RSADefaultKeySize
	}
	return rsa.GenerateKey(rand.Reader, keysize)
}

func (kf keyFactory) createEd25519Key() (crypto.PrivateKey, error) {
	seed := make([]byte, ed25519.SeedSize)
	if _, err := rand.Read(seed); err != nil {
		return nil, err
	}
	return ed25519.NewKeyFromSeed(seed), nil
}

func ensureKeyDirExists(path string, isPrivate bool) error {
	fi, err := os.Stat(path)
	if err == nil {
		if !fi.IsDir() {
			return fmt.Errorf("not a directory")
		}
		if isPrivate && fi.Mode()&0044 != 0000 {
			return fmt.Errorf("private key path %q is not private. ensure only the user can read that directory")
		}
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	perm := os.FileMode(0700)
	if !isPrivate {
		perm = os.FileMode(0755)
	}
	return os.MkdirAll(path, perm)
}

func savePublicKey(path string, pemdata []byte) (string, error) {
	blk, _ := pem.Decode(pemdata)
	puk, err := x509.ParsePKIXPublicKey(blk.Bytes)
	if err != nil {
		return "", err
	}
	if err := ensureKeyDirExists(path, false); err != nil {
		return "", err
	}
	id, err := model.NewKeyIdFromKey(puk)
	if err != nil {
		return "", err
	}
	fp := filepath.Join(path, strings.Join([]string{id.String(), ".pem"}, ""))
	return fp, os.WriteFile(fp, pemdata, 0644)
}

func savePrivateKey(path string, pemdata []byte) (string, error) {
	blk, _ := pem.Decode(pemdata)
	prk, err := x509.ParsePKCS8PrivateKey(blk.Bytes)
	if err != nil {
		return "", err
	}
	if err := ensureKeyDirExists(path, true); err != nil {
		return "", err
	}
	puk, err := utils.PublicKeyFromPrivate(prk)
	if err != nil {
		return "", err
	}
	id, err := model.NewKeyIdFromKey(puk)
	if err != nil {
		return "", err
	}
	fp := filepath.Join(path, strings.Join([]string{id.String(), ".pem"}, ""))
	return fp, os.WriteFile(fp, pemdata, 0600)
}
