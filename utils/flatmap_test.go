package utils

import (
	"testing"
)

func TestNewFlatMap(t *testing.T) {
	testdata := makeTestMap()

	fm := NewFlatMap(testdata)
	if len(fm) != 6 {
		t.Errorf("Unexpected flatmap length, expected %d, found %d", 3, len(fm))
	}

	v, ok := fm["one"]
	if !ok {
		t.Errorf("Missing key %s", "one")
	}
	if v == nil {
		t.Errorf("Unexpected key one value is nil")
	} else if *v != "1" {
		t.Errorf("Unexpected key one value, expected %s, found %s", "1", *v)
	}

	v, ok = fm["two"]
	if !ok {
		t.Errorf("Missing key %s", "two")
	}
	if v == nil {
		t.Errorf("Unexpected key two value is nil")
	} else if *v != "two" {
		t.Errorf("Unexpected key two value, expected %s, found %s", "two", *v)
	}

	v, ok = fm["three"]
	if !ok {
		t.Errorf("Missing key %s", "three")
	}
	if v == nil {
		t.Errorf("Unexpected key three value is nil")
	} else if *v != "false" {
		t.Errorf("Unexpected key three value, expected %s, found %s", "false", *v)
	}

	v, ok = fm["four.a"]
	if !ok {
		t.Errorf("Missing key %s", "four.a")
	}
	if v == nil {
		t.Errorf("Unexpected key four.a value is nil")
	} else if *v != "true" {
		t.Errorf("Unexpected key four.a value, expected %s, found %s", "true", *v)
	}

	v, ok = fm["four.b"]
	if !ok {
		t.Errorf("Missing key %s", "four.b")
	}
	if v == nil {
		t.Errorf("Unexpected key four.b value is nil")
	} else if *v != "false" {
		t.Errorf("Unexpected key four.b value, expected %s, found %s", "false", *v)
	}

	v, ok = fm["four.c"]
	if !ok {
		t.Errorf("Missing key %s", "four.c")
	}
	if v == nil {
		t.Errorf("Unexpected key four.c value is nil")
	} else if *v != "-1" {
		t.Errorf("Unexpected key four.c value, expected %s, found %s", "-1", *v)
	}
}

func TestFlatMap_Expand(t *testing.T) {
	tm := makeTestMap()
	fm := NewFlatMap(tm)

	etm := fm.Expand()
	if len(etm) != len(tm) {
		t.Errorf("Unexpected expanded map length. expected %d, found %d", len(tm), len(etm))
	}
}

func makeTestMap() map[string]interface{} {
	return map[string]interface{}{
		"one":   1,
		"two":   "two",
		"three": false,
		"four": map[string]interface{}{
			"a": true,
			"b": false,
			"c": -1,
		},
	}
}
