package main

import (
	"bytes"
	"testing"
)

func TestLoader(t *testing.T) {
	var tt = []struct {
		testName, templateName, csvName, expect string
		success                                 bool
	}{
		{"basic", "basic.tpl", "basic.csv", "Bond, James Bond", true},
		{"fail1", "missing.tpl", "basic.csv", "", false},
		{"fail2", "basic.tpl", "missing.csv", "", false},
		{"escape", "escape.tpl", "escape.csv", "Michael O&#39;Brian", true},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
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
		})
	}
}
