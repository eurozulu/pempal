package cmd

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"pempal/keytools"
	"pempal/keytracker"
)

const encryptPwdMinLength = 6
const encryptPwdCipher = x509.PEMCipherAES256

// KeyCommand creates new keys.
// When used with no arguments, it generates a new private/public key pair using the keyAlgorithm and keyLength flags.
// When used with an argument, it searches the given arguments as key paths searching for private keys, attempting to show the corresponding public key file for each.
type KeyCommand struct {
	keyAlgorithm string
	keyLength    int

	quiet   bool
	public  bool
	encrypt bool
	passwd  string
}

func (cmd *KeyCommand) Description() string {
	return "creates a new key pair"
}

func (cmd *KeyCommand) Flags(f *flag.FlagSet) {
	algos := fmt.Sprintf("The key Algorithm to generate.\n\tmust be one of %v", keytools.PublicKeyAlgoNames[1:])
	f.StringVar(&cmd.keyAlgorithm, "a", "rsa", algos)
	f.IntVar(&cmd.keyLength, "l", 2048, "The length/curve to use for the key")
	f.BoolVar(&cmd.quiet, "q", false, "surpresses the confirmation prompt to generate new key or password request for encrypted keys to generate public keys")
	f.BoolVar(&cmd.public, "public", true, "When true, default, generates the corresponding public key pem. If private key already exists, requests password to decrypt")
	f.BoolVar(&cmd.public, "pubout", true, "same as 'public'")
	f.BoolVar(&cmd.encrypt, "encrypt", true, "When true, default, new keys are encrypted with a password.")
	f.StringVar(&cmd.passwd, "password", "", "the passwordto encrypt or decrypt a key.")
}

func (cmd *KeyCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	if len(args) == 0 {
		return cmd.createPrivateKey(out)
	}
	// With args, treat as keypath and create public keys for any found private keys
	return cmd.createPublicKeys(ctx, out, args)
}

func (cmd *KeyCommand) createPrivateKey(out io.Writer) error {
	pka := keytools.ParsePublicKeyAlgorithm(cmd.keyAlgorithm)
	if pka == x509.UnknownPublicKeyAlgorithm {
		return fmt.Errorf("%s is not a known PublicKeyAlgorithm. Use one of: %v", cmd.keyAlgorithm, keytools.PublicKeyAlgoNames[1:])
	}
	if !cmd.quiet {
		if !PromptConfirm(fmt.Sprintf("Generate new %s key of length %d", pka.String(), cmd.keyLength), false) {
			return nil
		}
	}
	prk, err := keytools.GenerateKey(pka, cmd.keyLength)
	if err != nil {
		return err
	}
	blk, err := keytools.MarshalPrivateKey(prk)
	if err != nil {
		return err
	}

	if cmd.public {
		// make public prior to encrypting
		if err = cmd.makePublicKey(keytracker.NewKey(blk), out); err != nil {
			return err
		}
	}
	if cmd.encrypt {
		if cmd.passwd == "" {
			if cmd.quiet {
				return fmt.Errorf("no password provided to encrypt new key")
			}
			pwd, err := PromptCreatePassword("Enter a password for the new key:", encryptPwdMinLength)
			if err != nil {
				return err
			}
			cmd.passwd = pwd
		}
		eb, err := x509.EncryptPEMBlock(rand.Reader, blk.Type, blk.Bytes, []byte(cmd.passwd), encryptPwdCipher)
		if err != nil {
			return err
		}
		blk = eb
	}
	if err = pem.Encode(out, blk); err != nil {
		return err
	}
	k := keytracker.NewKey(blk)
	if err != nil {
		return err
	}
	return cmd.makePublicKey(k, out)
}

func (cmd *KeyCommand) createPublicKeys(ctx context.Context, out io.Writer, keypath []string) error {
	kt := keytracker.KeyTracker{ShowLogs: Verbose}
	keyCh := kt.FindKeys(ctx, keypath...)
	for {
		select {
		case <-ctx.Done():
			return nil
		case k, ok := <-keyCh:
			if !ok {
				return nil
			}
			if err := cmd.makePublicKey(k, out); err != nil {
				return err
			}
		}
	}
}

func (cmd *KeyCommand) makePublicKey(k keytracker.Key, out io.Writer) error {
	var prk crypto.PrivateKey
	if !k.IsEncrypted() {
		pk, err := k.PrivateKey()
		if err != nil {
			return err
		}
		prk = pk
	} else {
		if cmd.passwd == "" {
			s := fmt.Sprintf("Enter the password for encrypted key or enter to skip key %s\n%s: ", k.Location(), k)
			pwd, err := PromptPassword(s)
			if err != nil {
				return err
			}
			if pwd == "" {
				return nil
			}
			cmd.passwd = pwd
		}
		pk, err := k.PrivateKeyDecrypted(cmd.passwd)
		if err != nil {
			return err
		}
		prk = pk
	}
	blk, err := keytools.MarshalPublicKey(keytools.PublicKeyFromPrivate(prk))
	if err != nil {
		return err
	}
	return pem.Encode(out, blk)
}
