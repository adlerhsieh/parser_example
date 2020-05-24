package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fpath := "./example/example.go"
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fpath, nil, parser.Mode(0))
	if err != nil {
		fmt.Println(err)
		return
	}

	funcNames := map[*ast.BlockStmt]*ast.Ident{}

	ast.Inspect(f, func(node ast.Node) bool {
		if funcDecl, ok := node.(*ast.FuncDecl); ok {
			blockStmt := funcDecl.Body
			funcNames[blockStmt] = funcDecl.Name
		}

		if blockStmt, ok := node.(*ast.BlockStmt); ok {
			// Find the line number of the beginning of a block statement.
			stmtStartingPosition := blockStmt.Pos()
			stmtLine := fset.Position(stmtStartingPosition).Line

			// Find the line number of the first statement in the block.
			firstStmt := blockStmt.List[0]
			firstStmtStartingPosition := firstStmt.Pos()
			firstStmtLine := fset.Position(firstStmtStartingPosition).Line

			// The difference should be one. Newlines exist when it is larger.
			if stmtLine+1 < firstStmtLine {
				// Retrieve the function name with the pointer key we saved earlier,
				// and print it.
				funcName := funcNames[blockStmt]
				fmt.Printf("Unnecessary newline at the beginning. %s() Line:%d\n", funcName, firstStmtLine)
			}
		}
		return true
	})
}
