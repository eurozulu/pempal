package resources

import "encoding/pem"

// PemResource represents a single location containing one or more PEMs
type PemResource struct {
	Path    string
	Content []*pem.Block
}
