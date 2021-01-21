package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"github.com/eurozulu/pempal"
	"github.com/eurozulu/pempal/encoding"
	"github.com/eurozulu/pempal/templates"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const defaultKeyAlgorithm = x509.RSA
const defaultKeyLengthRSA = 2048
const defaultEncoding = "pem"

var defaultKeyLengthECDSA = elliptic.P521()

type MakeKeyCommand struct {
	OutPath            string `flag:"out,o"`
	Encode             string `flag:"encode,en"`
	PublicKeyAlgorithm string `flag:"keyalgorithm,k"`
	KeyLength          string `flag:"length,l"`
	Encrypt            bool   `flag:"encrypt,e"`

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
	prOut := os.Stdout
	puOut := os.Stdout

	hasFileNames := len(args) > 0

	// private key path given
	if hasFileNames {
		f, err := os.OpenFile(args[0], os.O_APPEND|os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		defer func(fl *os.File) {
			if err := fl.Close(); err != nil {
				log.Println(err)
			}
		}(f)
		// If no public key path given, make one up
		if len(args) == 1 {
			args = append(args, strings.Join([]string{args[0], "pub"}, "."))
		}
	}
	// public key given (or made up)
	if len(args) > 1 {
		f, err := os.OpenFile(args[1], os.O_APPEND|os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		defer func(fl *os.File) {
			if err := fl.Close(); err != nil {
				log.Println(err)
			}
		}(f)
	}
	if len(args) > 2 {
		return fmt.Errorf("unexpected argument.  Expecting only 2 arguments max. found %d", len(args))
	}

	// Make the new key
	var key crypto.PrivateKey
	keyAlgo, err := kc.keyAlgorithm()
	if err != nil {
		return err
	}
	switch keyAlgo {
	case x509.RSA:
		key, err = kc.makeRSA()
	case x509.ECDSA:
		key, err = kc.makeECDSA()
	case x509.Ed25519:
		key, err = kc.makeEd25519()
	default:
		return fmt.Errorf("key algorithm %v is not supported", keyAlgo)
	}
	if err != nil {
		return err
	}

	// Template new keys
	var prName string
	var puName string
	if hasFileNames {
		prName = args[0]
		puName = args[1]
	}

	pwd := kc.Password
	// If password provided, assume key is encrypted
	if pwd != "" {
		kc.Encrypt = true
	}
	// If encrypt request but no password, ask for one (unless scripting, then error)
	if kc.Encrypt && pwd == "" {
		if kc.Script {
			return fmt.Errorf("new key encryption failed as no password provided")
		}
		pwd, err = PromptCreatePassword("Enter a new password for the key.  (Hit enter for unencrypted key)", 0)
		if err != nil {
			return err
		}
	}

	prT, err := privateKeyTemplate(prName, pwd, key)
	if err != nil {
		return err
	}
	puk := pempal.PublicKeyFromPrivate(key)
	puT, err := publicKeyTemplate(puName, puk)
	if err != nil {
		return err
	}

	keys := []templates.Template{prT, puT}
	if !kc.confirmKeys(prT) {
		return fmt.Errorf("aborted")
	}

	// Write out templates to encoder
	prEnc, err := kc.encoder(prOut)
	if err != nil {
		return fmt.Errorf("invalid 'encode' flag  %v", err)
	}
	if err := prEnc.Encode([]templates.Template{prT}); err != nil {
		return err
	}
	puEnc, _ := kc.encoder(puOut) // ignore error as would have triggered on private key
	return puEnc.Encode(keys)
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

func privateKeyTemplate(p, password string, prk crypto.PrivateKey) (*templates.PrivateKeyTemplate, error) {
	by, err := x509.MarshalPKCS8PrivateKey(prk)
	if err != nil {
		return nil, err
	}
	prkT := &templates.PrivateKeyTemplate{
		FilePath:    p,
		Passphrase:  password,
		IsEncrypted: password != "",
	}
	if err := prkT.UnmarshalBinary(by); err != nil {
		return nil, err
	}
	return prkT, nil
}

func publicKeyTemplate(p string, puk crypto.PublicKey) (*templates.PublicKeyTemplate, error) {
	by, err := x509.MarshalPKIXPublicKey(puk)
	if err != nil {
		return nil, err
	}
	prkT := &templates.PublicKeyTemplate{FilePath: p}
	if err := prkT.UnmarshalBinary(by); err != nil {
		return nil, err
	}
	return prkT, nil
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

func (kc MakeKeyCommand) confirmKeys(key templates.Template) bool {
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
	return PromptConfirm("Create new these keys?", false)
}
