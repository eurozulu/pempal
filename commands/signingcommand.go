package commands

import (
	"pempal/keycache"
)

// SigningCommand is a command which performs a signing operation, requiring a private key
// Any Command supporting this interface is given a new keycache
type SigningCommand interface {
	SetKeys(keys *keycache.KeyCache)
}
