package passman

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	handlesLength(t)
	handlesNosym(t)
}

func handlesLength(t *testing.T) {
	lengths := []int{10, 0, 3, 4, 5}
	for _, l := range lengths {
		g := Generate(l, false)
		if len(g) != l {
			t.Errorf("Expected length %d, got %d", l, len(g))
		}
	}
}

func handlesNosym(t *testing.T) {
	symbols := "[]{}!@#$%^&*()_+-=;'.,<>';:\""
	g := Generate(5, true)
	if strings.ContainsAny(g, symbols) {
		t.Errorf("Expected generated password not to contain symbols, got %s", g)
	}
}
