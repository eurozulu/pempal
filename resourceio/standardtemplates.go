package resourceio

var standardTemplates = map[string][]byte{
	"privatekey-rsa":           []byte(PrivateKeyRSA),
	"privatekey-ecdsa":         []byte(PrivateKeyECDSA),
	"certificate-ca":           []byte(CertificateCA),
	"certificate-intermediate": []byte(CertificateIntermediate),
}

const PrivateKeyRSA = `#type privatekey
public-key-algorithm: RSA
key-param: 2048
`
const PrivateKeyECDSA = `#type privatekey
public-key-algorithm: ECDSA
key-param: P256
`
const CertificateCA = `#extends certificate
is-ca: true
basic-constraints-valid: true
max-path-len: -1
max-path-len-zero: false
`
const CertificateIntermediate = `#extends certificate-ca
max-path-len: 1
`
