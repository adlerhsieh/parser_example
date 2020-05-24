package extraline

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "extraline",
	Doc:  "catch extra line",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	funcNames := map[*ast.BlockStmt]*ast.Ident{}

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if funcDecl, ok := node.(*ast.FuncDecl); ok {
				blockStmt := funcDecl.Body
				funcNames[blockStmt] = funcDecl.Name
			}

			if blockStmt, ok := node.(*ast.BlockStmt); ok {
				// Find the line number of the beginning of a block statement.
				stmtStartingPosition := blockStmt.Pos()
				stmtLine := pass.Fset.Position(stmtStartingPosition).Line

				// Find the line number of the first statement in the block.
				firstStmt := blockStmt.List[0]
				firstStmtStartingPosition := firstStmt.Pos()
				firstStmtLine := pass.Fset.Position(firstStmtStartingPosition).Line

				// The difference should be one. Newlines exist when it is larger.
				if stmtLine+1 < firstStmtLine {
					// Retrieve the function name with the pointer key we saved earlier,
					// and print it.
					funcName := funcNames[blockStmt]
					fmt.Printf("Unnecessary newline at the beginning: %s\n", funcName)
				}
			}
			return true
		})
	}

	return nil, nil
}
