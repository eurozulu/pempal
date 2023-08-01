package resources

import "strings"

// ResourceType is a specific type of x509 resource.
type ResourceType int

const (
	UnknownResourceType ResourceType = iota
	PublicKey
	PrivateKey
	CertificateRequest
	Certificate
	RevocationList
)

func (rt ResourceType) String() string {
	if rt < 0 || int(rt) >= len(ResourceTypeNames) {
		rt = UnknownResourceType
	}
	return ResourceTypeNames[rt]
}

func (rt ResourceType) PEMString() string {
	if rt < 0 || int(rt) >= len(pemTypeNames) {
		rt = UnknownResourceType
	}
	return pemTypeNames[rt]
}

var ResourceTypes = []ResourceType{
	UnknownResourceType, PublicKey, PrivateKey, CertificateRequest, Certificate, RevocationList,
}

var resourceTypes = []ResourceType{
	UnknownResourceType,
	PublicKey,
	PrivateKey,
	CertificateRequest,
	Certificate,
	RevocationList,
}

var ResourceTypeNames = []string{
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
	"crl":     RevocationList,
	"revoked": RevocationList,
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
	for i, rs := range ResourceTypeNames {
		if strings.EqualFold(s, rs) {
			return ResourceType(i)
		}
	}
	return UnknownResourceType
}

func ParsePEMType(s string) ResourceType {
	s = strings.ToUpper(s)
	for i, rs := range pemTypeNames {
		if strings.Contains(s, rs) {
			return ResourceType(i)
		}
	}
	return UnknownResourceType
}
