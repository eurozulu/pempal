package model

import (
	"github.com/eurozulu/pempal/logger"
	"github.com/go-yaml/yaml"
	"reflect"
	"strings"
)

const organisation_name = `
saubject.organisation: Your Organisation Name
`
const client_issuer_name = `
issuer.common-name: My Client Issuer Certificate CN
`

const default_key = `
#extends privatekey
public-key-algorithm: RSA
key-size: 2048
`
const default_certificate = `
#extends certificate
public-key-algorithm: RSA
key-param: 2048
signature-algorithm: SHA512-RSA
not-before: {{ now }}
not-after: {{ nowPlusDays 365 }}
is-ca: false
`

const client_certificate = `
#extends default_certificate
#imports organisation_name,client_issuer_name
serial-number: 1
subject: {{ .organisationName }}
issuer: {{ .clientIssueName }}
basic-constraints-valid: true
max-path-len: 0
max-path-len-zero: true
key-usage: KeyUsageDigitalSignature|KeyUsageContentCommitment|KeyUsageKeyEncipherment|KeyUsageDataEncipherment|KeyUsageKeyAgreement|KeyUsageCRLSign
`

// DefaultResourceTemplates contains the named resource templates for each resource type.
// These define the base properties used to create that resource type.
var DefaultResourceTemplates = map[string][]byte{
	"key":                 []byte(default_key),
	"default_certificate": []byte(default_certificate),
	"client_certificate":  []byte(client_certificate),
	"organisation_name":   []byte(organisation_name),
	"client_issuer_name":  []byte(client_issuer_name),
}

func init() {
	// Add base resource type templates
	for _, rt := range resourceTypes {
		dto := NewDTOForResourceType(rt)
		data, err := yaml.Marshal(dto)
		if err != nil {
			logger.Error("Failed to marshal resource template for %s, %v", reflect.TypeOf(dto).String(), err)
			continue
		}
		DefaultResourceTemplates[strings.ToLower(rt.String())] = data
	}
}
