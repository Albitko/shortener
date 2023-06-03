package lint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for os.Exit in main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	expr := func(x *ast.ExprStmt) {
		if call, ok := x.X.(*ast.CallExpr); ok {
			if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
				if i, ok := selector.X.(*ast.Ident); ok && i.Name == `os` {
					if selector.Sel.Name == `Exit` {
						pass.Reportf(selector.Pos(), "os.Exit exists in main body")
					}
				}
			}
		}
	}

	for _, file := range pass.Files {
		if file.Name.String() == "main" {
			ast.Inspect(file, func(node ast.Node) bool {
				switch x := node.(type) {
				case *ast.ExprStmt:
					expr(x)
				}
				return true
			})
		}

	}
	return nil, nil
}
