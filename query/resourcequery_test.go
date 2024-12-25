package query

import (
	"github.com/eurozulu/pempal/model"
	"slices"
	"testing"
)

func TestResourceQuery_ColumnNames(t *testing.T) {
	query := ResourceQuery{
		Fields: nil,
		Types:  nil,
	}
	checkColumnNames(query.ColumnNames(), defaultFields[0], t)
	query.Types = []model.ResourceType{model.PrivateKey}
	checkColumnNames(query.ColumnNames(), defaultFields[model.PrivateKey], t)
	query.Types = []model.ResourceType{model.Certificate}
	checkColumnNames(query.ColumnNames(), defaultFields[model.Certificate], t)
	query.Types = []model.ResourceType{model.Certificate, model.PrivateKey}
	checkColumnNames(query.ColumnNames(), append(defaultFields[model.Certificate], "public-key-algorithm"), t)

	query.Types = nil
	query.Fields = []string{"public-key-algorithm"}
	checkColumnNames(query.ColumnNames(), []string{"public-key-algorithm"}, t)
	query.Types = []model.ResourceType{model.Certificate}
	checkColumnNames(query.ColumnNames(), []string{"public-key-algorithm"}, t)
	query.Fields = []string{"+public-key-algorithm"}
	checkColumnNames(query.ColumnNames(), append(defaultFields[model.Certificate], "public-key-algorithm"), t)

}

func TestResourceQuery_Query(t *testing.T) {
	query := ResourceQuery{}
	found := query.QueryAll("../testing")
	if len(found) != 6 {
		t.Errorf("found %d resources, expected %d", 6, len(found))
	}
	p1 := found[0]
	if len(p1) != 3 {
		t.Errorf("found %d properties, expected %d", 3, len(p1))
	}

	query.Types = []model.ResourceType{model.Certificate}
	found = query.QueryAll("../testing")
	if len(found) != 3 {
		t.Errorf("found %d resources, expected %d", 3, len(found))
	}
	for _, p := range found {
		tp := p["type"].(string)
		if tp != "Certificate" {
			t.Errorf("found resource type %s, expected %s", tp, "Certificate")
		}
	}
}

func checkColumnNames(found, expected []string, t *testing.T) {
	if len(found) != len(expected) {
		t.Errorf("Found %d columns, expected %d", len(found), len(expected))
	}
	for _, expect := range expected {
		if !slices.Contains(found, expect) {
			t.Errorf("Did not find expected column: %s", expect)
		}
	}
}
