package pkg

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"sync"
)

// FileInfo contains information about a single Go file
type FileInfo struct {
	*ast.File
	FileName string
	FilePath string
	Package  string
	content  *string      // cached file content
	mu       sync.RWMutex // mutex for thread-safe cache access
}

// Imports returns the import statements of the file
func (f *FileInfo) Imports() string {
	// Get the file content
	content := f.String()
	if content == "" {
		return ""
	}

	// Parse the Go source code
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return ""
	}

	// Extract import declarations
	var imports []string
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			// Get the import declaration as string
			start := fset.Position(genDecl.Pos()).Offset
			end := fset.Position(genDecl.End()).Offset
			if start >= 0 && end <= len(content) && end > start {
				importStr := content[start:end]
				imports = append(imports, strings.TrimSpace(importStr))
			}
		}
	}

	// Join all import declarations
	return strings.Join(imports, "\n")
}

func (f *FileInfo) PackagePath() string {
	return f.Package
}

// String reads and returns the content of the file
func (f *FileInfo) String() string {
	// Check cache with read lock
	f.mu.RLock()
	if f.content != nil {
		cached := *f.content
		f.mu.RUnlock()
		return cached
	}
	f.mu.RUnlock()

	// Acquire write lock for file reading and caching
	f.mu.Lock()
	defer f.mu.Unlock()

	// Double-check pattern: another goroutine might have cached it
	if f.content != nil {
		return *f.content
	}

	// Read file content
	contentBytes, err := os.ReadFile(f.FileName)
	if err != nil {
		return ""
	}

	// Cache the content
	contentStr := string(contentBytes)
	f.content = &contentStr

	return contentStr
}
