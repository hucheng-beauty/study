package main

import (
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"go/token"
)

func AddLatencyForFunction() {
	// 创建一个 FileSet，用于跟踪文件的位置信息
	fs := token.NewFileSet()

	// 使用装饰器解析 Go 源代码文件
	node, err := decorator.ParseFile(fs, "main.go", nil, 0)
	if err != nil {
		panic(err)
	}

	// 遍历 DST 语法树，修改其中的函数声明
	dst.Inspect(node, func(n dst.Node) bool {
		if funcDecl, ok := n.(*dst.FuncDecl); ok {
			// now := time.Now()
			startStmt := dst.AssignStmt{
				Lhs: []dst.Expr{&dst.Ident{Name: "now"}},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{
					&dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X:   &dst.Ident{Name: "time"},
							Sel: dst.NewIdent("Now"),
						},
						Args: nil,
					},
				},
			}
			// 在函数体开头添加一行代码
			funcDecl.Body.List = append([]dst.Stmt{&startStmt}, funcDecl.Body.List...)

			// latency := time.Since(now)
			endStmt := &dst.AssignStmt{
				Lhs: []dst.Expr{
					dst.NewIdent("latency"), // 左边变量 latency
				},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{
					&dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X:   dst.NewIdent("time"), // 调用 time.Since()
							Sel: dst.NewIdent("Since"),
						},
						Args: []dst.Expr{
							dst.NewIdent("now"), // 参数 now
						},
					},
				},
			}
			// 在函数体末尾添加一行代码
			funcDecl.Body.List = append(funcDecl.Body.List, endStmt)
		}
		return true
	})

	// 输出修改后的 Go 源代码到 stdout
	errPrint := decorator.Print(node)
	if errPrint != nil {
		panic(errPrint)
	}
}

func main() {
	AddLatencyForFunction()
}
