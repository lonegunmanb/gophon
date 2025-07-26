// POC, DO NOT CHANGE IT OR USE IT, WILL BE DELETED IN THE FUTURE
package main

import (
	"log"
	"os"

	"github.com/gophon/pkg"
)

func main() {
	// 1. 检查命令行参数
	if len(os.Args) != 2 {
		log.Fatalf("Usage: go run main.go <path_to_package>")
	}
	pkgPath := os.Args[1]

	// 2. 使用新的扫描函数
	info, err := pkg.ScanPackage(pkgPath)
	if err != nil {
		log.Fatalf("Error scanning package: %v", err)
	}

	// 3. 显示扫描结果
	info.PrintDetailedInfo()
}
