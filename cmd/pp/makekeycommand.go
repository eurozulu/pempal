package main

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"github.com/eurozulu/pempal/encoding"
	"github.com/eurozulu/pempal/templates"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const defaultKeyAlgorithm = x509.RSA
const defaultKeyLengthRSA = 2048
const defaultEncoding = "pem"

var defaultKeyLengthECDSA = elliptic.P521()

type MakeKeyCommand struct {
	Encode             string `flag:"encode,en"`
	PublicKeyAlgorithm string `flag:"keyalgorithm,k"`
	KeyLength          string `flag:"length,l"`
	Encrypt            bool   `flag:"encrypt,e"`
	NoPublicKey        bool   `flag:"nopublickey,npub,npuk,n"`

	// Script flag, when set, prevents any user prompting and assumes confirmation when creating.
	Script   bool   `flag:"script, s"`
	Password string `flag:"password,p"`
}

// MakeKey creates a new key.
// It has two optional commands:
// makekey <filename of new private key> <filename of new public key>
// When only the first is specified, the public key is written to the same name, with a '.pub' extension.
// When neither are specified, output is to stdout in encoding set with 'encode' flag (default pem)
func (kc MakeKeyCommand) MakeKey(args ...string) error {
	if len(args) > 2 {
		return fmt.Errorf("unexpected argument.  Expecting only 2 arguments max. found %d", len(args))
	}
	kt, err := kc.NewKey()
	if err != nil {
		return err
	}

	if len(args) > 0 {
		p, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		kt.FilePath = p
	}

	var pukPath string
	if !kc.NoPublicKey && kt.FilePath != "" {
		if len(args) > 1 {
			pukPath = args[1]
		} else {
			pukPath = strings.Join([]string{kt.FilePath, "pub"}, ".")
		}
	}
	if !kc.confirmKeys(kt, pukPath) {
		return fmt.Errorf("aborted")
	}

	// Write out templates to encoder
	buf := bytes.NewBuffer(nil)
	enc, err := kc.encoder(buf)
	if err != nil {
		return fmt.Errorf("invalid 'encode' flag  %v", err)
	}
	if err := enc.Encode([]templates.Template{kt}); err != nil {
		return err
	}
	if err := writeOutBytes(kt.FilePath, buf.Bytes(), 0600); err != nil {
		return err
	}

	if kc.NoPublicKey {
		return nil
	}
	pkt := templates.NewPublicKeyTemplate(kt.PublicKey())
	buf.Reset()
	if err := enc.Encode([]templates.Template{pkt}); err != nil {
		return err
	}
	return writeOutBytes(pukPath, buf.Bytes(), 0600)
}

// NewKey generates a new Privatekey
func (kc MakeKeyCommand) NewKey() (*templates.PrivateKeyTemplate, error) {
	var key crypto.PrivateKey
	keyAlgo, err := kc.keyAlgorithm()
	if err != nil {
		return nil, err
	}
	switch keyAlgo {
	case x509.RSA:
		key, err = kc.makeRSA()
	case x509.ECDSA:
		key, err = kc.makeECDSA()
	case x509.Ed25519:
		key, err = kc.makeEd25519()
	default:
		err = fmt.Errorf("key algorithm %v is not supported", keyAlgo)
	}
	if err != nil {
		return nil, err
	}
	by, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}

	pwd := kc.Password
	// If password provided, assume key is encrypted
	if pwd != "" {
		kc.Encrypt = true
	}
	// If encrypt request but no password, ask for one (unless scripting, then error)
	if kc.Encrypt && pwd == "" {
		if kc.Script {
			return nil, fmt.Errorf("new key encryption failed as no password provided")
		}
		pwd, err = PromptCreatePassword("Enter a new password for the key.  (Hit enter for unencrypted key)", 0)
		if err != nil {
			return nil, err
		}
	}

	prkT := &templates.PrivateKeyTemplate{
		Passphrase: pwd,
		Encrypted:  kc.Encrypt,
	}
	if err := prkT.UnmarshalBinary(by); err != nil {
		return nil, err
	}
	return prkT, nil
}

func (kc MakeKeyCommand) keyAlgorithm() (x509.PublicKeyAlgorithm, error) {
	pka := defaultKeyAlgorithm
	if kc.PublicKeyAlgorithm != "" {
		var ka templates.PublicKeyAlgorithm
		if err := yaml.Unmarshal([]byte(kc.PublicKeyAlgorithm), &ka); err != nil {
			return 0, fmt.Errorf("%v use one of: %s, %s, %s", err,
				templates.PublicKeyAlgorithms[1],
				templates.PublicKeyAlgorithms[3],
				templates.PublicKeyAlgorithms[4])
		}
		pka = x509.PublicKeyAlgorithm(ka)
	}
	return pka, nil
}

func (kc MakeKeyCommand) encoder(out io.Writer) (encoding.TemplateEncoder, error) {
	ec := defaultEncoding
	if kc.Encode != "" {
		ec = kc.Encode
	}
	return encoding.NewEncoder(ec, out)
}

func (kc MakeKeyCommand) makeEd25519() (crypto.PrivateKey, error) {
	if kc.KeyLength != "" {
		return nil, fmt.Errorf("Key length does not apply for Ed25519 keys")
	}
	_, prk, err := ed25519.GenerateKey(rand.Reader)
	return prk, err
}

func (kc MakeKeyCommand) makeECDSA() (crypto.PrivateKey, error) {
	c, err := kc.keyLengthECDSA()
	if err != nil {
		return nil, err
	}
	return ecdsa.GenerateKey(c, rand.Reader)
}

func (kc MakeKeyCommand) keyLengthECDSA() (elliptic.Curve, error) {
	if kc.KeyLength == "" {
		return defaultKeyLengthECDSA, nil
	}
	switch kc.KeyLength {
	case "521", "P521", "p521":
		return elliptic.P521(), nil
	case "384", "P384", "p384":
		return elliptic.P384(), nil
	case "256", "P256", "p256":
		return elliptic.P256(), nil
	case "224", "P224", "p224":
		return elliptic.P224(), nil
	default:
		return nil, fmt.Errorf("invalid key length for ECDSA, use p521, p384, p256 or p224")
	}
}

func (kc MakeKeyCommand) makeRSA() (crypto.PrivateKey, error) {
	l, err := kc.keyLengthRsa()
	if err != nil {
		return nil, err
	}
	return rsa.GenerateKey(rand.Reader, l)
}

func (kc MakeKeyCommand) keyLengthRsa() (int, error) {
	l := defaultKeyLengthRSA
	if kc.KeyLength != "" {
		i, err := strconv.ParseInt(kc.KeyLength, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid keylength for a rsa key.  Defaults to %d  %v",
				defaultKeyLengthRSA, err)
		}
		l = int(i)
	}
	return l, nil
}

func (kc MakeKeyCommand) confirmKeys(key *templates.PrivateKeyTemplate, pubOut string) bool {
	if kc.Script {
		return true
	}
	enc, err := encoding.NewEncoder("yaml", os.Stdout)
	if err != nil {
		log.Println(err)
		return false
	}
	err = enc.Encode([]templates.Template{key})
	if err != nil {
		log.Println(err)
		return false
	}
	if pubOut != "" {
		fmt.Printf("Public key saved as: %s\n", pubOut)
	}
	return PromptConfirm("\nCreate new key pair?", false)
}

func writeOutBytes(p string, by []byte, perm os.FileMode) error {
	if p == "" {
		_, err := os.Stdout.Write(by)
		return err
	}
	return ioutil.WriteFile(p, by, perm)
}
