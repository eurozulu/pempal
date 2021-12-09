package pemresources

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"pempal/fileformats"
	"pempal/keytools"
	"strings"
)

type PrivateKey struct {
	PemResource
	PrivateKey         string                  `yaml:"private_key"`
	PrivateKeyHash     string                  `yaml:"private_key_hash,omitempty"`
	IsEncrypted        bool                    `yaml:"is_encrypted"`
	PublicKey          string                  `yaml:"public_key"`
	PublicKeyHash      string                  `yaml:"public_key_hash"`
	PublicKeyAlgorithm x509.PublicKeyAlgorithm `yaml:"public_key_algorithm"`
	PublicKeyLength    string                  `yaml:"public_key_length"`
}

func (kt *PrivateKey) LinkPublicKey(puk *PublicKey) error {
	if puk.LinkedId != "" && !strings.HasSuffix(puk.LinkedId, kt.PrivateKeyHash) {
		return fmt.Errorf("linked id %s does not match this key", puk.LinkedId)
	}
	if TrimLocation(puk.Location) != TrimLocation(kt.Location) {
		return fmt.Errorf("location %s does not match this key", puk.Location)
	}
	kt.PublicKey = puk.PublicKey
	kt.PublicKeyHash = puk.PublicKeyHash
	kt.PublicKeyAlgorithm = puk.PublicKeyAlgorithm
	kt.PublicKeyLength = puk.PublicKeyLength
	return nil
}

func (kt PrivateKey) PublicKeyTemplate() (*PublicKey, error) {
	if kt.PublicKey == "" {
		return nil, fmt.Errorf("Public key ot available on this private key.")
	}
	by, err := base64.StdEncoding.DecodeString(kt.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("public key is not valid base64  %w", err)
	}
	pk := &PublicKey{}
	copyHeadersWithoutEncryption(kt.PemHeaders, pk.PemHeaders)
	if err := pk.UnmarshalPem(&pem.Block{
		Type:    strings.Replace(kt.PemType, "PRIVATE", "PUBLIC", -1),
		Headers: pk.PemHeaders,
		Bytes:   by,
	}); err != nil {
		return nil, err
	}
	return pk, nil
}
func (kt PrivateKey) Decrypt(password string) (*PrivateKey, error) {
	if !kt.IsEncrypted {
		return &kt, nil
	}
	blk, err := kt.MarshalPem()
	if err != nil {
		return nil, err
	}
	der, err := x509.DecryptPEMBlock(blk, []byte(password))
	if err != nil {
		return nil, err
	}
	delete(blk.Headers, "DEK-Info")
	blk.Bytes = der

	pk := &PrivateKey{}
	if err = pk.UnmarshalPem(blk); err != nil {
		return nil, err
	}
	return pk, nil
}

func (kt PrivateKey) ResourceId() string {
	return kt.PublicKeyHash
}

func (kt PrivateKey) MarshalPem() (*pem.Block, error) {
	if kt.PrivateKey == "" {
		return nil, fmt.Errorf("no PrivateKey value is set")
	}
	by, err := base64.StdEncoding.DecodeString(kt.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("PrivateKey value is not a valid base64 encoded private key  %w", err)
	}
	blk, err := kt.PemResource.MarshalPem()
	if err != nil {
		return nil, err
	}
	blk.Bytes = by
	return blk, nil

}

func copyHeadersWithoutEncryption(from map[string]string, to map[string]string) {
	for k, v := range from {
		if k == "DEK-Info" {
			continue
		}
		to[k] = v
	}
}

func (kt *PrivateKey) UnmarshalPem(block *pem.Block) error {
	if !fileformats.PemTypesPrivateKey[block.Type] {
		return fmt.Errorf("not a private key pem")
	}
	if err := kt.PemResource.UnmarshalPem(block); err != nil {
		return err
	}

	kt.PrivateKey = base64.StdEncoding.EncodeToString(block.Bytes)
	kt.PrivateKeyHash = keytools.SHA1HashString(block.Bytes)

	kt.IsEncrypted = x509.IsEncryptedPEMBlock(block)
	if !kt.IsEncrypted {
		prk, err := fileformats.ParsePrivateKey(block.Bytes)
		if err != nil {
			return err
		}
		puk := keytools.PublicKeyFromPrivate(prk)
		if puk == nil {
			return fmt.Errorf("failed to read public key from private key")
		}
		by, err := fileformats.MarshalPublicKey(puk)
		if err != nil {
			return err
		}
		kt.PublicKey = base64.StdEncoding.EncodeToString(by.Bytes)
		kt.PublicKeyHash = keytools.SHA1HashString(by.Bytes)
		kt.PublicKeyAlgorithm = keytools.PublicKeyAlgorithm(puk)
		kt.PublicKeyLength = keytools.PublicKeyLength(puk)

	} else {
		kt.PublicKey = ""
		kt.PublicKeyAlgorithm = 0
		kt.PublicKeyLength = ""
		kt.PublicKeyHash = strings.Join([]string{"*", kt.PrivateKeyHash}, "")
	}
	return nil
}
