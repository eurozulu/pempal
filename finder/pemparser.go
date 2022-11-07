package finder

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"log"
	"pempal/finder/byteparsers"
	"pempal/pemtypes"
)

type pemParser struct {
	types   []pemtypes.PEMType
	parsers []byteparsers.ByteParser
}

func (pp pemParser) Parse(path string, data []byte) (Location, error) {
	var blocks []*pem.Block
	for len(data) > 0 {
		blks, rest := pp.parseEachFormat(data)
		if len(blks) == 0 {
			break
		}
		blocks = append(blocks, blks...)
		data = rest
	}
	if len(blocks) == 0 {
		return nil, fmt.Errorf("no pems found in %s", path)
	} else if len(bytes.TrimSpace(data)) > 0 {
		log.Printf("unexpected data found and ignored (%d bytes)", len(data))
	}
	return &pemLocation{
		path: path,
		pems: blocks,
	}, nil
}

func (pp pemParser) MatchPath(path string) bool {
	for _, p := range pp.parsers {
		if p.MatchPath(path) {
			return true
		}
	}
	return false
}

func (pp pemParser) FilterLocation(rl Location) Location {
	pl, ok := rl.(*pemLocation)
	if !ok {
		return nil
	}
	if len(pp.types) == 0 || pp.containsType(pemtypes.Unknown) {
		return rl
	}
	var blocks []*pem.Block
	for _, b := range pl.pems {
		if !pp.containsType(pemtypes.ParsePEMType(b.Type)) {
			continue
		}
		blocks = append(blocks, b)
	}
	if len(blocks) == 0 {
		return nil
	}
	return &pemLocation{
		path: rl.Path(),
		pems: blocks,
	}
}

func (pp pemParser) parseEachFormat(data []byte) ([]*pem.Block, []byte) {
	for _, p := range pp.parsers {
		blks, rest, _ := p.Parse(data)
		if len(blks) == 0 {
			continue
		}
		return blks, rest
	}
	return nil, data
}

func (pp pemParser) containsType(pt pemtypes.PEMType) bool {
	for _, p := range pp.types {
		if p == pt {
			return true
		}
	}
	return false
}

func newPemParser(pemType ...pemtypes.PEMType) *pemParser {
	if len(pemType) == 0 {
		pemType = []pemtypes.PEMType{pemtypes.Unknown}
	}
	return &pemParser{
		types:   pemType,
		parsers: byteparsers.NewByteParsers(pemType...),
	}
}
