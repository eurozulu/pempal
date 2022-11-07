package resourcefinder

import (
	"fmt"
	"os"
	"pempal/resourcefinder/pemparsers"
	"pempal/resources"
)

var ErrUnknownFormat = fmt.Errorf("unknown format")

// typeParsers maps single ResourceTypes to their respective Parsers
// these are the parsers specific to the type.  Common parser, are always used.
var typeParsers = map[resources.ResourceType][]pemparsers.ParserType{
	resources.Unknown:        {pemparsers.All},
	resources.Key:            {pemparsers.DERPrivateKey, pemparsers.DERPublicKey},
	resources.Certificate:    {pemparsers.DERCertificate},
	resources.PublicKey:      {pemparsers.DERPublicKey},
	resources.Request:        {pemparsers.DERCertificateRequest},
	resources.RevocationList: {pemparsers.DERRevokationList},
}

// commonParsers are the list of parsers always included in every scan.
var commonParsers = []pemparsers.ParserType{pemparsers.PEMBlocks}

type ResourceParser interface {
	CanParse(location string) bool
	ParseLocation(location string) (resources.Resources, error)
}

type resourceParser struct {
	parsers []pemparsers.PemParser
}

func (rp resourceParser) CanParse(location string) bool {
	for _, p := range rp.parsers {
		if p.Match(location) {
			return true
		}
	}
	return false
}

func (rp resourceParser) ParseLocation(location string) (resources.Resources, error) {
	// match parsers using filename/extension
	matchingParsers := rp.matchParsers(location)
	if len(matchingParsers) == 0 {
		return nil, ErrUnknownFormat
	}

	data, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}

	var blks pemparsers.PemBlocks
	// Keep offering data to all the matching parsers until no data left.
	// This is so order of parsers is not important
	for len(data) > 0 {
		var bs pemparsers.PemBlocks
		bs, data = rp.parseData(matchingParsers, data)
		if len(bs) == 0 {
			// no parser could parse data.
			break
		}
		blks = append(blks, bs...)
	}
	var res resources.Resources
	for _, b := range blks {
		res = append(res, resources.NewResource(location, b))
	}
	return res, nil
}

func (rp resourceParser) parseData(parsers []pemparsers.PemParser, data []byte) (pemparsers.PemBlocks, []byte) {
	var blks pemparsers.PemBlocks
	var err error
	for _, parser := range parsers {
		if len(data) == 0 {
			break
		}
		blks, data, err = parser.Parse(data)
		if err != nil {
			// ignore as parser doesn't recognise file
			continue
		}
	}
	return blks, data
}

func (rp resourceParser) matchParsers(path string) []pemparsers.PemParser {
	var matchingParsers []pemparsers.PemParser
	for _, parser := range rp.parsers {
		if !parser.Match(path) {
			continue
		}
		matchingParsers = append(matchingParsers, parser)
	}
	return matchingParsers
}

func parsersForTypes(resourceType ...resources.ResourceType) []pemparsers.PemParser {
	if len(resourceType) == 0 {
		return pemparsers.Parsers()
	}
	var ptypes []pemparsers.ParserType
	copy(ptypes, commonParsers)

	for _, rt := range resourceType {
		if rt == resources.Unknown {
			return pemparsers.Parsers()
		}
		ptypes = append(ptypes, typeParsers[rt]...)
	}
	return pemparsers.Parsers(ptypes...)
}
