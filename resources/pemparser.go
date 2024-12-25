package resources

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/utils"

	cryptobyte "golang.org/x/crypto/cryptobyte"
	cryptobyte_asn1 "golang.org/x/crypto/cryptobyte/asn1"
)

// PemParser parses a byte block into one or more Pem Blocks
type PemParser interface {
	CanParse(data []byte) bool
	FormatAsPem(data []byte) ([]*pem.Block, error)
}

// PemResourceParser is an implementation of PemParser for PEM files.
// i.e. it simply decodes the given bytes into blocks.
type PemResourceParser struct{}

func (rf PemResourceParser) CanParse(data []byte) bool {
	return bytes.Contains(data, []byte("-----BEGIN ")) &&
		bytes.Contains(data, []byte("-----END "))
}

func (rf PemResourceParser) FormatAsPem(data []byte) ([]*pem.Block, error) {
	var blocks []*pem.Block
	for {
		b, rest := pem.Decode(data)
		if b == nil {
			break
		}
		blocks = append(blocks, b)
		data = rest
	}
	if len(blocks) == 0 {
		return nil, fmt.Errorf("no pem resources found")
	}
	return blocks, nil
}

// DerResourceParser is an implementation of PemParser to read a single DER encoded byteblock
type DerResourceParser struct{}

func (rf DerResourceParser) FormatAsPem(data []byte) ([]*pem.Block, error) {
	var rt model.ResourceType
	if _, err := x509.ParseCertificate(data); err == nil {
		rt = model.Certificate
	} else if _, err := x509.ParseCertificateRequest(data); err == nil {
		rt = model.CertificateRequest
	} else if _, err := x509.ParsePKCS8PrivateKey(data); err == nil {
		rt = model.PrivateKey
	} else if _, err := x509.ParsePKIXPublicKey(data); err == nil {
		rt = model.PublicKey
	} else if _, err := x509.ParseRevocationList(data); err == nil {
		rt = model.RevokationList
	}
	if rt == model.UnknownResourceType {
		return nil, fmt.Errorf("unknown DER format")
	}

	blk, _ := pem.Decode(utils.DERToPem(data, rt))
	return []*pem.Block{blk}, nil
}

func (rf DerResourceParser) CanParse(data []byte) bool {
	input := cryptobyte.String(data)
	if !input.ReadASN1Element(&input, cryptobyte_asn1.SEQUENCE) {
		return false
	}
	if !input.ReadASN1(&input, cryptobyte_asn1.SEQUENCE) {
		return false
	}
	return true
}
