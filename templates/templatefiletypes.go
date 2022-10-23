package templates

import (
	"pempal/templates/parsers"
)

var TemplateFileTypes = map[string]TemplateParser{
	"yaml":     &parsers.YAMLParser{},
	"template": &parsers.YAMLParser{},
	"pem":      &parsers.PEMParser{},
	"der":      &parsers.DERCertificateParser{},
	"crt":      &parsers.DERCertificateParser{},
	"cer":      &parsers.DERCertificateParser{},
	"cert":     &parsers.DERCertificateParser{},
	"csr":      &parsers.CSRCertificateParser{},

	//"ppk":       true,
	//"pub":       true,
	//"key":       true,
	//"rsa":       true,
	//"dsa":       true,

	//"ca-bundle": true,
	//"pkcs8":     true,
	//"p7b":       true,
	//"p7s":       true,
	//"pkcs7":     true,
	//"pfx":       true,
	//"p12":       true,
	//"pkcs12":    true,
}
