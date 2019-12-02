package main

import (
	"fmt"
	"go/parser"
	"go/token"
)

func main() {
	fpath := "./example/example.go"
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fpath, nil, parser.ImportsOnly)
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Println(f.Imports)

	for _, i := range f.Imports {
		fmt.Println(*i.Path)
	}
}
