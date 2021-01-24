package templates

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal"
	"golang.org/x/crypto/ssh"
	"strings"
)

// PrivateKeyTemplate is for private keys encoded in pkcs8 or rsa
type PrivateKeyTemplate struct {
	IsEncrypted          bool               `yaml:"IsEncrypted,omitempty"`
	PublicKeyAlgorithm   PublicKeyAlgorithm `yaml:"PublicKeyAlgorithm,omitempty"`
	PublicKeyFingerprint string             `yaml:"PublicKeyFingerprint,omitempty"`
	FilePath             string             `yaml:"Location,omitempty"`
	Passphrase           string             `yaml:"-"`

	key      crypto.PrivateKey
	pemBlock *pem.Block // Keeps pem block for when key is encrypted
}

func (t *PrivateKeyTemplate) String() string {
	pkf := t.PublicKeyFingerprint
	if pkf == "" {
		pkf = "???"
	}
	pka := t.PublicKeyAlgorithm.String()
	if pka == "" {
		pka = "???"
	}

	e := "encrypted"
	if !t.IsEncrypted {
		e = "unencrypted"
	}
	return strings.Join([]string{TemplateType(t), pka, e, pkf, t.Location()}, "\t")
}

func (t *PrivateKeyTemplate) Decrypt(passphrase string) error {
	if t.pemBlock == nil {
		return fmt.Errorf("template %s has no binary pem data", t.FilePath)
	}
	t.Passphrase = passphrase
	return t.unmarshalPEM(t.pemBlock)
}

func (t PrivateKeyTemplate) Location() string {
	return t.FilePath
}

func (t *PrivateKeyTemplate) UnmarshalBinary(data []byte) error {
	bl, _ := pem.Decode(data)
	if bl != nil {
		return t.unmarshalPEM(bl)
	}

	k, err := x509.ParsePKCS8PrivateKey(data)
	if err != nil {
		if t.IsEncrypted {
			return nil
		}
		// Try to parse as an rsa key
		k, err = x509.ParsePKCS1PrivateKey(data)
		if err != nil {
			return err
		}
	}
	pk := pempal.PublicKeyFromPrivate(k)
	pkssh, err := ssh.NewPublicKey(pk)
	if err != nil {
		return err
	}
	t.PublicKeyFingerprint = ssh.FingerprintSHA256(pkssh)
	t.PublicKeyAlgorithm = PublicKeyAlgorithm(pempal.PublicKeyAlgorithm(pk))
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

func (t *PrivateKeyTemplate) unmarshalPEM(bl *pem.Block) error {
	t.IsEncrypted = x509.IsEncryptedPEMBlock(bl)
	if t.IsEncrypted && t.Passphrase != "" {
		by, err := x509.DecryptPEMBlock(bl, []byte(t.Passphrase))
		if err != nil {
			return err
		}
		bl.Bytes = by
	}
	t.pemBlock = bl
	return t.UnmarshalBinary(bl.Bytes)
}

type SSHPrivateKeyTemplate struct {
	PrivateKeyTemplate
}

func (t *SSHPrivateKeyTemplate) String() string {
	s := strings.Split(t.PrivateKeyTemplate.String(), "\t")
	s[0] = TemplateType(t)
	return strings.Join(s, "\t")
}

func (t *SSHPrivateKeyTemplate) UnmarshalPEM(bl *pem.Block) error {
	t.pemBlock = bl
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
		by, err := x509.MarshalPKCS8PrivateKey(k)
		if err != nil {
			return err
		}
		return t.PrivateKeyTemplate.UnmarshalBinary(by)
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
func (t *SSHPrivateKeyTemplate) Decrypt(passphrase string) error {
	if t.pemBlock == nil {
		return fmt.Errorf("template %s has no binary pem data", t.FilePath)
	}
	t.Passphrase = passphrase
	return t.UnmarshalPEM(t.pemBlock)
}
