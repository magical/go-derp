// "Parsing with Derivatives" for Go.
package derp

import "fmt"

type nodeType int
type node struct {
	ty   nodeType
	a, b *node
	t    int

	walk, walkempty uint32

	empty bool
	memo  *node
	name  string
}

var walk, walkempty uint32

// Null is the null set. It represents a failed parse.
// Empty is the empty string. It represents a successful parse.
// Term is a terminal production. It matches a token.
// Alt represents an alternative between two productions.
// Cat represents the concatenation of two productions.
const (
	null nodeType = iota
	empty
	term
	alt
	cat
)

// Deriv computes the derivative of a grammar with respect to a token.
//
// D null = null
// D empty = null
// D term(t) = if token == t then empty else null
// D alt(a, b) = alt(D a, D b)
// D cat(a, b) = if nullable(a) then alt(cat(D a, b), D b)
//                              else     cat(D a, b)
func deriv(g *node, token int) *node {
	walk++
	return _deriv(g, token)
}
func _deriv(g *node, token int) *node {
	if g.walk == walk && g.memo != nil {
		return g.memo
	}
	g.walk = walk
	switch g.ty {
	case null:
		return g
	case empty:
		return &node{ty: null}
	case term:
		if g.t == token {
			return &node{ty: empty}
		}
		return &node{ty: null}
	case alt:
		g.memo = &node{ty: alt}
		g.memo.a = _deriv(g.a, token)
		g.memo.b = _deriv(g.b, token)
		return g.memo
	case cat:
		if nullable(g.a) {
			g.memo = &node{ty: alt}
			g.memo.a = &node{ty: cat}
			g.memo.a.a = _deriv(g.a, token)
			g.memo.a.b = g.b
			g.memo.b = _deriv(g.b, token)
		} else {
			g.memo = &node{ty: cat}
			g.memo.a = _deriv(g.a, token)
			g.memo.b = g.b
		}
		return g.memo
	default:
		panic("unreachable")
	}
}

// Compact simplifies a grammar.
//
// alt(a, null) => a
// alt(null, b) => b
// cat(a, empty) => a
// cat(empty, b) => b
//
// TODO: combine with deriv?
func compact(g *node) *node {
	walk++
	return _compact(g)
}
func _compact(g *node) *node {
	if g.walk == walk && g.memo != nil {
		return g.memo
	}
	g.walk = walk
	switch g.ty {
	case null, empty, term:
		return g
	}
	if isnull(g) {
		g.memo = &node{ty: null}
		return g.memo
	}
	if isempty(g) {
		g.memo = &node{ty: empty}
		return g.memo
	}
	switch g.ty {
	case alt:
		if isnull(g.a) {
			g.memo = nil
			return _compact(g.b)
		}
		if isnull(g.b) {
			g.memo = nil
			return _compact(g.a)
		}
		g.memo = g
		g.a = _compact(g.a)
		g.b = _compact(g.b)
		return g
	case cat:
		if isempty(g.a) {
			g.memo = nil
			return _compact(g.b)
		}
		if isempty(g.b) {
			g.memo = nil
			return _compact(g.a)
		}
		g.memo = g
		g.a = _compact(g.a)
		g.b = _compact(g.b)
		return g
	default:
		panic("unreachable")
	}
}

// DerivEmpty computes the derivative of a grammar with respect to the empty string.
//
// D null = null
// D empty =
// D term(t) = null
// D alt(a, b) = set(D a, D b)
// D cat(a, b) =
func derivEmpty(grammar *node) *node { return nil }

// isnull reports whether a grammar is equivalent to the empty set.
func isnull(g *node) bool {
	walkempty++
	return _isnull(g)
}
func _isnull(g *node) bool {
	switch g.ty {
	case null:
		return debug("null", true)
	case empty:
		return debug("empty", false)
	case term:
		return debug("term", false)
	case alt:
		if g.walkempty == walkempty {
			return debug("assume alt", g.empty)
		}
		g.walkempty = walkempty
		g.empty = true
		g.empty = _isnull(g.a) && _isnull(g.b)
		return debug("alt", g.empty)
	case cat:
		if g.walkempty == walkempty {
			return debug("assume cat", g.empty)
		}
		g.walkempty = walkempty
		g.empty = true
		g.empty = _isnull(g.a) || _isnull(g.b)
		return debug("cat", g.empty)
	default:
		panic("unreachable")
	}
}

func debug(name string, b bool) bool {
	/*if b {
		println(name, "isnull")
	} else {
		println(name, "!isnull")
	}*/
	return b
}

// isempty reports whether a grammar is equivalent to the empty string.
func isempty(g *node) bool {
	walkempty++
	return _isempty(g)
}
func _isempty(g *node) bool {
	switch g.ty {
	case null:
		return false
	case empty:
		return true
	case term:
		return false
	case alt:
		if g.walkempty == walkempty {
			return g.empty
		}
		g.walkempty = walkempty
		g.empty = true
		g.empty = _isempty(g.a) && _isempty(g.b)
		return g.empty
	case cat:
		if g.walkempty == walkempty {
			return g.empty
		}
		g.walkempty = walkempty
		g.empty = true
		g.empty = _isempty(g.a) && _isempty(g.b)
		return g.empty
	default:
		panic("unreachable")
	}
}

// Nullable reports whether a grammar contains the empty string.
//
// E null = false
// E empty = true
// E term(t) = false
// E alt(a, b) = E a || E b
// E cat(a, b) = E a && E b
func nullable(g *node) bool {
	walkempty++
	return _nullable(g)
}
func _nullable(g *node) bool {
	switch g.ty {
	case null:
		return false
	case empty:
		return true
	case term:
		return false
	case alt:
		if g.walkempty == walkempty {
			return g.empty
		}
		g.walkempty = walkempty
		g.empty = false
		g.empty = _nullable(g.a) || _nullable(g.b)
		return g.empty
	case cat:
		if g.walkempty == walkempty {
			return g.empty
		}
		g.walkempty = walkempty
		g.empty = false
		g.empty = _nullable(g.a) && _nullable(g.b)
		return g.empty
	default:
		panic("unreachable")
	}
}

func size(g *node) int {
	walk++
	return _size(g)
}
func _size(g *node) (n int) {
loop:
	if g.walk == walk {
		return n
	}
	g.walk = walk
	n++
	switch g.ty {
	case alt, cat:
		n += _size(g.a)
		g = g.b
		goto loop
	default:
	}
	return n
}

// S constructs the grammar S = 1 | S + S.
func S() *node {
	s := &node{ty: alt}
	s.a = &node{ty: term, t: '1'}
	s.b = &node{ty: cat,
		a: s,
		//a: &node{ty: term, t: '1'},
		b: &node{ty: cat,
			a: &node{ty: term, t: '+'},
			b: s}}
	return s
}

// Match derives a grammar with respect to a string. It returns the new grammar
// and whether the string was a valid string in the original grammar.
func Match(g *node, s string) (*node, bool) {
	//isnull(g.a)
	//println("---")
	//isnull(g.b)
	g = compact(g)
	//return g, false
	fmt.Println(size(g), g)
	gg := g
	_ = gg
	for _, r := range s {
		g = deriv(g, int(r))
		gg = g
		sz := size(g)
		g = compact(g)
		fmt.Println(sz, size(g), g.ty)
	}
	return gg, nullable(g)
}
