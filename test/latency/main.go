package main

import (
    "bytes"
    "fmt"
    "go/ast"
    "go/parser"
    "go/printer"
    "go/token"
    "log"
)

func main() {
    // 解析源代码并生成 AST
    fileSet := token.NewFileSet()
    node, err := parser.ParseFile(fileSet, "main.go",
        nil, parser.ParseComments)
    if err != nil {
        log.Panicln(err)
    }

    // 遍历 AST 并修改函数体
    ast.Inspect(node, func(n ast.Node) bool {
        // 找到函数声明
        if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "main" {
            // start := time.Now()
            startTime := &ast.AssignStmt{
                Lhs: []ast.Expr{ast.NewIdent("startTime")},
                Tok: token.DEFINE,
                Rhs: []ast.Expr{&ast.CallExpr{
                    Fun: &ast.SelectorExpr{
                        X:   ast.NewIdent("time"),
                        Sel: ast.NewIdent("Now"),
                    },
                }},
            }
            fn.Body.List = append([]ast.Stmt{startTime}, fn.Body.List...)

            // fmt.Println("latency: ", time.Since(startTime))
            endTime := &ast.ExprStmt{
                X: &ast.CallExpr{
                    Fun: &ast.SelectorExpr{
                        X:   ast.NewIdent("fmt"),
                        Sel: ast.NewIdent("Println"),
                    },
                    Args: []ast.Expr{
                        &ast.BasicLit{
                            Kind: token.STRING, Value: `"latency: "`},
                        &ast.CallExpr{
                            Fun: &ast.SelectorExpr{
                                X:   ast.NewIdent("time"),
                                Sel: ast.NewIdent("Since"),
                            },
                            Args: []ast.Expr{ast.NewIdent("startTime")},
                        },
                    },
                },
            }
            fn.Body.List = append(fn.Body.List, endTime)
        }
        return true
    })

    // 打印修改后的代码
    var buf bytes.Buffer
    err = printer.Fprint(&buf, fileSet, node)
    if err != nil {
        fmt.Println(err)
        return
    }

    // 输出修改后的代码
    fmt.Println(buf.String())

    // // 保存到文件（如果需要）
    // err = os.WriteFile("modified_code.go", buf.Bytes(), 0644)
    // if err != nil {
    //     fmt.Println("Error writing to file:", err)
    //     return
    // }
}
