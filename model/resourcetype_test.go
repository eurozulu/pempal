package model

import (
	"testing"
)

func TestResourceType_ParseResourceType(t *testing.T) {
	for i, ts := range []string{
		"unknown", "publickey", "privatekey", "certificaterequest", "certificate", "revocationlist",
	} {
		expect := ResourceType(i)
		found := ParseResourceType(ts)
		if expect != found {
			t.Errorf("Unexpected pemResource type parsed from '%s'.  Expected %s, found %s", ts, expect, found)
		}

	}
	// test empty
	rt := ParseResourceType("")
	if rt != Unknown {
		t.Errorf("Unexpected pemResource type parsed from empty string.  Expected %s, found %s", Unknown, rt)
	}
}

func TestResourceType_ParsePemResourceType(t *testing.T) {
	for i, ts := range []string{
		"", "Public Key", "Private Key", "Certificate Request", "Certificate", "x509 crl",
	} {
		expect := ResourceType(i)
		found := ParsePEMType(ts)
		if expect != found {
			t.Errorf("Unexpected pemResource type parsed from PEM type '%s'.  Expected %s, found %s", ts, expect, found)
		}

	}
	// test empty
	rt := ParsePEMType("")
	if rt != Unknown {
		t.Errorf("Unexpected pemResource type pem parsed from empty string.  Expected %s, found %s", Unknown, rt)
	}
}
