package formathelpers

import (
	"encoding/hex"
)

type HexByteFormat []byte

func (pkf HexByteFormat) String() string {
	return hex.EncodeToString(pkf)
}

func ParseHexBytes(s string) (HexByteFormat, error) {
	return hex.DecodeString(s)
}
