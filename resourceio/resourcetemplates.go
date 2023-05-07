package resourceio

import (
	"crypto/x509"
	"pempal/templates"
)

var standardTemplates = map[string][]byte{
	"privatekey":               []byte(PrivateKey),
	"privatekey-rsa":           []byte(PrivateKeyRSA),
	"privatekey-ecdsa":         []byte(PrivateKeyECDSA),
	"certificaterequest":       []byte(CertificateRequest),
	"certificate":              []byte(Certificate),
	"certificate-default":      []byte(DefaultCertificate),
	"certificate-ca":           []byte(CertificateCA),
	"certificate-intermediate": []byte(CertificateIntermediate),
	"revokationlist":           []byte(RevokationList),
	"dn":                       []byte(DN),
}

const PrivateKey = `
resource-type: privatekey
public-key-algorithm: 
is-encrypted: true
`
const CertificateRequest = `
resource-type: certificaterequest
subject: null
signature: 
signature-algorithm:
public-key-algorithm:
public-key:
version:
`
const Certificate = `
resource-type: certificate
subject: null
issuer: null
signature: 
signature-algorithm:
public-key-algorithm:
public-key:
version:
serial-number:
not-before:
not-after:
is-ca: false
`
const DN = `common-name:
country: []
organization: []
organizational-unit: []
locality: []
province: []
street-address: []
postal-code: []
serial-number:
`

var a x509.PublicKeyAlgorithm

const DefaultCertificate = `
#extends certificate
#imports DN "subject"
#imports DN "issuer"

version: 1
subject: {{ range $key, $value := .subject }}
  {{ $key }}: {{ $value }}{{ end }}
  
issuer: {{ range $key, $value := .issuer }}
  {{ $key }}: {{ $value }}{{ end }}

signature-algorithm: SHA256WithRSA
public-key-algorithm: RSA

not-before: {{ now }}
not-after: {{ nowPlusDays 365 }}
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

const RevokationList = `#revokationliste
`

const PrivateKeyRSA = `#extends privatekey
public-key-algorithm: RSA
key-param: 2048
`
const PrivateKeyECDSA = `#extends privatekey
public-key-algorithm: ECDSA
key-param: P256
`

func NewResourceTemplateManager(root string) (templates.TemplateManager, error) {
	tm, err := templates.NewTemplateManager(root)
	if err != nil {
		return nil, err
	}
	for name, data := range standardTemplates {
		tm.AddTemplate(name, data)
	}
	return tm, nil
}
