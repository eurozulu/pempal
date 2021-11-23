package fileformats

import "encoding/pem"

// FileFormat represents a parser which parses raw bytes into one or more PEM blocks
type FileFormat interface {
	Format(by []byte) ([]*pem.Block, error)
}

var pemFormatter = &pemFileFormat{}

// FileFormats is a map of all the supported file name extensions, mapped to their respective format implementations
var FileFormats = map[string]FileFormat{
	// empty key captures files with no extension
	"":    unknownFormat{},
	"pem": pemFormatter,

	"der":  derCertificateFormat{},
	"cer":  derCertificateFormat{},
	"crt":  derCertificateFormat{},
	"cert": derCertificateFormat{},

	"csr": derCertificateRequestFormat{},

	"ppk": derKeyFormat{},
	"pub": derKeyFormat{},
	"key": derKeyFormat{},
	"rsa": derKeyFormat{},
	"dsa": derKeyFormat{},
}
