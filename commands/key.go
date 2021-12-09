package commands

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"pempal/fileformats"
	"pempal/keycache"
	"pempal/keytools"
	"pempal/pemresources"
	"pempal/templates"
	"strconv"
)

const minPasswordLength = 6
const passwordCipher = x509.PEMCipherAES256

type KeyCommand struct {
	KeyAlgorithm string
	KeyLength    int

	Encrypt  bool
	Password string

	Recursive bool

	LinkPublic bool

	keyCache *keycache.KeyCache
}

func (cmd *KeyCommand) Description() string {
	return "outputs an existing key pair or generates a new one\n"
}

func (cmd *KeyCommand) Flags(fs *flag.FlagSet) {
	fs.StringVar(&cmd.KeyAlgorithm, "keyalgorithm", "rsa", "sets the public key algorithtm")
	fs.StringVar(&cmd.KeyAlgorithm, "ka", "rsa", "same as keyalgorithm")
	fs.IntVar(&cmd.KeyLength, "keylength", 2048, "sets the public key length")
	fs.IntVar(&cmd.KeyLength, "kl", 2048, "same as keylength")
	fs.BoolVar(&cmd.LinkPublic, "link", false, "inserts a custom header in the pem of the newly generated public key to identify its encrypted private key counterpart")

	fs.BoolVar(&cmd.Recursive, "r", false, "search sub directories")
	fs.BoolVar(&cmd.Encrypt, "encrypt", true, "When true, output keys are encrypted using password or prompt user")
	fs.StringVar(&cmd.Password, "password", "", "password to use on decrypting existing passwords to to use when encrypting output keys")
}

func (cmd *KeyCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	pk := &pemresources.PrivateKey{}
	tb := templates.NewTemplateBuilder(pk)
	if err := tb.Add(args...); err != nil {
		return err
	}
	if _, err := tb.Build(); err != nil {
		return err
	}

	// IF no key set, generate a new one
	if pk.PublicKeyHash == "" {
		return cmd.createPrivateKey(pk, out)
	}

	prk := cmd.keyCache.KeyByID(pk.PublicKeyHash)
	if prk == nil {
		return fmt.Errorf("failed to find private key for %s %s", pk.PublicKeyHash, pk.Location)
	}
	return cmd.showKey(prk, out)
}

func (cmd *KeyCommand) showKey(k *pemresources.PrivateKey, out io.Writer) error {
	linkId := k.PublicKeyHash
	if k.IsEncrypted {
		if cmd.Password == "" {
			if Script {
				return fmt.Errorf("key %s is encrypted and no password was specified", k.PublicKeyHash)
			}
			pwd, err := PromptPassword(fmt.Sprintf("Enter password for key: %s\t%s\nHit enter to skip", k.PublicKeyHash, k.Location))
			if err != nil {
				return err
			}
			if pwd == "" {
				return fmt.Errorf("aborted")
			}
			cmd.Password = pwd
		}
		uk, err := k.Decrypt(cmd.Password)
		if err != nil {
			return err
		}
		k = uk
	}
	puk, err := k.PublicKeyTemplate()
	if err != nil {
		return err
	}
	blkPub, err := puk.MarshalPem()
	if err != nil {
		return err
	}
	if cmd.LinkPublic && linkId != puk.PublicKeyHash {
		if blkPub.Headers == nil {
			blkPub.Headers = map[string]string{}
		}
		blkPub.Headers[pemresources.LinkedKeyHeaderKey] = linkId
	}
	blkPrv, err := k.MarshalPem()
	if err != nil {
		return err
	}

	return formatWriter.Marshal([]*pem.Block{blkPrv, blkPub}, out)
}

func (cmd *KeyCommand) createPrivateKey(prk *pemresources.PrivateKey, out io.Writer) error {
	pka := prk.PublicKeyAlgorithm
	if pka == x509.UnknownPublicKeyAlgorithm {
		pka = keytools.ParsePublicKeyAlgorithm(cmd.KeyAlgorithm)
	}
	if pka == x509.UnknownPublicKeyAlgorithm {
		return fmt.Errorf("%s is not a known PublicKeyAlgorithm. Use one of: %v", pka, keytools.PublicKeyAlgoNames[1:])
	}
	if prk.PublicKeyLength == "" {
		prk.PublicKeyLength = strconv.Itoa(cmd.KeyLength)
	}
	kl, err := strconv.Atoi(prk.PublicKeyLength)
	if err != nil {
		return fmt.Errorf("PublicKeyLength (%s) is not a valid number", prk.PublicKeyLength)
	}
	if !Script {
		if !PromptConfirm(fmt.Sprintf("Generate a new %s key of length %d?", pka, kl), false) {
			return fmt.Errorf("aborted")
		}
	}
	k, err := keytools.GenerateKey(pka, kl)
	if err != nil {
		return err
	}
	blk, err := fileformats.MarshalPrivateKey(k)
	if err != nil {
		return err
	}

	if cmd.Encrypt {
		if cmd.Password == "" {
			if Script {
				return fmt.Errorf("new key generation failed as no password supplied to encrypt new key")
			}
			pwd, err := PromptCreatePassword("Enter a password to encrypt the new key: ", minPasswordLength)
			if err != nil {
				return err
			}
			cmd.Password = pwd
		}
		blk, err = x509.EncryptPEMBlock(rand.Reader, blk.Type, blk.Bytes, []byte(cmd.Password), passwordCipher)
		if err != nil {
			return err
		}
	}
	pkt := &pemresources.PrivateKey{}
	if err = pkt.UnmarshalPem(blk); err != nil {
		return err
	}
	return cmd.showKey(pkt, out)
}

func (cmd *KeyCommand) SetKeys(keys *keycache.KeyCache) {
	cmd.keyCache = keys
}
