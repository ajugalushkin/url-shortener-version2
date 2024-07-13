package osexitcheck

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer структура анализатора
var Analyzer = &analysis.Analyzer{
	Name: "osexitlint",
	Doc:  "prohibiting the use of a direct call to os.Exit",
	Run:  run,
}

// run функция анализатора
func run(pass *analysis.Pass) (interface{}, error) {
	isMainPackage := func(x *ast.File) bool {
		return x.Name.Name == "main"
	}

	isMainFunc := func(x *ast.FuncDecl) bool {
		return x.Name.Name == "main"
	}

	isOsExit := func(x *ast.SelectorExpr, isMain bool) bool {
		if !isMain || x.X == nil {
			return false
		}
		ident, ok := x.X.(*ast.Ident)
		if !ok {
			return false
		}
		if ident.Name == "os" && x.Sel.Name == "Exit" {
			pass.Reportf(ident.NamePos, "os.Exit called in main func in main package")
			return true
		}
		return false
	}

	i := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.SelectorExpr)(nil),
	}

	mainInspecting := false

	i.Preorder(nodeFilter, func(n ast.Node) {
		switch x := n.(type) {
		case *ast.File:
			// если пакет не main выходим
			if !isMainPackage(x) {
				return
			}
		case *ast.FuncDecl:
			// проверяем что функция main и что до этого обрабатывали пакет main
			fn := isMainFunc(x)
			if mainInspecting && !fn {
				mainInspecting = false
				return
			}
			mainInspecting = fn
		case *ast.SelectorExpr:
			// проверяем на вызов os.Exit
			if isOsExit(x, mainInspecting) {
				return
			}
		}
	})

	return nil, nil
}
