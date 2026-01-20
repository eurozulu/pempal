package model

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

type Fingerprint [20]byte

func (f Fingerprint) String() string {
	return hex.EncodeToString(f[:])
}

func (f Fingerprint) Equals(s string) bool {
	return strings.EqualFold(f.String(), s)
}

func (f Fingerprint) Match(s string) bool {
	return strings.Contains(f.String(), s)
}

func NewFingerPrint(data []byte) Fingerprint {
	return sha1.Sum(data)
}

func ParseFingerPrint(s string) (Fingerprint, error) {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return Fingerprint{}, err
	}
	if len(bytes) != 20 {
		return Fingerprint{}, fmt.Errorf("Fingerprint must be 20 bytes")
	}
	return Fingerprint(bytes), nil
}

func IsFingerPrint(s string) bool {
	_, err := ParseFingerPrint(s)
	return err == nil
}

func IsPartialFingerPrint(s string) bool {
	if s == "" {
		return false
	}
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return false
	}
	return len(bytes) <= 20
}
