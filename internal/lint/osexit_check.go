package lint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OsExitAnalyzer - analyzes the presence of os.Exit in main func.
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for os.Exit in main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	findOs := func(s *ast.SelectorExpr) {
		if i, ok := s.X.(*ast.Ident); ok && i.Name == `os` {
			if s.Sel.Name == `Exit` {
				pass.Reportf(s.Pos(), "os.Exit exists in main body")
			}
		}
	}

	expr := func(x *ast.ExprStmt) {
		if call, ok := x.X.(*ast.CallExpr); ok {
			if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
				findOs(selector)
			}
		}
	}

	for _, file := range pass.Files {
		if file.Name.String() == "main" {
			ast.Inspect(file, func(node ast.Node) bool {
				if v, ok := node.(*ast.ExprStmt); ok {
					expr(v)
				}
				return true
			})
		}
	}
	return nil, nil
}
