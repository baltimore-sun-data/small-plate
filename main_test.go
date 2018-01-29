package main

import (
	"bytes"
	"testing"
)

func TestLoader(t *testing.T) {
	var tt = []struct {
		templateName, csvName, expect string
		success                       bool
	}{
		{"basic.tpl", "basic.csv", "Bond, James Bond", true},
		{"missing.tpl", "basic.csv", "", false},
		{"basic.tpl", "missing.csv", "", false},
	}
	for _, tc := range tt {
		var buf bytes.Buffer
		err := run(
			"testfiles/"+tc.templateName,
			"testfiles/"+tc.csvName,
			&buf)
		if (err == nil) != tc.success {
			t.Fatalf("expected success == %v; got error %v", tc.success, err)
		}
		if tc.expect != buf.String() {
			t.Fatalf("expected result == %q; got %q", tc.expect, buf.String())
		}
	}
}
