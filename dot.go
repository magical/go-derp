package derp

// DOT-file generator for derp grammars.

import "io"
import "fmt"

var dotnum uint32

// Dot writes out a grammar as a DOT file.
func Dot(g *node, w io.Writer) {
	fmt.Fprintln(w, "digraph {")
	fmt.Fprintln(w, "\tcenter=true;")
	fmt.Fprintln(w, "\troot [shape=doublecircle];")
	fmt.Fprintln(w, "\troot -> g1;")
	walk++
	dotnum = 0
	dotwalk(g, w)
	fmt.Fprintln(w, "}")
}

func dotwalk(g *node, w io.Writer) {
	if g.walk == walk {
		return
	}
	g.walk = walk
	dotnum++
	g.name = fmt.Sprintf("g%d", dotnum)
	switch g.ty {
	case null:
		fmt.Fprintf(w, "\t%s [label=null];\n", g.name)
	case empty:
		fmt.Fprintf(w, "\t%s [label=empty];\n", g.name)
	case term:
		fmt.Fprintf(w, "\t%s [shape=record,label=\"term|%q\"];\n", g.name, g.t)
	case alt:
		dotwalk(g.a, w)
		dotwalk(g.b, w)
		fmt.Fprintf(w, "\t%s -> %s [label=a];\n", g.name, g.a.name)
		fmt.Fprintf(w, "\t%s -> %s [label=b];\n", g.name, g.b.name)
		fmt.Fprintf(w, "\t%s [label=alt];\n", g.name)
	case cat:
		dotwalk(g.a, w)
		dotwalk(g.b, w)
		fmt.Fprintf(w, "\t%s:L -> %s;\n", g.name, g.a.name)
		fmt.Fprintf(w, "\t%s:R -> %s;\n", g.name, g.b.name)
		fmt.Fprintf(w, "\t%s [shape=record,label=\"{cat|{<L>L|<R>R}}\"];\n", g.name)
	}
}
