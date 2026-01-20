package model

import (
	"encoding"
	"fmt"
)

// PemResource represents a single Pem Item:
// Resources are Certificate, Private or Public key, CSR or CRL
// Resources can also be marshalled and unmarshalled as text to/from pem encoding.
type PemResource interface {
	// The type of Resource
	// A Certificate, Key, CSR or CRL.
	ResourceType() ResourceType
	// Fingerprint returns the unique id for this reqource.
	// The Fingerprint is a sha1 hash of the der encoded resource.
	Fingerprint() Fingerprint
	// String returns a summary of the Resource
	fmt.Stringer
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}
