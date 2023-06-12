package model

import (
	"bytes"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/templates"
	"github.com/go-yaml/yaml"
	"reflect"
	"strings"
)

const defaultkey = `
#extends privatekey
public-key-algorithm: RSA
key-param: 2048
`
const client_certificate_request = `
#imports - key
#imports clientName
#extends certificaterequest
version: 1.0
public-key-algorithm: RSA
signature-algorithm: SHA512WithRSA
public-key: {{ .key }}
subject: {{ .clientName }}
`
const client_certificate = `
#imports issuerName
#extends clientcertificaterequest,certificate
serial-number: 1
issuer: {{ .issuerName }}
not-before: {{ now }}
not-after: {{ nowPlusDays 365 }}
is-ca: false
basic-constraints-valid: true
max-path-len: 0
max-path-len-zero: true
key-usage: KeyUsageDigitalSignature|KeyUsageContentCommitment|KeyUsageKeyEncipherment|KeyUsageDataEncipherment|KeyUsageKeyAgreement|KeyUsageCRLSign
`
const newClientCertificate = `
#extends defaultKey
#imports defaultIssuer issuerName
---
#extends clientcertificate
`

// DefaultResourceTemplates contains the named resource templates for each resource type.
// These define the base properties used to create that resource type.
var DefaultResourceTemplates map[string][]byte

func init() {
	DefaultResourceTemplates = map[string][]byte{}

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

func EqualsResourceTemplate(data []byte) string {
	for k, v := range DefaultResourceTemplates {
		if bytes.Equal(data, v) {
			return k
		}
	}
	return ""
}

// DetectResourceType scans the given templates from first to last to locate the first 'resource template'.
// A resource template is one of the in built templates which relates to a specific resource type.
// see #DefaultResourceTemplates
func DetectResourceType(temps ...templates.Template) ResourceType {
	for _, t := range temps {
		s := EqualsResourceTemplate(t.Bytes())
		if s == "" {
			continue
		}
		return ParseResourceType(s)
	}
	return Unknown
}
