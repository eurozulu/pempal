package finder

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"os"
	"pempal/pemtypes"
)

// PemParser parses a given block of bytes as one or more x509 resources
// Parser may be limited by PemType as to which resources it will parse.
// By default it will parse all known pem types, providing one ore more types will lijmit it to parsing resources of those types
type PemParser interface {
	Parse(data []byte) ([]pemtypes.PemResource, error)
}

type pemParser struct {
	types        []pemtypes.PEMType
	reportErrors bool
}

// Parse takes the full content of a resource/file and attempts to split it into its content components
func (pp pemParser) Parse(data []byte) ([]pemtypes.PemResource, error) {
	// if data is a recognised container, Split data into component parts
	var dataBlocks [][]byte
	wrapper := checkIfContainer(data)
	if wrapper != nil {
		dataBlocks = wrapper.Contents(data, pp.types...)
	} else {
		// not recognised container, treat whole block as a single resource
		dataBlocks = [][]byte{data}
	}

	var res []pemtypes.PemResource
	for _, dataBlk := range dataBlocks {
		r, err := pp.parseResource(dataBlk)
		if err != nil {
			if pp.reportErrors {
				fmt.Fprintln(os.Stderr, "%v", err)
			}
			continue
		}
		res = append(res, r)
	}
	return res, nil
}

func (pp pemParser) parseResource(data []byte) (pemtypes.PemResource, error) {
	types := pp.types
	if len(types) == 0 {
		types = pemtypes.AllPemTypes[:]
	}

	// If looks like pem, attempt to parse as pem
	if looksLikePem(data) {
		res, err := pp.parseResourceAsPem(data)
		if err == nil && res != nil {
			return res, nil
		}
	}
	// not a pem, attempt to parse as der
	for _, pt := range types {
		pr := pemtypes.NewPemResource(pt)
		if err := pr.UnmarshalBinary(data); err == nil {
			return pr, nil
		} else if pp.reportErrors {
			fmt.Fprintln(os.Stderr, "%v", err)
		}
	}
	return nil, fmt.Errorf("unknown format")
}

func (pp pemParser) parseResourceAsPem(data []byte) (pemtypes.PemResource, error) {
	blks := pemtypes.ReadPEMBlocks(data, pp.types...)
	if len(blks) == 0 {
		return nil, fmt.Errorf("no pems found")
	}
	pr := pemtypes.NewPemResource(pemtypes.ParsePEMType(blks[0].Type))
	if err := pr.UnmarshalText(data); err != nil {
		return nil, err
	}
	return pr, nil
}

// container represents a location such as a file which contains one or more resources
type container interface {
	CanParse(data []byte) bool
	Contents(data []byte, pemtype ...pemtypes.PEMType) [][]byte
}

type pemContainer struct {
}

func (p pemContainer) CanParse(data []byte) bool {
	return looksLikePem(data)
}

func (p pemContainer) Contents(data []byte, pemtype ...pemtypes.PEMType) [][]byte {
	// split possible miltiple pem blocks into individual blocks of data.
	blks := pemtypes.ReadPEMBlocks(data, pemtype...)
	dataBlocks := make([][]byte, len(blks))
	for i, blk := range blks {
		dataBlocks[i] = pem.EncodeToMemory(blk)
	}
	return dataBlocks
}

var containers = [...]container{&pemContainer{}}

func looksLikePem(data []byte) bool {
	i := bytes.Index(data, []byte("-----BEGIN "))
	if i < 0 {
		return false
	}
	return bytes.Contains(data[i:], []byte("-----END "))
}

func checkIfContainer(data []byte) container {
	for _, c := range containers {
		if c.CanParse(data) {
			return c
		}
	}
	return nil
}
