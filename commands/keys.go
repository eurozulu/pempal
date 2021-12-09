package commands

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"pempal/pemresources"
	"sort"
	"strings"
	"text/tabwriter"
)

// KeysCommand finds and lists all the private keys in the keypath
type KeysCommand struct {
	ShowIndex   bool
	ShowHeaders bool
}

func (cmd *KeysCommand) Description() string {
	lines := bytes.NewBufferString(fmt.Sprintf("lists all the available private keys in the given path(s) and in $%s if it is set.\n", ENV_KEYPATH))
	lines.WriteString("Available keys are private keys which can be identified by their public key pair.\n")
	lines.WriteString("The keys command matches the private keys to their public key, to give them an identity.\n")
	lines.WriteString("Private keys which are encrypted can not be used to generate the public key, therefore.\n")
	lines.WriteString("require a supporting public key pem linked to them.\n")
	lines.WriteString("public keys are linked either by location or a 'link' header\n")
	lines.WriteString("By location, if a private key and a public key share the same filepath, excluding any file extension, they are assumed to be a pair.\n")
	lines.WriteString("If the public key pem can not be stored in the same location, a linking key can be generated using 'key'.\n")
	lines.WriteString("A linked public key has a special header, identifying the encyprted private key it belongs to.\n")
	lines.WriteString("Keys linked in this way can be in any location, provided they both appear in the keypath.\n")
	lines.WriteString("Unencrypted private keys (for those that like to live dangeriously) do not require a matching public pair as they can generate their own. \n")
	lines.WriteString("\n")
	lines.WriteString("Output of keys shows:\n")
	lines.WriteString("<Public Key Hash>\t\t<Pem type>\t<encrypted status>\t<location of the private key>:\n")
	lines.WriteString("Public Key Hash is the SHA1 hash of the public key for that private key.\n")
	lines.WriteString("If any private, encrypted key can not be linked to its public key pair, it is unidentified.\n" +
		"Unidentified private keys show a hash of the encrypted key itself, preceeded with a \"*\"\n")

	return lines.String()
}

func (cmd *KeysCommand) Flags(f *flag.FlagSet) {
	f.BoolVar(&cmd.ShowIndex, "i", false, "Adds the index location to the path for files containing more than one resource")
	f.BoolVar(&cmd.ShowHeaders, "h", false, "output any pem headers for each key")
}

func (cmd *KeysCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	// append keypath to any given
	args = GetKeyPath(args)
	if len(args) == 0 {
		return fmt.Errorf("must provide at least one location to search for keys or set the %s environment variable with the path to search.", ENV_KEYPATH)
	}

	keyscan := pemresources.Keys{
		ShowLogs:      Verbose,
		HideAnonymous: false,
	}

	keyCh := keyscan.ScanKeys(ctx, args...)
	var keys []*pemresources.PrivateKey
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case k, ok := <-keyCh:
			if !ok {
				return nil
			}
			keys = append(keys, k)
		}
	}

	keys = SortKeys(keys)
	tw := tabwriter.NewWriter(out, 4, 4, 2, ' ', 0)
	for _, k := range keys {
		var s = bytes.NewBuffer(nil)
		s.WriteString(k.PublicKeyHash)
		s.WriteString("\t")
		s.WriteString(k.PemType)
		s.WriteString("\t")
		if k.IsEncrypted {
			s.WriteString("encrypted")
		}
		s.WriteString("\t")
		if cmd.ShowHeaders {
			s.WriteString(collectHeaders(k.PemHeaders))
			s.WriteString("\t")
		}
		l := k.Location
		if !cmd.ShowIndex {
			l = strings.TrimRight(l, "0123456789:")
		}
		s.WriteString(l)
		s.WriteString("\t")
		fmt.Fprintln(tw, s.String())
	}
	return tw.Flush()
}

func collectHeaders(heads map[string]string) string {
	buf := bytes.NewBuffer(nil)
	for k, v := range heads {
		if buf.Len() > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(k)
		buf.WriteString(" = ")
		buf.WriteString(v)
	}
	return buf.String()
}

func SortKeys(keys []*pemresources.PrivateKey) []*pemresources.PrivateKey {
	// sort results by their filepaths
	sort.Slice(keys, func(i, j int) bool {
		il := keys[i].Location
		jl := keys[j].Location
		ls := []string{il, jl}
		sort.Strings(ls)
		return ls[0] == il
	})
	return keys
}
