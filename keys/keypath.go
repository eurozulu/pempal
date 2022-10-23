package keys

import (
	"bytes"
	"context"
	"pempal/filepath"
	"strings"
)

var keyFileTypes = []string{"", "key", "rsa", "prk"}

// KeyFilePath represents the collection of directories containing the keys.
// Keypath is non-recursive, so any subdirctories found in a key path directory will be ignored.
type KeyFilePath interface {
	Path() []string
	Find(name string) Key
	FindAll(ctx context.Context) <-chan Key
}

type keyParser[T Key] struct{}

func (kp keyParser[T]) Parse(p string) (Key, error) {
	k, err := NewKey(p)
	return k, err
}

type keyMatcher[T Key] struct{}

func (km keyMatcher[T]) Match(name string, k Key) bool {
	// first check if its a file path location match
	if k.Location() == name {
		return true
	}

	// next check if its a public key hash match
	pkhash, err := PublicKeyHash(k.PublicKey())
	if err == nil && bytes.Equal([]byte(name), pkhash) {
		return true
	}

	// next check if the name matches
	n := filepath.NameFromPath(k.Location())
	if strings.ToLower(name) == n {
		return true
	}
	return false
}

// NewKeyPath creates a new KeyFilePath using the given, colon delimited, list of directory names
func NewKeyPath(keypath string) KeyFilePath {
	if !strings.HasPrefix(keypath, "./:") {
		keypath = strings.Join(append([]string{}, "./", keypath), ":")
	}
	var parser filepath.ObjectParser[Key]
	parser = &keyParser[Key]{}
	matcher := &keyMatcher[Key]{}
	return filepath.NewObjectPath(keypath, parser, matcher, keyFileTypes...)
}
