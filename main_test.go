package main

import (
	"bytes"
	"testing"
)

func TestLoader(t *testing.T) {
	var tt = []struct {
		templateName, csvName, expect string

		err error
	}{
		{"basic.tpl", "basic.csv", "1: 3\n", nil},
	}
	for _, tc := range tt {
		var buf bytes.Buffer
		err := run(
			"testfiles/"+tc.templateName,
			"testfiles/"+tc.csvName,
			&buf)
		if err != tc.err {
			t.Fatalf("expected err == %v; got %v", tc.err, err)
		}
		if tc.expect != buf.String() {
			t.Fatalf("expected result == %q; got %q", tc.expect, buf.String())
		}
	}
}
