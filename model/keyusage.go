package model

import (
	"crypto/x509"
	"fmt"
	"github.com/eurozulu/pempal/tools"
	"strings"
)

type KeyUsage x509.KeyUsage

var keyUsageNames = []string{
	"UnknownKeyUsage",
	"KeyUsageDigitalSignature",
	"KeyUsageContentCommitment",
	"KeyUsageKeyEncipherment",
	"KeyUsageDataEncipherment",
	"KeyUsageKeyAgreement",
	"KeyUsageCertSign",
	"KeyUsageCRLSign",
	"KeyUsageEncipherOnly",
	"KeyUsageDecipherOnly",
}

func (k KeyUsage) String() string {
	var names []string
	ku := int(k)
	for i, name := range keyUsageNames {
		if i == 0 {
			continue
		}
		p := 2 ^ i
		if (p & ku) != p {
			continue
		}
		names = append(names, name)
	}
	if len(names) == 0 {
		return ""
	}
	return strings.Join(names, ",")
}

func (k KeyUsage) MarshalText() (text []byte, err error) {
	return []byte(k.String()), nil
}

func (k *KeyUsage) UnmarshalText(text []byte) error {
	names := tools.TrimSlice(strings.Split(string(text), "|"))
	ku := int(*k)
	for _, name := range names {
		u, err := ParseKeyUsage(name)
		if err != nil {
			return err
		}
		ku |= int(u)
	}
	*k = KeyUsage(ku)
	return nil
}

func ParseKeyUsage(s string) (KeyUsage, error) {
	if ss := strings.Split(s, ","); len(ss) > 1 {
		var k KeyUsage
		for _, name := range ss {
			ku, err := ParseKeyUsage(name)
			if err != nil {
				return 0, err
			}
			k |= ku
		}
		return k, nil
	}
	for i, n := range keyUsageNames {
		if strings.EqualFold(n, s) {
			return KeyUsage(i), nil
		}
	}
	return KeyUsage(0), fmt.Errorf("invalid key usage: %q", s)
}
