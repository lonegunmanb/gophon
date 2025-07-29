package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/lonegunmanb/gophon/pkg"
)

func main() {
	var (
		pkgPath    = flag.String("pkg", "", "Package path to scan (e.g., 'testharness' or '' for root)")
		basePkgUrl = flag.String("base", "", "Base package URL (e.g., 'github.com/lonegunmanb/gophon/pkg')")
		destDir    = flag.String("dest", "./index", "Destination directory for generated index files")
		help       = flag.Bool("help", false, "Show help message")
	)

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "gophon - Go Project Code Indexing Tool\n\n")
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		_, _ = fmt.Fprintf(os.Stderr, "\nExamples:\n")
		_, _ = fmt.Fprintf(os.Stderr, "  # Index the entire project\n")
		_, _ = fmt.Fprintf(os.Stderr, "  %s -base=github.com/lonegunmanb/gophon/pkg -dest=./output\n\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "  # Index a specific package\n")
		_, _ = fmt.Fprintf(os.Stderr, "  %s -pkg=testharness -base=github.com/lonegunmanb/gophon/pkg -dest=./output\n\n", os.Args[0])
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *basePkgUrl == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Error: -base flag is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Convert destination path to absolute path
	absDestDir, err := filepath.Abs(*destDir)
	if err != nil {
		log.Fatalf("Failed to resolve destination directory: %v", err)
	}

	fmt.Printf("Gophon Code Indexer\n")
	fmt.Printf("===================\n")
	fmt.Printf("Package path: %s\n", *pkgPath)
	fmt.Printf("Base URL: %s\n", *basePkgUrl)
	fmt.Printf("Destination: %s\n", absDestDir)
	fmt.Printf("\nGenerating index files...\n")

	// Call IndexSourceCode to generate index files
	err = pkg.IndexSourceCode(*pkgPath, *basePkgUrl, absDestDir)
	if err != nil {
		log.Fatalf("Failed to generate index files: %v", err)
	}

	fmt.Printf("✓ Index files generated successfully in: %s\n", absDestDir)
}
