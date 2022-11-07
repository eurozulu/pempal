package byteparsers

import (
	"encoding/pem"
)

type ByteParser interface {
	MatchPath(path string) bool
	Parse(data []byte) ([]*pem.Block, []byte, error)
}
