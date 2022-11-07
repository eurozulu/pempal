package byteparsers

// parserType defines all the available format parsets
type parserType int

const (
	All parserType = iota
	PEM
	DERPrivateKey
	DERPublicKey
	DERCertificate
	DERCertificateRequest
	DERRevokationList
)

// order of array should match the order of the iota consts.
var allParsers = [...]ByteParser{
	&pemParser{},
	&derPrivateKeyParser{},
	&derPublicKeyParser{},
	&derCertificateParser{},
	&derCSRParser{},
	&derCRLParser{},
}
