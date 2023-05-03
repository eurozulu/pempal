package model

import (
	"pempal/resources"
	"testing"
)

func TestResourceType_ParseResourceType(t *testing.T) {
	for i, ts := range []string{
		"unknown", "publickey", "privatekey", "certificaterequest", "certificate", "revocationlist",
	} {
		expect := resources.ResourceType(i)
		found := resources.ParseResourceType(ts)
		if expect != found {
			t.Errorf("Unexpected pemResource type parsed from '%s'.  Expected %s, found %s", ts, expect, found)
		}

	}
	// test empty
	rt := resources.ParseResourceType("")
	if rt != resources.Unknown {
		t.Errorf("Unexpected pemResource type parsed from empty string.  Expected %s, found %s", resources.Unknown, rt)
	}
}

func TestResourceType_ParsePemResourceType(t *testing.T) {
	for i, ts := range []string{
		"", "Public Key", "Private Key", "Certificate Request", "Certificate", "x509 crl",
	} {
		expect := resources.ResourceType(i)
		found := resources.ParsePEMType(ts)
		if expect != found {
			t.Errorf("Unexpected pemResource type parsed from PEM type '%s'.  Expected %s, found %s", ts, expect, found)
		}

	}
	// test empty
	rt := resources.ParsePEMType("")
	if rt != resources.Unknown {
		t.Errorf("Unexpected pemResource type pem parsed from empty string.  Expected %s, found %s", resources.Unknown, rt)
	}
}
