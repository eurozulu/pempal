package cmd

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"pempal/keytracker"
	"sort"
	"strings"
	"text/tabwriter"
)

const ENV_KeyPath = "PP_KEYPATH"

var KeyPath = strings.TrimSpace(os.Getenv(ENV_KeyPath))

// KeysCommand finds and lists all the private keys in the keypath
type KeysCommand struct {
	Recursive bool
}

func (cmd *KeysCommand) Description() string {
	lines := bytes.NewBufferString(fmt.Sprintf("lists all the available private keys in the given path(s) and in $%s if it is set.\n", ENV_KeyPath))
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
	f.BoolVar(&cmd.Recursive, "r", false, "search subdirectories recursively")
}

func (cmd *KeysCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	// append keypath to any given
	args = GetKeyPath(args)
	if len(args) == 0 {
		return fmt.Errorf("must provide at least one location to search for keys or set the %s environment variable with the path to search.", ENV_KeyPath)
	}
	keys := SortKeys(Keys(ctx, args, cmd.Recursive))
	//TODO, fix column sizing
	tw := tabwriter.NewWriter(out, 2, 1, 4, ' ', 0)
	for _, s := range KeyList(keys) {
		fmt.Fprintf(out, "%s\n", s)
	}
	return tw.Flush()
}

func KeyList(keys []keytracker.Key) []string {
	names := make([]string, len(keys))
	for i, k := range keys {
		enc := ""
		if k.IsEncrypted() {
			enc = "(encrypted)"
		}
		names[i] = fmt.Sprintf("%s\t%s\t%s\t%s", k.String(), k.Type(), enc, k.Location())
	}
	return names
}

func Keys(ctx context.Context, keypath []string, recursive bool) []keytracker.Key {
	kt := keytracker.KeyTracker{ShowLogs: VerboseFlag, Recursive: recursive}
	keyCh := kt.FindKeys(ctx, keypath...)

	var found []keytracker.Key
	for {
		select {
		case <-ctx.Done():
			return nil
		case id, ok := <-keyCh:
			if !ok {
				return found
			}
			found = append(found, id)
		}
	}
}

func SortKeys(keys []keytracker.Key) []keytracker.Key {
	// sort results by their filepaths
	sort.Slice(keys, func(i, j int) bool {
		il := keys[i].Location()
		jl := keys[j].Location()
		ls := []string{il, jl}
		sort.Strings(ls)
		return ls[0] == il
	})
	return keys
}

func GetKeyPath(p []string) []string {
	if len(p) == 0 {
		p = []string{os.ExpandEnv("$PWD")}
	}
	if KeyPath == "" {
		return p
	}
	return append(p, strings.Split(KeyPath, ":")...)
}
