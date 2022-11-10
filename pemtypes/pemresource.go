package pemtypes

import (
	"encoding"
	"encoding/pem"
	"fmt"
)

// PemResource represents a single 'pem' item, such as a certificate, csr, public or private key etc.
// Each is capable of Binary (Un)Marsahaling and Text (Un)Marsahaling.
// Binary is usually ASN.1 DER, depending on resource.
// Text Marshals PEM encoded text
type PemResource interface {
	fmt.Stringer
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

func NewPemResource(pt PEMType) PemResource {
	switch pt {
	case Certificate:
		return &certificateType{}
	case PrivateKey:
		return &privateKeyType{}
	case PrivateKeyEncrypted:
		// use empty pem block to indicate it should be encrypted
		return &privateKeyType{encrypted: &pem.Block{}}
	case PublicKey:
		return &publicKeyType{}

	case Request:
		return &requestType{}
	case RevocationList:
		return &revokeListType{}
	case Name:
		return &dnameType{}
	default:
		return nil
	}
}
