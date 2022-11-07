package pemparsers

import (
	"encoding/pem"
)

type PemBlocks []*pem.Block

// PemParser represents a signle format parser, which parses bytes into one or more PemBlocks
type PemParser interface {
	// Match compares the given file path to confirm if this parser can parse it.
	// Offers scanner an oppertunity to skip loading a file, if no parsers will parse it.
	// Usually checks the file extension.
	Match(path string) bool

	// Parse attempts to read the given data as a known format, to read one or more PEMBlocks related resources.
	// data format is defined by the implentation.
	// A single data block may contain multiple resources, all of which a single parser may not be able to parses.
	// Parsers can therefore perform 'partial' parses, extracting the data they can parse and returning the data they cannot.
	// The result is a block of pems, it could parse, the bytes it didnt use, or an error.
	Parse(data []byte) (PemBlocks, []byte, error)
}

func stringIndex(s string, ss []string) int {
	for i, sz := range ss {
		if s == sz {
			return i
		}
	}
	return -1
}
