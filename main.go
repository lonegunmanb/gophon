// POC, DO NOT CHANGE IT OR USE IT, WILL BE DELETED IN THE FUTURE
package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/packages"
)

func main() {
	// 1. 检查命令行参数
	if len(os.Args) != 2 {
		log.Fatalf("Usage: go run method_finder.go <path_to_package>")
	}
	pkgPath := os.Args[1]

	// 2. 配置加载模式
	// 我们需要语法树(Syntax)和文件信息(Files)来定位代码
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax,
	}

	// 3. 加载指定的包
	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil {
		log.Fatalf("Error loading package %s: %v", pkgPath, err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	// 4. 对加载的包进行分析
	for _, pkg := range pkgs {
		fmt.Printf("--- Analyzing Package: %s ---\n\n", pkg.ID)
		analyzePackage(pkg)
	}
}

// analyzePackage 遍历包的 AST 并输出所需信息
func analyzePackage(pkg *packages.Package) {
	// fset 用于将代码位置 (token.Pos) 转换为可读的文件名和行列号
	fset := pkg.Fset

	// 遍历包中的所有 Go 源文件
	for _, file := range pkg.Syntax {
		fmt.Printf("File: %s\n", fset.File(file.Pos()).Name())

		// 使用 ast.Inspect 遍历文件中的所有节点
		ast.Inspect(file, func(n ast.Node) bool {
			// 根据节点类型进行处理
			switch decl := n.(type) {

			// case 1: 通用声明 (寻找类型定义)
			case *ast.GenDecl:
				if decl.Tok == token.TYPE {
					for _, spec := range decl.Specs {
						// 类型断言为 TypeSpec
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							position := fset.Position(typeSpec.Pos())
							fmt.Printf("  [Type]   Name: %-30s Location: %s\n", typeSpec.Name.Name, position)
						}
					}
				}

			// case 2: 函数/方法声明 (只寻找方法)
			case *ast.FuncDecl:
				// 这是关键：decl.Recv 不为 nil 才是一个方法
				if decl.Recv != nil && decl.Recv.List != nil {
					position := fset.Position(decl.Pos())

					// 构造接收者名称，例如 (s *MyStruct)
					var receiverType string
					// 检查接收者是否为指针类型
					if star, ok := decl.Recv.List[0].Type.(*ast.StarExpr); ok {
						receiverType = fmt.Sprintf("*%s", star.X)
					} else {
						// 否则为值类型
						receiverType = fmt.Sprintf("%s", decl.Recv.List[0].Type)
					}

					methodName := fmt.Sprintf("(%s) %s", receiverType, decl.Name.Name)
					fmt.Printf("  [Method] Name: %-30s Location: %s\n", methodName, position)
				}
			}
			return true // 继续遍历
		})
		fmt.Println("---------------------------------")
	}
}
