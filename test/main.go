package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
)

/*
	_, file, line, _ := runtime.Caller(0)
	fmt.Printf("Current line is %s:%d\n", file, line)

	func main() {
	fmt.Println("hello world!")
	_, file, line, _ := runtime.Caller(0)
	fmt.Printf("Current line is %s:%d\n", file, line)
}
*/

type nodeInserter struct {
	fset *token.FileSet
}

func (ni nodeInserter) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return ni
	}

	switch stmt := node.(type) {
	case *ast.ImportSpec:
		fmt.Println(stmt.Path.Value)
	case *ast.FuncDecl:
		stmt.Body.List = append(stmt.Body.List, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{Name: "_"},
				&ast.Ident{Name: "file"},
				&ast.Ident{Name: "line"},
				&ast.Ident{Name: "_"},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "runtime",
						},
						Sel: &ast.Ident{Name: "Caller"},
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Value: "0",
						},
					},
				},
			},
		})

		stmt.Body.List = append(stmt.Body.List, &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "fmt",
					},
					Sel: &ast.Ident{
						Name: "Printf",
					},
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Value: "\"Current line is %s:%d\\n\"",
					},
					&ast.Ident{
						Name: "file",
					},
					&ast.Ident{
						Name: "line",
					},
				},
			},
		})
	}

	return nil
}

func main() {
	src := `
	package main
	
	func main() {
		fmt.Println("hello world!")
	}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		log.Fatal(err)
	}

	inserter := &nodeInserter{fset: fset}
	ast.Walk(inserter, file)

	if err := format.Node(os.Stdout, fset, file); err != nil {
		log.Fatal(err)
	}
}
