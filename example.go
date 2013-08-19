// +build ignore

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/magical/go-derp"
)

func main() {
	g, ok := derp.Match(derp.S(), os.Args[1])
	fmt.Println(ok)

	f, err := os.Create("out.dot")
	if err != nil {
		fmt.Println(err)
		return
	}
	w := bufio.NewWriter(f)
	derp.Dot(g, w)
	w.Flush()
	f.Close()
}
