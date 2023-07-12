package commands

import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/eurozulu/pempal/identity"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/utils"
	"io"
	"strings"
)

var columnNames = []string{
	"identity",
	"algorithm",
	"encrypted",
	"location",
}

// keysCommand displays the available keys on the 'keypath'
// On its own, with no arguments, it lists all the identified private keys found on the keypath.
// private keys are identified by a hash of their public key.
// i.e. encrypted keys, which do not contain a public key pem in the same file are not identified.
// When arguments are given, these are used to replace the keypath for the search.
// arguments map be a directory or single file.
// If that key is known, it is output to the standard out.
// Flags:
// -all | -a  When given, shows the unidentified (encrypted) identity otherwise hidden.
// -new ["<key properties>"]  optional properties specify the key algorithm and optional length. defaults to RSA 2048
type keysCommand struct {
	AllKeys    bool `yaml:"all,omitempty"`
	keyManager identity.Keys
}

func (cmd keysCommand) Execute(args []string, out io.Writer) error {
	if len(args) == 0 {
		args = ResolvePath(CommonFlags.KeyPath)
	}
	keyManager := identity.NewKeys(args)
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()

	cout := createColumnOutput(out)
	if logger.Level() >= logger.LevelInfo {
		if err := writeColumnNames(cout); err != nil {
			return err
		}
	}
	for k := range keyManager.AllKeys(ctx) {
		if !matchKeyWithArgs(k, args) {
			continue
		}
		if err := cmd.writeKey(cout, k); err != nil {
			return err
		}
	}
	return nil
}

func (cmd *keysCommand) writeKey(out io.Writer, k identity.Key) error {
	id := k.Identity()
	if id == "" && !cmd.AllKeys {
		return nil
	}
	var ka string
	pka := utils.PublicKeyAlgorithmFromKey(k.PublicKey())
	if pka != x509.UnknownPublicKeyAlgorithm {
		ka = pka.String()
	}
	_, err := fmt.Fprintf(out, "%s,%s,%s,%s\n", id, ka, k.IsEncrypted(), k.Location())
	return err
}

func writeColumnNames(out io.Writer) error {
	_, err := fmt.Fprintf(out, strings.Join(columnNames, ","))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out)
	return err
}

func matchKeyWithArgs(k identity.Key, args []string) bool {
	if len(args) == 0 {
		return true
	}
	for _, arg := range args {
		if arg == k.Identity().String() || strings.Contains(k.Location(), arg) {
			return true
		}
	}
	return false
}

func createColumnOutput(out io.Writer) *utils.ColumnOutput {
	widths := make([]int, len(columnNames)-1)
	for i, n := range columnNames[:len(columnNames)-1] {
		widths[i] = len(n)
	}
	cout := utils.NewColumnOutput(out, widths...)
	cout.DataDelimiter = ","
	return cout
}
