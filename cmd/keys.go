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

type KeysCommand struct {
	listCerts bool
	showHash  bool
}

func (cmd *KeysCommand) Description() string {
	return fmt.Sprintf("lists the private keys available in the given path(s) and in $%s if it is set. See keypath for more details", ENV_KeyPath)
}

func (cmd *KeysCommand) Flags(f *flag.FlagSet) {
	f.BoolVar(&cmd.listCerts, "id", false, "lists the keys associated certificates")
	f.BoolVar(&cmd.showHash, "hash", false, "Displays a sha1 hash of the public key")
}

func (cmd *KeysCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	// append keypath to any gioven
	if KeyPath != "" {
		args = append(args, strings.Split(KeyPath, ":")...)
	}
	if len(args) == 0 {
		return fmt.Errorf("must provide at least one location to search for keys or set the %s environment variable with the path to search.", ENV_KeyPath)
	}
	keys := Keys(ctx, args)
	//TODO, fix column sizing
	tw := tabwriter.NewWriter(out, 2, 1, 4, ' ', 0)
	for _, key := range keys {
		fmt.Fprintf(out, "%s\t%s\t%s", key.Type(), key.String(), key.Location())
		pl := key.PublicLocation()
		if pl != "" {
			fmt.Fprintf(out, "\n\tPublic key location: %s", pl)
		}
	}
	return tw.Flush()
}

func Keys(ctx context.Context, keypath []string) []keytracker.Key {
	kt := keytracker.KeyTracker{ShowLogs: Verbose}
	keyCh := kt.FindKeys(ctx, keypath...)

	var found []keytracker.Key
	for {
		select {
		case <-ctx.Done():
			return nil
		case id, ok := <-keyCh:
			if !ok {
				return sortKeys(found)
			}
			found = append(found, id)
		}
	}
}

func sortKeys(keys []keytracker.Key) []keytracker.Key {
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
