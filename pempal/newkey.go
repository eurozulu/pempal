package pempal

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/pempal/templates"
	"strconv"
	"strings"
)

func FindKeys(q, password string) ([]*QueryResult, error) {
	pq := PEMQuery{
		Query:         []string{q},
		CaseSensitive: false,
		Types:         []string{"PRIVATE KEY"},
		Password:      password,
	}
	return pq.QueryPaths(CertPath(), true)
}

func NewKey(kt *templates.NewKeyTemplate) (*pem.Block, error) {
	bl, err := generateKey(kt)
	if err != nil {
		return nil, err
	}
	if kt.Password != "" {
		kt.IsEncrypted = true
	}
	if !kt.IsEncrypted {
		return bl, err
	}

	return EncryptPemBlock(bl, kt.Password, x509.PEMCipher(kt.PEMCipher))
}

func generateKey(t *templates.NewKeyTemplate) (*pem.Block, error) {
	switch x509.PublicKeyAlgorithm(t.PublicKeyAlgorithm) {
	case x509.RSA:
		kl, err := strconv.Atoi(t.PublicKeyLength)
		if err != nil {
			return nil, fmt.Errorf("'%s' is an invalid keylength for an rsa key. must be integer")
		}
		prk, err := rsa.GenerateKey(rand.Reader, kl)
		if err != nil {
			return nil, fmt.Errorf("failed to generate new rsa key  %v", err)
		}
		return MakePrivateKeyPemBlock(prk)

	case x509.ECDSA:
		prk, err := ecdsa.GenerateKey(curveFromLength(t.PublicKeyLength), rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to generate new ECDSA key  %v", err)
		}
		return MakePrivateKeyPemBlock(prk)

	case x509.Ed25519:
		_, prk, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to generate new Ed25519 key  %v", err)
		}
		return MakePrivateKeyPemBlock(prk)
	default:
		return nil, fmt.Errorf("unsupported key type. must be 'RSA', 'ECDSA' or 'ECDSA'")
	}
}

func curveFromLength(l string) elliptic.Curve {
	if strings.Contains(l, "224") {
		return elliptic.P224()
	}
	if strings.Contains(l, "256") {
		return elliptic.P256()
	}
	if strings.Contains(l, "348") {
		return elliptic.P384()
	}
	if strings.Contains(l, "521") {
		return elliptic.P521()
	}
	return elliptic.P384()
}

// MakePrivateKeyPemBlock wraps the given private key into a PEM block
func MakePrivateKeyPemBlock(prk crypto.PrivateKey) (*pem.Block, error) {
	by, err := x509.MarshalPKCS8PrivateKey(prk)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal new rsa key  %v", err)
	}
	return &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: by,
	}, nil
}

// EncryptPemBlock encrypts the given pem block using the given password and cipher
func EncryptPemBlock(bl *pem.Block, password string, cipher x509.PEMCipher) (*pem.Block, error) {
	if password == "" {
		pw, err := PromptCreatePassword("Enter password to encrypt new key:", 0)
		if err != nil {
			return nil, err
		}
		password = pw
	}
	return x509.EncryptPEMBlock(rand.Reader, bl.Type, bl.Bytes, []byte(password), cipher)
}
