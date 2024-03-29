package identity

import (
	"crypto"
	"crypto/md5"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
)

type Identity string

func (p Identity) String() string {
	b, _ := pem.Decode([]byte(p))
	if b == nil {
		return ""
	}
	hash := md5.Sum(b.Bytes)
	return hex.EncodeToString(hash[:])
}

func (p Identity) PublicKey() (crypto.PublicKey, error) {
	if p == "" {
		return nil, fmt.Errorf("identity is empty")
	}
	blk, _ := pem.Decode([]byte(p))
	if blk == nil {
		return nil, fmt.Errorf("Failed to unmarshal identifier '%s'", p)
	}
	return x509.ParsePKIXPublicKey(blk.Bytes)
}

func IsIdentity(data string) bool {
	return Identity(data).String() != ""
}
