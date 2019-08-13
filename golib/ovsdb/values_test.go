package ovsdb

import "testing"

func TestUnmarshalJson(t *testing.T) {
	validValus := []string{
		`null`,
		`10`,
		`true`,
		`"abcdefg"`,
		`["uuid", "xsaxsaxsa"]`,
		`["set", []]`,
		`["set", [1, true, "xsaxsa", null]]`,
		`["map", []]`,
		`["map", [["key", 1], ["key2", "value2"], ["key3", null]]]`,
		`["map", [["key1", ["set", [1, true]]]]]`,
	}
	for _, s := range validValus {
		if v, err := UnmarshalJson([]byte(s)); err != nil {
			t.Fatal(err)
		} else {
			t.Log(v.Value())
		}
	}
}
