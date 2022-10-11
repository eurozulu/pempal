package pemresources

import (
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"pempal/fileformats"
	"pempal/keytools"
	"strings"
)

const LinkedKeyHeaderKey = "linkedkey"

type PublicKey struct {
	PemResource
	PublicKeyAlgorithm x509.PublicKeyAlgorithm
	PublicKeyLength    string
	PublicKeyHash      string
	PublicKey          string
	LinkedId           string
}

func (kt *PublicKey) MarshalPem() (*pem.Block, error) {
	blk, err := kt.PemResource.MarshalPem()
	if err != nil {
		return nil, err
	}
	if kt.PublicKey == "" {
		return nil, fmt.Errorf("Public key has no value")
	}
	by, err := base64.StdEncoding.DecodeString(kt.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public string as base64  %s", err)
	}
	blk.Bytes = by
	if !fileformats.PemTypesPublicKey[blk.Type] {
		blk.Type = fileformats.PEM_PUBLIC_KEY
	}
	if kt.LinkedId != "" {
		if blk.Headers == nil {
			blk.Headers = map[string]string{}
		}
		blk.Headers[LinkedKeyHeaderKey] = strings.TrimSpace(kt.LinkedId)
	}
	return blk, nil
}

func (kt *PublicKey) UnmarshalPem(block *pem.Block) error {
	if err := kt.PemResource.UnmarshalPem(block); err != nil {
		return err
	}
	puk := PublicKeyFromBlock(block)
	if puk == nil {
		return fmt.Errorf("cpould not find a public key in the %s block", block.Type)
	}
	blk, err := fileformats.MarshalPublicKey(puk)
	if err != nil {
		return err
	}
	if len(block.Headers) > 0 {
		kt.LinkedId = block.Headers[LinkedKeyHeaderKey]
	}

	kt.PublicKey = base64.StdEncoding.EncodeToString(blk.Bytes)
	kt.PublicKeyHash = keytools.SHA1HashString(blk.Bytes)
	kt.PublicKeyAlgorithm = keytools.PublicKeyAlgorithm(puk)
	kt.PublicKeyLength = keytools.PublicKeyLength(puk)
	return nil
}

// PublicKeyFromBlock extracts the Public key from the given resource.
// given resource may be a Private key (unencrypted), public key, Certificate ro Certificate request.
func PublicKeyFromBlock(blk *pem.Block) crypto.PublicKey {
	if fileformats.PemTypesPrivateKey[blk.Type] {
		prk, err := fileformats.ParsePrivateKey(blk.Bytes)
		if err != nil {
			return nil
		}
		return keytools.PublicKeyFromPrivate(prk)
	}
	if fileformats.PemTypesPublicKey[blk.Type] {
		puk, err := fileformats.ParsePublicKey(blk.Bytes)
		if err != nil {
			return nil
		}
		return puk
	}
	if fileformats.PemTypesCertificate[blk.Type] {
		c, err := x509.ParseCertificate(blk.Bytes)
		if err != nil {
			return nil
		}
		return c.PublicKey
	}
	if fileformats.PemTypesCertificateRequest[blk.Type] {
		c, err := x509.ParseCertificateRequest(blk.Bytes)
		if err != nil {
			return nil
		}
		return c.PublicKey
	}
	return nil
}
