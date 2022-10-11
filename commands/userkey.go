package commands

import (
	"fmt"
	"os"
	"pempal/keycache"
	"pempal/pemresources"
	"strings"
)

const ENV_KEYPATH = "PEMPAL_KEYPATH"
const ENV_USER = "PEMPAL_USERKEY"

// KeyPath holds a comma delimited list of paths to search for keys
var KeyPath = os.ExpandEnv(os.Getenv(ENV_KEYPATH))

// UserKey holds the publicKeyHash of the private key to perform signings
var UserKey = os.ExpandEnv(os.Getenv(ENV_USER))

var keyCache *keycache.KeyCache

func GetKeyPath(p []string) []string {
	if len(p) == 0 {
		p = []string{os.ExpandEnv("$PWD")}
	}
	if KeyPath == "" {
		return p
	}
	return append(p, strings.Split(KeyPath, ":")...)
}

func GetUserKey() (*pemresources.PrivateKey, error) {
	if UserKey == "" {
		return selectUserKey()
	}
	var prk *pemresources.PrivateKey
	if !strings.Contains(UserKey, "/") {
		prk = keyCache.KeyByID(UserKey)
	}
	if prk == nil {
		prk = keyCache.KeyByPath(UserKey)
	}
	if prk == nil {
		return nil, fmt.Errorf("%s user key could not be found", UserKey)
	}
	return prk, nil
}

// selectUserKey lists all the known keys and asks the user to select one.
// If user selects a key, it is returned, setting UserKEy to its keyhash.
// User selects to generate new key, returns nil, nil
// If user aborts with zero, an error is returned.
func selectUserKey() (*pemresources.PrivateKey, error) {
	if Script {
		return nil, fmt.Errorf("No user key specified in command line. Can not script run without this")
	}
	keys := keyCache.Keys(true)
	// Ask user to select or generate a key
	keys = SortKeys(keys)

	names := keyLocations(keys)
	var prompt string
	if len(names) == 0 {
		prompt = "No keys found to sign request, create one or zero to abort"
	} else {
		prompt = "Select the key to sign request, create a new one or zero to abort"
	}
	names = append([]string{"Generate new key"}, names...)

	choice := PromptChooseList(prompt, names)
	if choice < 0 {
		return nil, fmt.Errorf("aborted")
	}
	// align choice with keys slice (has extra "generate new key" entru
	choice--
	if choice >= 0 {
		k := keys[choice]
		UserKey = k.PublicKeyHash
		return k, nil
	}

	// request to create new key, return nil without err
	return nil, nil

}

func keyLocations(keys []*pemresources.PrivateKey) []string {
	locs := make([]string, len(keys))
	for i, k := range keys {
		locs[i] = k.Location
	}
	return locs
}
