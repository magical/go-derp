package derp

import (
	"strings"
	"testing"
)

func newempty() *node         { return &node{ty: empty} }
func newnull() *node          { return &node{ty: null} }
func newterm(t int) *node     { return &node{ty: term, t: t} }
func newalt(a, b *node) *node { return &node{ty: alt, a: a, b: b} }
func newcat(a, b *node) *node { return &node{ty: cat, a: a, b: b} }

func TestDeriv(t *testing.T) {
	var tests = []struct {
		g    *node
		want nodeType
	}{
		{newempty(), null},
		{newnull(), null},
		{newterm('+'), empty},
		{newterm('-'), null},
	}
	for _, tt := range tests {
		g := deriv(tt.g, '+')
		if g.ty != tt.want {
			t.Errorf("deriv(%v, '+'): want %v, got %v", tt.g, tt.want, g.ty)
		}
	}
}

func TestNullable(t *testing.T) {
	var tests = []struct {
		want bool
		g    *node
	}{
		{true, newempty()},
		{false, newnull()},
		{false, newterm('+')},
		{true, newalt(newempty(), newempty())},
		{false, newalt(newnull(), newnull())},
		{true, newalt(newnull(), newempty())},
		{true, newalt(newempty(), newnull())},
	}
	for _, tt := range tests {
		got := nullable(tt.g)
		if got != tt.want {
			t.Errorf("nullable(%v): want %v, got %v", tt.g, tt.want, got)
		}
	}
}

func TestIsNull(t *testing.T) {
	if isnull(S().b) {
		t.Errorf("isnull(S().b): want %v, got %v", false, true)
	}
}

func TestPathological(t *testing.T) {
	t.Skip("exponential")
	goodString := strings.Repeat("1+", 100) + "1"
	badString := strings.Repeat("1+", 100)
	if _, ok := Match(S(), goodString); !ok {
		t.Errorf("failed to recognize goodString")
	}
	if _, ok := Match(S(), badString); ok {
		t.Errorf("failed to reject badString")
	}
}
