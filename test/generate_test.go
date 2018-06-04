package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

var symbols = "[]{}!@#$%^&*()_+-=;'.,<>';:\""

var gentests = []struct {
	in  []string
	len int
}{
	{[]string{"generate", "--length=23"}, 23},
	{[]string{"generate", "--length=5"}, 5},
	{[]string{"generate", "-l", "15"}, 15},
}

var symtests = []struct {
	in []string
}{
	{[]string{"generate", "--nosym"}},
	{[]string{"generate", "-n"}},
}

func TestGenerate(t *testing.T) {
	for _, tt := range gentests {
		gen := new(test)
		var b bytes.Buffer
		gen.init(os.Stdin, &b, tt.in...)
		err := gen.run()
		if err != nil {
			t.Errorf("could not run generate command\n%s", err)
		}
		res := strings.TrimRight(b.String(), "\n")
		if len(res) != tt.len {
			t.Errorf("expected len %d, got %d", tt.len, len(res))
		}
	}

	for _, tt := range symtests {
		gen := new(test)
		var b bytes.Buffer
		gen.init(os.Stdin, &b, tt.in...)
		err := gen.run()
		if err != nil {
			t.Errorf("could not run generate command\n%s", err)
		}
		res := strings.TrimRight(b.String(), "\n")
		if strings.ContainsAny(res, symbols) {
			t.Error("expected no symbols from `generate -n`")
		}
	}
}
