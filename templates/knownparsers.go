package templates

import (
	"encoding/pem"
	"pempal/keytools"
	"pempal/templates/parsers"
)

var defaultParser = &parsers.UnknownParser{}
var keyParser = &parsers.PkiParser{}
var certParser = &parsers.CertificateParser{}
var reqParser = &parsers.CsrParser{}

var knownParsers = buildKnownParsers()

func buildKnownParsers() map[string]parsers.TemplateParser {
	m := map[string]parsers.TemplateParser{
		"": defaultParser,
	}
	combineParserMap(m, keytools.CertificateTypes, certParser)
	combineParserMap(m, keytools.CSRTypes, reqParser)
	combineParserMap(m, keytools.PrivateKeyTypes, keyParser)
	combineParserMap(m, keytools.PublicKeyTypes, keyParser)
	return m
}
func combineParserMap(m map[string]parsers.TemplateParser, m2 map[string]bool, p parsers.TemplateParser) {
	for k, v := range m2 {
		if v {
			m[k] = p
		}
	}
}

func ParseBlock(blk *pem.Block) (Template, error) {
	return NewTemplateParser(blk.Type).Parse(blk)
}

func NewTemplateParser(pemtype string) parsers.TemplateParser {
	tr := knownParsers[pemtype]
	if tr == nil {
		tr = defaultParser
	}
	return tr
}
