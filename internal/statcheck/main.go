package statcheck

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var ExitOnMain = &analysis.Analyzer{
	Name: "exitonmain",
	Doc:  "check for call os.Exit() on func main() in packet main",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	isExit := func(x *ast.CallExpr) bool {
		// @TODO
		return false
	}
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}
	for _, file := range pass.Files {
		// функцией ast.Inspect проходим по всем узлам AST
		inMain := false
		ast.Inspect(file, func(node ast.Node) bool {
			if inMain {
				switch x := node.(type) {
				case *ast.ExprStmt: // определение функции
					if call, ok := x.X.(*ast.CallExpr); ok {
						if isExit(call) {
							pass.Reportf(call.Pos(), "Unexpected call os.Exit() on func main() in packet main")
						}
					}
				}
			} else {
				switch x := node.(type) {
				case *ast.FuncDecl: // определение функции
					if x.Name.Name != "main" {
						return false
					}
					inMain = true
				}
			}
			return true
		})
	}
	return nil, nil
}
