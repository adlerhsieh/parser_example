package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	txt := `
		package example

		import "fmt"

		func Foo() {

			fmt.Println("Hello Foo")
		}
	`

	// fpath := "./example/example.go"
	fset := token.NewFileSet()
	// f, err := parser.ParseFile(fset, fpath, "", parser.Mode(0))
	f, err := parser.ParseFile(fset, "", txt, parser.Mode(0))
	if err != nil {
		fmt.Println(err)
		return
	}

	// ast.Inspect(f, func(n ast.Node) bool {
	// 	fmt.Println(n)
	// 	// Find Return Statements
	// 	ret, ok := n.(*ast.ReturnStmt)
	// 	if ok {
	// 		fmt.Printf("return statement found on line %d:\n\t", fset.Position(ret.Pos()).Line)
	// 		printer.Fprint(os.Stdout, fset, ret)
	// 		return true
	// 	}
	// 	return true
	// })

	// imports := f.Imports

	// for _, i := range imports {
	// path := *i.Path
	// fmt.Println(path.Value)
	// }

	v := visitor{
		fset:      fset,
		funcNames: make(map[*ast.BlockStmt]*ast.Ident),
	}

	ast.Walk(v, f)
}

type visitor struct {
	// We need this field to save the fileset as a reference for line numbers.
	fset *token.FileSet
	// When a function is detected as NNL, function name is retrieved here.
	funcNames map[*ast.BlockStmt]*ast.Ident
}

func (v visitor) Visit(node ast.Node) ast.Visitor {
	// nil node is skipped as it is irrelevant to our goal
	if node == nil {
		return nil
	}

	// fmt.Printf("%T\n", node)
	// if i, ok := node.(*ast.Ident); ok {
	// 	fmt.Println(i.Name)
	// }

	// Once we find a function, we save the function name in
	// a map using its body statement as a pointer key.
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		blockStmt := funcDecl.Body
		v.funcNames[blockStmt] = funcDecl.Name
	}

	if blockStmt, ok := node.(*ast.BlockStmt); ok {
		// Find the line number of the beginning of a block statement.
		stmtStartingPosition := blockStmt.Pos()
		stmtLine := v.fset.Position(stmtStartingPosition).Line

		// Find the line number of the first statement in the block.
		firstStmt := blockStmt.List[0]
		firstStmtStartingPosition := firstStmt.Pos()
		firstStmtLine := v.fset.Position(firstStmtStartingPosition).Line

		// The difference should be one. Newlines exist when it is larger.
		if stmtLine+1 < firstStmtLine {
			// Retrieve the function name with the pointer key we saved earlier,
			// and print it.
			funcName := v.funcNames[blockStmt]
			fmt.Printf("Unnecessary newline at the beginning: %s\n", funcName)
		}
	}

	// if stmt, ok := node.(*ast.ExprStmt); ok {
	// 	fmt.Println(stmt.X)
	// }

	return v
}
