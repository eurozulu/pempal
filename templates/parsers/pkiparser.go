package parsers

import (
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"pempal/keytools"
)

const PKI_KeyAlgorithm = X509PublicKeyAlgorithm
const PKI_KeyLength = "PublicKeyLength"
const PKI_PrivateKey = "PrivateKey"
const PKI_PublicKey = X509PublicKey
const PKI_PublicKeyHash = "PublicKeyHash"

var ALLPKINames = []string{
	PKI_KeyAlgorithm,
	PKI_KeyLength,
	PKI_PrivateKey,
	PKI_PublicKey,
	PKI_PublicKeyHash,
}

type PkiParser struct{}

func (cp PkiParser) KnownNames() []string {
	return ALLPKINames
}

func (cp PkiParser) Parse(b *pem.Block) (map[string]interface{}, error) {
	if keytools.PublicKeyTypes[b.Type] {
		return parsePublicKey(b)
	}
	if keytools.PrivateKeyTypes[b.Type] {
		return parsePrivateKey(b)
	}
	return nil, fmt.Errorf("unknown key format")
}

func parsePrivateKey(blk *pem.Block) (map[string]interface{}, error) {
	prk, err := keytools.ParsePrivateKey(blk)
	if err != nil {
		return nil, err
	}

	puk := keytools.PublicKeyFromPrivate(prk)
	m, err := mapPublicKey(puk)
	if err != nil {
		return nil, err
	}
	m[PEM_TYPE] = blk.Type
	m[PKI_PrivateKey] = "true"
	return m, nil
}

func parsePublicKey(blk *pem.Block) (map[string]interface{}, error) {
	puk, err := keytools.ParsePublicKeyPem(blk)
	if err != nil {
		return nil, err
	}
	m, err := mapPublicKey(puk)
	if err != nil {
		return nil, err
	}
	m[PEM_TYPE] = blk.Type
	return m, nil
}

func mapPublicKey(puk crypto.PublicKey) (map[string]interface{}, error) {
	ka := keytools.PublicKeyAlgorithm(puk)
	if ka == x509.UnknownPublicKeyAlgorithm {
		return nil, fmt.Errorf("unknown public key algorithm")
	}
	pem, err := keytools.MarshalPublicKey(puk)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{
		PKI_KeyAlgorithm:  ka.String(),
		PKI_KeyLength:     keytools.PublicKeyLength(puk),
		PKI_PublicKey:     base64.StdEncoding.EncodeToString(pem.Bytes),
		PKI_PublicKeyHash: keytools.PublicKeySha1Hash(puk),
	}
	return m, nil
}
