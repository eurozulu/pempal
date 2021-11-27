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

	quiet      bool
	public     bool
	private    bool
	encrypt    bool
	passwd     string
	linkPublic bool
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
	f.BoolVar(&cmd.linkPublic, "linkpublic", false, "Generate a hash id in the public key header to link the public key to an encrypted private key.")
	f.BoolVar(&cmd.linkPublic, "lp", false, "Same as 'linkpublic'")
	f.BoolVar(&cmd.private, "private", true, "When true, default, outputs the private key")
	f.BoolVar(&cmd.encrypt, "encrypt", true, "When true, default, new keys are encrypted with a password.")
	f.StringVar(&cmd.passwd, "password", "", "the passwordto encrypt or decrypt a key.")
}

func (cmd *KeyCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	if len(args) == 0 {
		return cmd.runNewKey(out)
	}

	keyCh := cmd.openPrivateKeys(ctx, args)
	for k := range keyCh {
		puk := keytools.PublicKeyFromPrivate(k)
		if err := cmd.encodePublicKey(puk, out); err != nil {
			return err
		}
		if err := cmd.encodePrivateKey(k, out); err != nil {
			return err
		}
	}
	return nil
}

func (cmd *KeyCommand) runNewKey(out io.Writer) error {
	k, err := cmd.createPrivateKey()
	if err != nil {
		return err
	}
	if err = cmd.encodePublicKey(keytools.PublicKeyFromPrivate(k), out); err != nil {
		return err
	}
	cmd.private = true
	return cmd.encodePrivateKey(k, out)
}

func (cmd *KeyCommand) createPrivateKey() (crypto.PrivateKey, error) {
	pka := keytools.ParsePublicKeyAlgorithm(cmd.keyAlgorithm)
	if pka == x509.UnknownPublicKeyAlgorithm {
		return nil, fmt.Errorf("%s is not a known PublicKeyAlgorithm. Use one of: %v", cmd.keyAlgorithm, keytools.PublicKeyAlgoNames[1:])
	}
	if !cmd.quiet {
		if !PromptConfirm(fmt.Sprintf("Generate new %s key of length %d", pka.String(), cmd.keyLength), false) {
			return nil, fmt.Errorf("aborted")
		}
	}
	return keytools.GenerateKey(pka, cmd.keyLength)
}

func (cmd *KeyCommand) openPrivateKeys(ctx context.Context, keypath []string) <-chan crypto.PrivateKey {
	ch := make(chan crypto.PrivateKey)
	go func(ch chan<- crypto.PrivateKey) {
		defer close(ch)
		kt := keytracker.KeyScanner{ShowLogs: Verbose}
		keyCh := kt.FindKeys(ctx, keypath...)
		for {
			select {
			case <-ctx.Done():
				return
			case k, ok := <-keyCh:
				if !ok {
					return
				}
				prk, err := cmd.decryptPrivateKey(k)
				if !handleError(err) {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case ch <- prk:
				}
			}
		}
	}(ch)
	return ch
}

func (cmd *KeyCommand) encodePublicKey(puk crypto.PublicKey, out io.Writer) error {
	if !cmd.public {
		return nil
	}
	blk, err := keytools.MarshalPublicKey(puk)
	if err != nil {
		return err
	}
	return pem.Encode(out, blk)
}

func (cmd *KeyCommand) encodePrivateKey(prk crypto.PrivateKey, out io.Writer) error {
	if !cmd.private {
		return nil
	}
	blk, err := keytools.MarshalPrivateKey(prk)
	if err != nil {
		return err
	}

	if cmd.encrypt {
		blk, err = cmd.encryptPrivateKey(blk)
		if err != nil {
			return err
		}
	}
	return pem.Encode(out, blk)
}

func (cmd *KeyCommand) encryptPrivateKey(blk *pem.Block) (*pem.Block, error) {
	if cmd.passwd == "" {
		if cmd.quiet {
			return nil, fmt.Errorf("no password provided to encrypt new key")
		}
		pwd, err := PromptCreatePassword("Enter a password for the new key:", encryptPwdMinLength)
		if err != nil {
			return nil, err
		}
		cmd.passwd = pwd
	}
	return x509.EncryptPEMBlock(rand.Reader, blk.Type, blk.Bytes, []byte(cmd.passwd), encryptPwdCipher)
}

func (cmd *KeyCommand) decryptPrivateKey(k keytracker.Key) (crypto.PrivateKey, error) {
	if !k.IsEncrypted() {
		return k.PrivateKey()
	}
	// encrypted key, ask for key
	if cmd.quiet || cmd.passwd == "" {
		s := fmt.Sprintf("Enter the password for encrypted key or enter to skip key %s\n%s: ", k.Location(), k)
		pwd, err := PromptPassword(s)
		if err != nil {
			return nil, err
		}
		if pwd == "" {
			return nil, fmt.Errorf("skipped")
		}
		cmd.passwd = pwd
	}
	return k.PrivateKeyDecrypted(cmd.passwd)
}
