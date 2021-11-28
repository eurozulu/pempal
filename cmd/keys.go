package cmd

import (
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
	ListCerts bool
	ShowHash  bool
	Recursive bool
}

func (cmd *KeysCommand) Description() string {
	return fmt.Sprintf("lists the private keys available in the given path(s) and in $%s if it is set. See keypath for more details", ENV_KeyPath)
}

func (cmd *KeysCommand) Flags(f *flag.FlagSet) {
	f.BoolVar(&cmd.ListCerts, "id", false, "lists the keys associated certificates")
	f.BoolVar(&cmd.ShowHash, "hash", false, "Displays a sha1 hash of the public key")
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
		fmt.Fprintf(out, s)
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
	kt := keytracker.KeyTracker{ShowLogs: Verbose, Recursive: recursive}
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
