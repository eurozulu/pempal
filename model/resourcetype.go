package model

import "strings"

type ResourceType int

const (
	Unknown ResourceType = iota
	PublicKey
	PrivateKey
	CertificateRequest
	Certificate
	RevokationList
)

func (rt ResourceType) String() string {
	if rt < 0 || int(rt) >= len(resourceTypeNames) {
		rt = Unknown
	}
	return resourceTypeNames[rt]
}

func (rt ResourceType) PEMString() string {
	if rt < 0 || int(rt) >= len(pemTypeNames) {
		rt = Unknown
	}
	return pemTypeNames[rt]
}

var resourceTypes = []ResourceType{
	Certificate, PrivateKey, CertificateRequest, PublicKey, RevokationList,
}

var resourceTypeNames = []string{
	"Unknown",
	"PublicKey",
	"PrivateKey",
	"CertificateRequest",
	"Certificate",
	"RevocationList",
}
var pemTypeNames = []string{
	"UNKNOWN",
	"PUBLIC KEY",
	"PRIVATE KEY",
	"CERTIFICATE REQUEST", // Note PEMs are parsed using 'Contains' therefore 'CERTIFICATE REQUEST' MUST appear before 'CERTIFICATE'
	"CERTIFICATE",
	"X509 CRL",
}

var resourceTypeAliasis = map[string]ResourceType{
	"cert":    Certificate,
	"key":     PrivateKey,
	"puk":     PublicKey,
	"request": CertificateRequest,
	"csr":     CertificateRequest,
	"crl":     RevokationList,
	"revoked": RevokationList,
}

func ContainsType(t ResourceType, types []ResourceType) bool {
	for _, rt := range types {
		if t == rt {
			return true
		}
	}
	return false
}

func ParseResourceType(s string) ResourceType {
	alias, ok := resourceTypeAliasis[s]
	if ok {
		return alias
	}
	for i, rs := range resourceTypeNames {
		if strings.EqualFold(s, rs) {
			return ResourceType(i)
		}
	}
	return Unknown
}

func ParsePEMType(s string) ResourceType {
	s = strings.ToUpper(s)
	for i, rs := range pemTypeNames {
		if strings.Contains(s, rs) {
			return ResourceType(i)
		}
	}
	return Unknown
}
