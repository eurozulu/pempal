package templates

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

var givenNames = []string{}

var expectNames = []string{
	"certificate",
	"certificate/allkeyusage",
	"certificate/cert",
	"certificate/intermediatecertificate",
	"certificate/rootcacertificate",
	"certificate/servercertificate",
	"certificaterequest",
	"certificaterequest/csr",
	"privatekey",
	"privatekey/ed25519key",
	"privatekey/key",
	"privatekey/rsakey",
	"privatekey/serverkey",
	"revokationlist",
	"revokationlist/crl",
}
var expectNamesMore = []string{
	"org",
	"certificate/org",
	"org/orgservercertificate",
	"issuedByInter",
	"certificate/servercertificate/testservercertificate",
}

func TestTemplateLib_GetTemplateNames(t *testing.T) {
	tl := NewTemplateLib()
	names := tl.GetTemplateNames()
	l := len(expectNames)

	if len(names) != l {
		t.Errorf("unexpected number of names found.  Expected %d, got %d", l, len(names))
	}
	for _, name := range names {
		if !slices.Contains(expectNames, name) {
			t.Errorf("unexpected name found., got %s", name)
		}
	}
	tl = NewTemplateLib("../testing/templates")
	names = tl.GetTemplateNames()
	l = len(expectNames) + len(expectNamesMore)

	if len(names) != l {
		t.Errorf("unexpected number of names found.  Expected %d, got %d", l, len(names))
	}
}

func TestTemplateLib_ResolveName(t *testing.T) {
	tl := NewTemplateLib("../testing/templates").(*templateLib)
	n, err := tl.ResolveName("unknownname")
	if err == nil {
		t.Errorf("expected error for unknown template name, got none")
	}
	if n != "" {
		t.Errorf("expected empty string for unknown template name, got %s", n)
	}

	n, err = tl.ResolveName("samename")
	if err == nil {
		t.Errorf("expected error for ambigious template name rootca, got none")
	}
	if !strings.Contains(err.Error(), "is ambiguous") {
		t.Errorf("expected error for ambiguous template name rootca, got %q", err.Error())
	}

	for i, name := range givenNames {
		n, err := tl.ResolveName(name)
		if err != nil {
			t.Errorf("Unexpected error while resolving name %s: %s", name, err)
			continue
		}
		if n != expectNames[i] {
			t.Errorf("Expected name %s to resolve to %s, got %s", name, expectNames[i], n)
		}
	}
}

func TestTemplateLib_GetTemplates(t *testing.T) {
	//tl := NewTemplateLib("../testing/templates").(*templateLib)
	tl := NewTemplateLib().(*templateLib)

	for i, name := range givenNames {
		tps, err := tl.GetTemplates(name)
		if err != nil {
			t.Errorf("Unexpected error while resolving name %s: %s", name, err)
			continue
		}
		expN := strings.Split(expectNames[i], string(filepath.Separator))
		if len(expN) != len(tps) {
			t.Errorf("unexpected number of templates for %q.  Expected %d, got %d", name, len(expN), len(tps))
			continue
		}
		for tpi, tp := range tps {
			expectName := strings.Join(expN[:tpi+1], string(filepath.Separator))
			if tp.Name() != expectName {
				t.Errorf("unexpected name for %q.  Expected %s, got %s", name, expectName, tp.Name())
			}
		}
	}
}

type testBaseTemplateOne struct {
	name string
}

func (t testBaseTemplateOne) Name() string {
	return t.name
}

type testBaseTemplateTwo struct {
	name string
}

func (t testBaseTemplateTwo) Name() string {
	return t.name
}

var testBaseTemplates = []Template{
	&testBaseTemplateOne{name: "t1"},
	&testBaseTemplateTwo{name: "t2"},
}
