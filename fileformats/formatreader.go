package fileformats

import (
	"encoding/pem"
	"fmt"
)

var FormatTypes = []string{
	"pem",
	"der",
	"yaml",
	"pkcs12",
	"pkcs7",
}

var FormatReaders = map[string][]FormatReader{
	"yaml":   {&yamlTemplateReader{}},
	"pem":    {&pemReader{}},
	"der":    {&derCertificateReader{}, &derCertificateRequestReader{}, &derPrivateKeyReader{}, &derPublicKeyReader{}},
	"pkcs12": {&pkcs12Format{}},
	"pkcs7":  {&pkcs7Format{}},
}

// FormatReader reads a raw file contents and attempts to parse it into one or more PEM blocks.
type FormatReader interface {
	Unmarshal(by []byte) ([]*pem.Block, error)
}

type formatReader struct {
	formats []FormatReader
}

func (f formatReader) Unmarshal(by []byte) ([]*pem.Block, error) {
	var blks []*pem.Block
	for _, fr := range f.formats {
		blks, _ = fr.Unmarshal(by)
		if len(blks) > 0 {
			return blks, nil
		}
	}
	return nil, fmt.Errorf("unknown format")
}

// NewFormatReader gets a new Reader for the given formats.
// formats should be "pem", "der" or "yaml" or any combination of them all.
// If no formats are given, all formats available are used.
// When unmarsahlling bytes, each format is tried, in the order given, until one produces pem blocks.
func NewFormatReader(formats ...string) FormatReader {
	var fmts []FormatReader
	// No args, use ALL formats
	if len(formats) == 0 {
		formats = FormatTypes
	}
	for _, f := range formats {
		fm, ok := FormatReaders[f]
		if !ok {
			panic(fmt.Errorf("%s is not a known format", f))
		}
		fmts = append(fmts, fm...)
	}
	if len(fmts) == 1 {
		return fmts[0]
	}
	return &formatReader{formats: fmts}
}
