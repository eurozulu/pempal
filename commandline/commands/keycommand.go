package commands

import (
	"bytes"
	"context"
	"fmt"
	"github.com/eurozulu/pempal/commandline/commonflags"
	"github.com/eurozulu/pempal/identity"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/utils"
	"golang.org/x/term"
	"io"
	"os"
	"sort"
	"strings"
	"syscall"
)

type keyCommand struct {
	Password   []byte `yaml:"password,omitempty"`
	Encrypt    bool   `yaml:"encrypt,omitempty"`
	Decrypt    bool   `yaml:"decrypt"`
	keyManager identity.Keys
}

func (cmd keyCommand) Execute(args []string, out io.Writer) error {
	cmd.keyManager = identity.NewKeys(commonflags.ResolvePath(commonflags.CommonFlags.KeyPath))
	keys, err := cmd.keysFromArguments(args)

	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i].Identity().String(), keys[j].Identity().String()) < 0
	})

	if err != nil {
		return err
	}
	output := commonflags.CommonFlags.Output

	for _, k := range keys {
		if !commonflags.CommonFlags.Quiet {
			if err := cmd.showKeyProperties(k, out); err != nil {
				return err
			}
		}

		if cmd.Encrypt && !k.IsEncrypted() {
			k, err = cmd.encryptKey(k)
			if err != nil {
				return err
			}
			// ensure output of new key
			if output == nil {
				output = new(string)
			}
		} else if cmd.Decrypt && k.IsEncrypted() {
			k, err = cmd.decryptKey(k)
			if err != nil {
				return err
			}
			// ensure output of new key
			if output == nil {
				output = new(string)
			}
		}
		if output != nil {
			_, err := fmt.Fprint(out, k.String())
			return err
		}
	}
	return nil
}

func (cmd keyCommand) showKeyProperties(key identity.Key, out io.Writer) error {
	var enc string
	if key.IsEncrypted() {
		enc = "encrypted"
	} else {
		enc = "unencrypted"
	}

	col := utils.NewColumnOutput(out, 32, 7, 11)
	col.ColumnDelimiter = "  "
	_, err := col.WriteSlice([]string{
		key.Identity().String(),
		key.PublicKeyAlgorithm().String(),
		enc,
		key.Location(),
	})
	col.WriteString("\n")
	return err
}

func (cmd keyCommand) encryptKey(key identity.Key) (identity.Key, error) {
	pw := cmd.Password
	if len(pw) == 0 {
		np, err := cmd.requestNewPassword(key)
		if err != nil {
			return nil, err
		}
		if len(np) == 0 {
			logger.Info("no password set, encryption of key %s aborted", key.Identity().String())
			return key, nil
		}
		pw = np
	}
	ekey, err := key.Encrypt(pw)
	if err != nil {
		return nil, err
	}
	return ekey, nil
}

func (cmd keyCommand) decryptKey(key identity.Key) (identity.Key, error) {
	pw, err := cmd.password(fmt.Sprintf("Enter password to decrypt key %s (%s)", key.Identity().String(), key.Location()))
	if err != nil {
		return nil, err
	}
	dkey, err := key.Decrypt(pw)
	if err != nil {
		return nil, err
	}
	return dkey, nil
}

func (cmd keyCommand) password(prompt string) ([]byte, error) {
	if len(cmd.Password) > 0 {
		return cmd.Password, nil
	}
	fmt.Fprint(os.Stderr, "\n%s: ", prompt)
	return term.ReadPassword(int(syscall.Stdin))
}

func (cmd keyCommand) requestNewPassword(key identity.Key) ([]byte, error) {
	for {
		pw, err := cmd.password(fmt.Sprintf("Enter password to encrypt key %s (%s)", key.Identity().String(), key.Location()))
		if err != nil {
			return nil, err
		}
		pw2, err := cmd.password(fmt.Sprintf("Re-enter password to encrypt key (%s)", key.Identity().String()))
		if err != nil {
			return nil, err
		}
		fmt.Println()
		if !bytes.Equal(pw, pw2) {
			fmt.Println("Passwords do not match!  Try again or just hit enter to abort")
			continue
		}
		return pw, nil
	}
}

func (cmd keyCommand) allKeys() []identity.Key {
	var keys []identity.Key
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	for k := range cmd.keyManager.AllKeys(ctx) {
		keys = append(keys, k)
	}
	return keys
}

func (cmd keyCommand) keysFromArguments(args []string) ([]identity.Key, error) {
	if len(args) == 0 {
		return cmd.allKeys(), nil
	}

	var keys []identity.Key
	for _, arg := range args {
		kz, err := cmd.keysFromArgument(arg)
		if err != nil {
			return nil, err
		}
		keys = append(keys, kz...)
	}
	return keys, nil
}

func (cmd keyCommand) keysFromArgument(arg string) ([]identity.Key, error) {
	// try by identity
	if k, err := cmd.keyManager.KeyByIdentity(arg); err == nil {
		return []identity.Key{k}, nil
	}
	// try by name
	if kz, err := cmd.keyManager.KeysByName(arg); err == nil {
		return kz, nil
	}
	// search all paths
	loc, err := commonflags.CommonFlags.FindInPath(arg, false)
	if err != nil {
		return nil, err
	}
	k, err := identity.NewKey(loc.Location(), loc.ResourcesAsPem(resources.PrivateKey))
	if err != nil {
		return nil, err
	}
	return []identity.Key{k}, nil
}
