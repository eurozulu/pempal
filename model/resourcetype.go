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
	"CERTIFICATE REQUEST",
	"CERTIFICATE",
	"X509 CRL",
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
