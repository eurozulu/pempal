package templates

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal"
	"golang.org/x/crypto/ssh"
)

// PrivateKeyTemplate is for private keys encoded in pkcs8 or rsa
type PrivateKeyTemplate struct {
	IsEncrypted          bool   `yaml:"IsEncrypted,omitempty"`
	PublicKeyFingerprint string `yaml:"PublicKeyFingerprint,omitempty"`
	FilePath             string `yaml:"-"`
	Passphrase           string `yaml:"-"`

	key    crypto.PrivateKey
	pemKey *pem.Block
}

func (t PrivateKeyTemplate) Location() string {
	return t.FilePath
}

func (t *PrivateKeyTemplate) String() string {
	e := "encrypted"
	if !t.IsEncrypted {
		e = "unencrypted"
	}
	pkp := t.PublicKeyFingerprint
	if pkp == "" {
		pkp = "???"
	}
	return fmt.Sprintf("%s\t%s\t", pkp, e)
}

func (t *PrivateKeyTemplate) UnmarshalBinary(data []byte) error {
	k, err := x509.ParsePKCS8PrivateKey(data)
	if err != nil {
		if t.IsEncrypted {
			return nil
		}
		return err
	}
	pk, err := ssh.NewPublicKey(pempal.PublicKeyFromPrivate(k))
	if err != nil {
		return err
	}
	t.PublicKeyFingerprint = ssh.FingerprintSHA256(pk)
	t.key = k
	return nil
}

func (t *PrivateKeyTemplate) MarshalBinary() (data []byte, err error) {
	if t.key == nil {
		if t.IsEncrypted {
			return nil, fmt.Errorf("template %s is encrypted and cannot be marshalled without passphrase", t.FilePath)
		}
		return nil, fmt.Errorf("template %s has no binary key data", t.FilePath)
	}
	return x509.MarshalPKCS8PrivateKey(t.key)
}

func (t *PrivateKeyTemplate) UnmarshalPEM(bl *pem.Block) error {
	t.pemKey = bl
	t.IsEncrypted = x509.IsEncryptedPEMBlock(bl)
	if t.IsEncrypted && t.Passphrase != "" {
		by, err := x509.DecryptPEMBlock(bl, []byte(t.Passphrase))
		if err != nil {
			return err
		}
		bl.Bytes = by
	}
	return t.UnmarshalBinary(bl.Bytes)
}

func (t *PrivateKeyTemplate) MarshalPEM() (*pem.Block, error) {
	if t.key == nil {
		if t.pemKey == nil {
			return nil, fmt.Errorf("template %s has no pem key data", t.FilePath)
		}
		return t.pemKey, nil
	}

	by, err := t.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if t.IsEncrypted {
		return x509.EncryptPEMBlock(rand.Reader, "PRIVATE KEY", by, []byte(t.Passphrase), x509.PEMCipherAES256)
	} else {
		return &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: by,
		}, nil
	}
}

type SSHPrivateKeyTemplate struct {
	PrivateKeyTemplate
}

func (t SSHPrivateKeyTemplate) UnmarshalPEM(bl *pem.Block) error {
	by := pem.EncodeToMemory(bl)
	k, err := ssh.ParseRawPrivateKey(by)
	if err != nil {
		_, t.IsEncrypted = err.(*ssh.PassphraseMissingError)
		if !t.IsEncrypted {
			return err
		}
	} else {
		t.IsEncrypted = false
	}

	if t.IsEncrypted && t.Passphrase != "" {
		k, err = ssh.ParseRawPrivateKeyWithPassphrase(by, []byte(t.Passphrase))
		if err != nil {
			return err
		}
	}

	if k != nil {
		pk, err := ssh.NewPublicKey(pempal.PublicKeyFromPrivate(k))
		if err != nil {
			return err
		}
		t.PublicKeyFingerprint = ssh.FingerprintSHA256(pk)
	} else {
		t.PublicKeyFingerprint = ""
	}
	return nil
}

func (t SSHPrivateKeyTemplate) MarshalPEM() (*pem.Block, error) {
	by, err := t.MarshalBinary()
	if err != nil {
		return nil, err
	}

	if t.IsEncrypted {
		return x509.EncryptPEMBlock(rand.Reader, "OPENSSH PRIVATE KEY", by, []byte(t.Passphrase), x509.PEMCipherAES256)
	} else {
		return &pem.Block{
			Type:  "OPENSSH PRIVATE KEY",
			Bytes: by,
		}, nil
	}
}

func (t SSHPrivateKeyTemplate) MarshalBinary() (data []byte, err error) {
	if t.key == nil {
		return nil, fmt.Errorf("template %s has no binary key data", t.FilePath)
	}
	return x509.MarshalPKCS8PrivateKey(t.key)
}

func (t SSHPrivateKeyTemplate) UnmarshalBinary(data []byte) error {
	k, err := x509.ParsePKCS8PrivateKey(data)
	if err != nil {
		if t.IsEncrypted {
			return nil
		}
		return err
	}
	pk, err := ssh.NewPublicKey(pempal.PublicKeyFromPrivate(k))
	if err != nil {
		return err
	}
	t.PublicKeyFingerprint = ssh.FingerprintSHA256(pk)
	t.key = k
	return nil
}
