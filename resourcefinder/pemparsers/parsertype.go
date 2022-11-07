package pemparsers

// ParserType defines all the available pem parsets
type ParserType int

const (
	All ParserType = iota
	PEMBlocks
	DERPrivateKey
	DERPublicKey
	DERCertificate
	DERCertificateRequest
	DERRevokationList
)

var allParsers = [...]PemParser{
	&pemBlocksParser{},
	&derPrivateKeyParser{},
	&derPublicKeyParser{},
	&derCertificateParser{},
	&derCSRParser{},
	&derCRLParser{},
}

func Parsers(ptype ...ParserType) []PemParser {
	if len(ptype) == 0 {
		return allParsers[:]
	}

	var parsers []PemParser
	found := map[ParserType]bool{}
	for _, pt := range ptype {
		if pt == All {
			return allParsers[:]
		}
		if found[pt] {
			continue
		}
		parsers = append(parsers, allParsers[pt-1])
	}
	return parsers
}
