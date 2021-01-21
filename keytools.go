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
	"golang.org/x/crypto/ssh"
	"strconv"
)

func ComparePublicKeys(pk1 crypto.PublicKey, pk2 crypto.PublicKey) bool {
	switch v := pk1.(type) {
	case *rsa.PublicKey:
		return v.Equal(pk2)
	case *ecdsa.PublicKey:
		return v.Equal(pk2)
	case *ed25519.PublicKey:
		return v.Equal(pk2)
	default:
		return false
	}
}

func PublicKeyAlgorithm(pk crypto.PublicKey) x509.PublicKeyAlgorithm {
	switch pk.(type) {
	case *rsa.PublicKey:
		return x509.RSA
	case *ecdsa.PublicKey:
		return x509.ECDSA
	case *ed25519.PublicKey:
		return x509.Ed25519
	default:
		return x509.UnknownPublicKeyAlgorithm
	}
}

func PublicKeyLength(pk crypto.PublicKey) string {
	switch v := pk.(type) {
	case *rsa.PublicKey:
		return strconv.Itoa(v.N.BitLen())
	case *ecdsa.PublicKey:
		return strconv.Itoa(v.Curve.Params().BitSize)
	case *ed25519.PublicKey:
		return ""
	default:
		return ""
	}
}

func PublicKeyFromPrivate(pk crypto.PrivateKey) crypto.PublicKey {
	switch v := pk.(type) {
	case *rsa.PrivateKey:
		return v.Public()
	case *ecdsa.PrivateKey:
		return v.Public()
	case *ed25519.PrivateKey:
		return v.Public()
	default:
		return nil
	}
}

func EncryptPEMKey(by []byte, passphrase string, pemCipher x509.PEMCipher) ([]byte, error) {
	bl, _ := pem.Decode(by)
	ebl, err := x509.EncryptPEMBlock(rand.Reader, bl.Type, bl.Bytes, []byte(passphrase), pemCipher)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(ebl), nil
}

func MarshalPublicKeyToSSH(pk crypto.PublicKey, comment string) ([]byte, error) {
	spk, err := ssh.NewPublicKey(pk)
	if err != nil {
		return nil, err
	}
	by := ssh.MarshalAuthorizedKey(spk)
	if comment != "" {
		by = append(by, ' ')
		by = append(by, []byte(comment)...)
	}
	return by, nil
}

// ParseSSHPublicKey parses the file bytes from a SSH public key into a key and any comments.
func ParseSSHPublicKey(by []byte) (string, crypto.PublicKey, error) {
	spk, cmt, _, _, err := ssh.ParseAuthorizedKey(by)
	if err != nil {
		return "", nil, err
	}
	pk := (spk.(ssh.CryptoPublicKey)).CryptoPublicKey()
	return cmt, pk, nil
}

func GenerateKey(ka x509.PublicKeyAlgorithm, keyLength int) (crypto.PrivateKey, error) {
	switch ka {
	case x509.RSA:
		return rsa.GenerateKey(rand.Reader, keyLength)
	case x509.ECDSA:
		return ecdsa.GenerateKey(curveFromLength(keyLength), rand.Reader)
	case x509.Ed25519:
		_, pk, err := ed25519.GenerateKey(rand.Reader)
		return pk, err
	default:
		return nil, fmt.Errorf("unsupported key type. must be 'RSA', 'ECDSA' or 'ECDSA'")
	}
}

func curveFromLength(l int) elliptic.Curve {
	switch l {
	case 224:
		return elliptic.P224()
	case 256:
		return elliptic.P256()
	case 348:
		return elliptic.P384()
	case 521:
		return elliptic.P521()
	default:
		return elliptic.P384()
	}
}
