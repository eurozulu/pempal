package model

import (
	"encoding/base64"
)

type Base64Binary []byte

func (b Base64Binary) MarshalText() (text []byte, err error) {
	return []byte(EncodeAsBase64(b)), nil
}

func (b *Base64Binary) UnmarshalText(text []byte) error {
	by, err := DecodeAsBase64(string(text))
	if err != nil {
		return err
	}
	*b = by
	return nil
}

func EncodeAsBase64(text []byte) string {
	return base64.StdEncoding.EncodeToString(text)
}

func DecodeAsBase64(text string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(text)
}
