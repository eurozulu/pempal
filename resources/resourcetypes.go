package resources

/*
 - CERTIFICATE [RFC5280]
 - X509 CRL [RFC5280]
- CERTIFICATE REQUEST [RFC2986]
PKCS7 [RFC2315]
CMS [RFC5652]
- PRIVATE KEY [RFC5208] [RFC5958]
ENCRYPTED PRIVATE KEY [RFC5958]
ATTRIBUTE CERTIFICATE [RFC5755]
- PUBLIC KEY [RFC5280]
*/

import (
	"log"
	"regexp"
	"strings"
)

type ResourceType int

const (
	Unknown ResourceType = iota
	PrivateKey
	PrivateKeyEncrypted
	Request
	Certificate
	RevocationList
	PublicKey
	Name
)

var allResourceTypes = [...]ResourceType{PrivateKey, PrivateKeyEncrypted, Request, Certificate, RevocationList, PublicKey, Name}
var resourceTypeNames = [...]string{"", "PRIVATE KEY", "ENCRYPTED PRIVATE KEY", "CERTIFICATE REQUEST", "CERTIFICATE", "X509 CRL", "PUBLIC KEY", "DISTINGUISHED NAME"}

func (r ResourceType) String() string {
	return resourceTypeNames[r]
}

var resourceTypePatterns = map[string]ResourceType{
	"*PRIVATE KEY*":         PrivateKey,
	"*PUBLIC KEY*":          PublicKey,
	"*CERTIFICATE*REQUEST*": Request,
	"*CRL*":                 RevocationList,
}

func ParseResourceType(s string) ResourceType {
	su := strings.ToUpper(s)
	// look for direct match first
	if index := strIndex(su, resourceTypeNames[:]); index >= 0 {
		return ResourceType(index)
	}
	// Look if there's a pattern match
	for k, v := range resourceTypePatterns {
		m, err := regexp.MatchString(k, su)
		if err != nil {
			log.Printf("Resource type pattern '%s' match problem.  %v", k, err)
			continue
		}
		if m {
			return v
		}
	}
	return Unknown
}

func strIndex(s string, ss []string) int {
	for i, sz := range ss {
		if sz == s {
			return i
		}
	}
	return -1
}
