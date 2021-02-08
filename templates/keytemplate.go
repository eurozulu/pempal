package templates

import (
	"crypto"
	"fmt"
	"strings"

	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"gopkg.in/yaml.v3"
)

const defaultPublicKeyAlgorithm = x509.RSA
const defaultPublicKeySize = "2048"
const defaultPEMCipher = x509.PEMCipherAES256

type (
	// PrivateKeyTemplate represents an existing private key.
	PrivateKeyTemplate struct {
		Password string
		keypem   *pem.Block
	}

	// PublicKeyTemplate represents an existing public key, with no private key counterpart
	PublicKeyTemplate struct {
		key crypto.PublicKey
	}

	// NewKeyTemplate represents a key configuration, prior to generating.
	NewKeyTemplate struct {
		PublicKeyAlgorithm PublicKeyAlgorithmTemplate `yaml:"PublicKeyAlgorithm"`
		PublicKeyLength    string                     `yaml:"PublicKeyLength,omitempty"`
		IsEncrypted        bool                       `yaml:"IsEncrypted"`
		PEMCipher          PEMCipherTemplate          `yaml:"PEMCipher,omitempty"`
		Password           string                     `yaml:"-"`
	}

	yamlPrivateKey struct {
		IsEncrypted          bool                       `yaml:"IsEncrypted"`
		PublicKeyAlgorithm   PublicKeyAlgorithmTemplate `yaml:"PublicKeyAlgorithm,omitempty"`
		PublicKeyLength      string                     `yaml:"PublicKeyLength,omitempty"`
		PublicKeyFingerprint string                     `yaml:"PublicKeyFingerprint,omitempty"`
	}

	yamlPublicKey struct {
		PublicKeyAlgorithm   PublicKeyAlgorithmTemplate `yaml:"PublicKeyAlgorithm"`
		PublicKeyLength      string                     `yaml:"PublicKeyLength,omitempty"`
		PublicKeyFingerprint string                     `yaml:"PublicKeyFingerprint,omitempty"`
		PublicKey            string                     `yaml:"PublicKey,omitempty"`
	}
)

func NewPrivateKeyTemplate(keypem *pem.Block) *PrivateKeyTemplate {
	return &PrivateKeyTemplate{keypem: keypem}
}

// NewNewKeyTemplate creates a New NewKeyTemplate with the default key properties
// - PublicKeyAlgorithm: rsa, PublicKeyLength: 2048, IsEncrypted: true, PEMCipher: PEMCipherAES256
func NewNewKeyTemplate() *NewKeyTemplate {
	return &NewKeyTemplate{
		PublicKeyAlgorithm: PublicKeyAlgorithmTemplate(defaultPublicKeyAlgorithm),
		PublicKeyLength:    defaultPublicKeySize,
		IsEncrypted:        true,
		PEMCipher:          PEMCipherTemplate(defaultPEMCipher),
	}
}

func (pk PrivateKeyTemplate) IsEncrypted() bool {
	return x509.IsEncryptedPEMBlock(pk.keypem)
}

func (pk PrivateKeyTemplate) PrivateKey() (crypto.PrivateKey, error) {
	by := pk.keypem.Bytes
	if pk.IsEncrypted() && pk.Password != "" {
		der, err := x509.DecryptPEMBlock(pk.keypem, []byte(pk.Password))
		if err != nil {
			return nil, err
		}
		by = der
	}
	return x509.ParsePKCS8PrivateKey(by)
}

func (pk PrivateKeyTemplate) PublicKey() (crypto.PublicKey, error) {
	prk, err := pk.PrivateKey()
	if err != nil {
		return nil, err
	}
	return PublicKeyFromPrivate(prk), nil
}

func (pk PrivateKeyTemplate) PublicKeyAlgorithm() PublicKeyAlgorithmTemplate {
	puk, err := pk.PublicKey()
	if err != nil {
		return 0
	}
	return PublicKeyAlgorithmTemplate(PublicKeyAlgorithm(puk))
}

func (pk PrivateKeyTemplate) PublicKeyLength() string {
	puk, err := pk.PublicKey()
	if err != nil {
		return ""
	}
	return PublicKeyLength(puk)
}

func (pk PrivateKeyTemplate) PublicKeyFingerprint() string {
	puk, err := pk.PublicKey()
	if err != nil {
		return ""
	}
	return PublicKeyFingerprint(puk)
}

func (pk PrivateKeyTemplate) PEMBlock() *pem.Block {
	by := make([]byte, len(pk.keypem.Bytes))
	copy(by, pk.keypem.Bytes)
	pb := &pem.Block{
		Type:    pk.keypem.Type,
		Headers: pk.keypem.Headers,
		Bytes:   by,
	}
	return pb
}

func (pk PrivateKeyTemplate) String() string {
	return strings.Join([]string{pk.PublicKeyAlgorithm().String(), pk.PublicKeyLength()}, "\t")
}

func (pk PrivateKeyTemplate) MarshalYAML() (interface{}, error) {
	return &yamlPrivateKey{
		IsEncrypted:          pk.IsEncrypted(),
		PublicKeyAlgorithm:   pk.PublicKeyAlgorithm(),
		PublicKeyLength:      pk.PublicKeyLength(),
		PublicKeyFingerprint: pk.PublicKeyFingerprint(),
	}, nil
}

func (nkt NewKeyTemplate) String() string {
	return strings.Join([]string{nkt.PublicKeyAlgorithm.String(), nkt.PublicKeyLength}, "\t")
}

func (pkt PublicKeyTemplate) String() string {
	return strings.Join([]string{pkt.PublicKeyAlgorithm().String(), pkt.PublicKeyLength()}, "\t")
}
func (pkt PublicKeyTemplate) PublicKey() crypto.PublicKey {
	return pkt.PublicKey()
}
func (pkt PublicKeyTemplate) PublicKeyAlgorithm() PublicKeyAlgorithmTemplate {
	return PublicKeyAlgorithmTemplate(PublicKeyAlgorithm(pkt.key))
}
func (pkt PublicKeyTemplate) PublicKeyLength() string {
	return PublicKeyLength(pkt.key)
}
func (pkt PublicKeyTemplate) PublicKeyFingerprint() string {
	return PublicKeyFingerprint(pkt.key)
}

func (pkt PublicKeyTemplate) MarshalYAML() (interface{}, error) {
	pkb, err := x509.MarshalPKIXPublicKey(pkt.key)
	if err != nil {
		return nil, err
	}
	return &yamlPublicKey{
		PublicKeyAlgorithm:   pkt.PublicKeyAlgorithm(),
		PublicKeyLength:      pkt.PublicKeyLength(),
		PublicKeyFingerprint: pkt.PublicKeyFingerprint(),
		PublicKey:            fmt.Sprintf("%x", pkb),
	}, nil
}

func (pkt *PublicKeyTemplate) UnmarshalYAML(value *yaml.Node) error {
	var yp yamlPublicKey
	if err := value.Decode(&yp); err != nil {
		return err
	}
	var by []byte
	_, err := fmt.Sscanf(yp.PublicKey, "%x", &by)
	if err != nil {
		return err
	}
	k, err := x509.ParsePKIXPublicKey(by)
	if err != nil {
		return err
	}
	pkt.key = k
	return nil
}

func PublicKeyFingerprint(k crypto.PublicKey) string {
	by, err := x509.MarshalPKIXPublicKey(k)
	if err != nil {
		return ""
	}
	h := sha256.New()
	h.Write(by)
	return string(h.Sum(nil))
}
