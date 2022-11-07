package pemtypes

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

type PEMType int

const (
	Unknown PEMType = iota
	PrivateKey
	PrivateKeyEncrypted
	Request
	Certificate
	RevocationList
	PublicKey
	Name
)

var PemTypeNames = [...]string{"", "PRIVATE KEY", "ENCRYPTED PRIVATE KEY", "CERTIFICATE REQUEST", "CERTIFICATE", "X509 CRL", "PUBLIC KEY", "DISTINGUISHED NAME"}

func (r PEMType) String() string {
	return PemTypeNames[r]
}

// alternative type names mapped into regular ones.
var resourceTypePatterns = map[string]PEMType{
	"KEY":                      PrivateKey,
	"CERT":                     Certificate,
	"CSR":                      Request,
	"REQUEST":                  Request,
	".*REVOKE.*LIST.*":         RevocationList,
	".*PRIVATE.*KEY.*":         PrivateKey,
	".*ENCRYPTED.*KEY.*":       PrivateKeyEncrypted,
	".*PUBLIC.*KEY.*":          PublicKey,
	".*CERTIFICATE.*REQUEST.*": Request,
	".*CRL.*":                  RevocationList,
}

func ParsePEMType(s string) PEMType {
	su := strings.ToUpper(s)
	// look for direct match first
	if index := strIndex(su, PemTypeNames[:]); index >= 0 {
		return PEMType(index)
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
