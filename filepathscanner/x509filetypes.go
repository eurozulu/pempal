package filepathscanner

// X509FileTypes lists all the file extensions to consider as candidates for parsing into resources.
var X509FileTypes = map[string]bool{
	// empty key captures files with no extension
	"":          true,
	"pem":       true,
	"ca-bundle": true,
	"crt":       true,
	"der":       true,
	"cer":       true,
	"cert":      true,
	"csr":       true,
	"ppk":       true,
	"pub":       true,
	"key":       true,
	"rsa":       true,
	"dsa":       true,
	"pkcs8":     true,
	"p7b":       true,
	"p7s":       true,
	"pkcs7":     true,
	"pfx":       true,
	"p12":       true,
	"pkcs12":    true,
}
