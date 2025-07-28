package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/lonegunmanb/gophon/pkg"
)

func main() {
	var pkgPath = flag.String("pkg", "", "Package path to scan (e.g., 'test-harness')")
	var basePkgUrl = flag.String("base", "", "Base package URL (e.g., 'github.com/lonegunmanb/gophon/pkg')")
	flag.Parse()

	if *basePkgUrl == "" {
		fmt.Fprintf(os.Stderr, "Error: -base flag is required\n")
		flag.Usage()
		os.Exit(1)
	}

	// Run ScanSinglePackage
	result, err := pkg.ScanSinglePackage(*pkgPath, *basePkgUrl)
	if err != nil {
		log.Fatalf("Failed to scan package: %v", err)
	}

	// Print all type names and positions
	fmt.Printf("Found %d types:\n\n", len(result.Types))
	for _, typeInfo := range result.Types {
		fmt.Printf("Type: %s\n", typeInfo.Name)
		fmt.Printf("  SourceCode: %s\n", typeInfo.String())
		fmt.Println()
	}
	for _, function := range result.Functions {
		fmt.Printf("Function: %s\n", function.Name)
		fmt.Printf("  SourceCode: %s\n", function.String())
		fmt.Println()
	}
}
