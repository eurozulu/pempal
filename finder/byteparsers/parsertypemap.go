package byteparsers

import (
	"pempal/pemtypes"
)

// parserTypeMap maps pem types to the parsers which may be able to parse that type
var parserTypeMap = map[pemtypes.PEMType][]parserType{
	pemtypes.Unknown:             {},
	pemtypes.PrivateKey:          {PEM, DERPrivateKey, DERPublicKey},
	pemtypes.PrivateKeyEncrypted: {PEM},
	pemtypes.PublicKey:           {PEM, DERPublicKey},
	pemtypes.Certificate:         {PEM, DERCertificate},
	pemtypes.Request:             {PEM, DERCertificateRequest},
	pemtypes.RevocationList:      {PEM, DERRevokationList},
}

func NewByteParsers(pemtypes ...pemtypes.PEMType) []ByteParser {
	// collect a set of parsers for all the pem types
	var parserTypes []parserType
	found := map[parserType]bool{}

	for _, pmt := range pemtypes {
		prts := parserTypeMap[pmt]
		// unknown pem type invokes all the available parsers
		if len(prts) == 0 {
			return allParsers[:]
		}
		for _, prt := range prts {
			if found[prt] {
				// ignore duplicate  parser
				continue
			}
			found[prt] = true
			parserTypes = append(parserTypes, prt)
		}
	}
	parsers := make([]ByteParser, len(parserTypes))
	for i, prt := range parserTypes {
		parsers[i] = allParsers[prt-1]
	}
	return parsers
}
