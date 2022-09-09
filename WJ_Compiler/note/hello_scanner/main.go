package main

import (
	"example/syntax"
	"fmt"
	"strings"
)

func main() {
    fmt.Println("HELLO WORLD!")
	const src = "if (+foo\t+=..123/***/0.9_0e-0i'a'`raw`\"string\"..f;//$"
	tokens := []token{_If, _Lparen, _Operator, _Name, _AssignOp, _Dot, _Literal, _Literal, _Literal, _Literal, _Literal, _Dot, _Dot, _Name, _Semi, _EOF}
    
    got :=syntax.ChanType
	got.init(strings.NewReader(src), errh, 0)
	for _, want := range tokens {
		got.next()
		if got.tok != want {
			continue
		}
	}
}
