package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func SaveAst2File(fileName string, fileSet *token.FileSet, node *ast.File) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// 格式化代码
	var buf bytes.Buffer
	err = format.Node(&buf, fileSet, node)
	if err != nil {
		return err
	}

	// 将结果写入文件
	_, err = file.WriteString(buf.String())
	return err

	return nil
}

func main() {
	srcFileName := "interface_definition.go"
	outFilename := "out.go"

	// 解析 Go 源码文件
	fileSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fileSet, srcFileName, nil, parser.ParseComments)
	if err != nil {
		log.Panicln(err.Error())
	}
	outFile, err := parser.ParseFile(fileSet, outFilename,
		fmt.Sprintf("package %s", astFile.Name.Name), parser.ParseComments)
	if err != nil {
		return
	}

	// 遍历 AST
	ast.Inspect(astFile, func(node ast.Node) bool {
		if genDecl, ok := node.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {

		}
		return true
	})

	// 保存 ast 至文件中
	if err := SaveAst2File(outFilename, fileSet, outFile); err != nil {
		log.Panicln(err.Error())
	}
}

// CreateZeroExpr 创建零值的 AST 表达式
func CreateZeroExpr(name string) ast.Expr {
	switch name {
	case "bool":
		return &ast.Ident{Name: "false"}
	case "int", "int8", "int16", "int32", "int64":
		return &ast.BasicLit{Kind: token.INT, Value: "0"}
	case "uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
		return &ast.BasicLit{Kind: token.INT, Value: "0"}
	case "float32", "float64":
		return &ast.BasicLit{Kind: token.FLOAT, Value: "0.0"}
	case "complex64", "complex128":
		return &ast.BasicLit{Kind: token.IMAG, Value: "0.0i"}
	case "string":
		return &ast.BasicLit{Kind: token.STRING, Value: `""`}
	default:
		return &ast.Ident{Name: "nil"}
	}
}
