package pkg

import (
	"fmt"
	"github.com/spf13/afero"
	"go/ast"
	"go/token"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"
)

var sourceFs = afero.NewOsFs()

// ProgressInfo represents progress information during package scanning
type ProgressInfo struct {
	Completed  int     // Number of packages completed
	Total      int     // Total number of packages discovered so far
	Current    string  // Currently processing package path
	Percentage float64 // Completion percentage (completed/total * 100)
}

// ScanSinglePackage scans the specified package and returns comprehensive information
func ScanSinglePackage(pkgPath, basePkgUrl string) (*PackageInfo, error) {
	// Use relative path for packages.Load to work with local filesystem
	var loadPath string
	if pkgPath == "" {
		loadPath = "."
	} else {
		loadPath = "./" + pkgPath
	}

	cfg := &packages.Config{
		Mode: packages.NeedFiles | packages.NeedName | packages.NeedImports | packages.NeedTypes | packages.NeedSyntax,
	}

	pkgs, err := packages.Load(cfg, loadPath)
	if err != nil {
		return nil, err
	}

	if len(pkgs) == 0 {
		return &PackageInfo{}, nil
	}

	pkg := pkgs[0]

	// Determine the actual package path based on the declared package name
	// Split the directory path and replace the last part with the actual package name
	pathParts := []string{basePkgUrl}
	if pkgPath != "" {
		sep := string(filepath.Separator)
		parts := strings.Split(pkgPath, sep)
		if !strings.Contains(pkgPath, sep) && strings.Contains(pkgPath, "/") {
			parts = strings.Split(pkgPath, "/")
		}
		if len(parts) > 1 {
			pathParts = append(pathParts, parts[:len(parts)-1]...)
		}
		pathParts = append(pathParts, pkg.Name)
	} else {
		pathParts = append(pathParts, pkg.Name)
	}
	actualPkgPath := strings.Join(pathParts, "/")

	var files []*FileInfo
	var constants []*ConstantInfo
	var variables []*VariableInfo
	var types []*TypeInfo
	var functions []*FunctionInfo

	// Extract constants, variables, and types from AST
	for _, file := range pkg.Syntax {
		if file == nil {
			continue
		}

		fileName := pkg.Fset.Position(file.Pos()).Filename
		fileInfo := &FileInfo{
			File:     file,
			FileName: fileName,
			Package:  actualPkgPath,
		}
		files = append(files, fileInfo)

		// Walk through declarations to find constants, variables, types, and functions
		for _, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok {
				switch genDecl.Tok {
				case token.CONST:
					constants = append(constants, extractDeclarations(actualPkgPath, genDecl, pkg, fileInfo, func(name string, pkgPath string, rangeInfo *Range) *ConstantInfo {
						return &ConstantInfo{
							GenDecl: genDecl,
							Name:    name,
							Range:   rangeInfo,
						}
					})...)
				case token.VAR:
					variables = append(variables, extractDeclarations(actualPkgPath, genDecl, pkg, fileInfo, func(name string, pkgPath string, rangeInfo *Range) *VariableInfo {
						return &VariableInfo{
							GenDecl: genDecl,
							Name:    name,
							Range:   rangeInfo,
						}
					})...)
				case token.TYPE:
					types = append(types, extractTypeDeclarations(genDecl, pkg, fileInfo)...)
				}
			} else if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				functions = append(functions, extractFunctionDeclarations(funcDecl, pkg, fileInfo)...)
			}
		}
	}

	return &PackageInfo{
		Files:     files,
		Constants: constants,
		Variables: variables,
		Types:     types,
		Functions: functions,
	}, nil
}

// Generic function to extract declarations from AST
func extractDeclarations[T any](pkgPath string, genDecl *ast.GenDecl, pkg *packages.Package, fileInfo *FileInfo, createFunc func(name string, pkgPath string, rangeInfo *Range) *T) []*T {
	var results []*T
	for _, spec := range genDecl.Specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			for _, name := range valueSpec.Names {
				// Skip blank identifier variables
				if name.Name == "_" {
					continue
				}

				// Get line numbers for the declaration
				startPos := pkg.Fset.Position(spec.Pos())
				endPos := pkg.Fset.Position(spec.End())

				rangeInfo := &Range{
					FileInfo:  fileInfo,
					StartLine: startPos.Line,
					EndLine:   endPos.Line,
				}

				result := createFunc(name.Name, pkgPath, rangeInfo)
				results = append(results, result)
			}
		}
	}
	return results
}

// Extract type declarations from AST
func extractTypeDeclarations(genDecl *ast.GenDecl, pkg *packages.Package, fileInfo *FileInfo) []*TypeInfo {
	var results []*TypeInfo
	for _, spec := range genDecl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok {
			// Get line numbers for the type declaration
			startPos := pkg.Fset.Position(typeSpec.Pos())
			endPos := pkg.Fset.Position(typeSpec.End())

			rangeInfo := &Range{
				FileInfo:  fileInfo,
				StartLine: startPos.Line,
				EndLine:   endPos.Line,
			}

			results = append(results, &TypeInfo{
				Name:    typeSpec.Name.Name,
				Range:   rangeInfo,
				GenDecl: genDecl,
			})
		}
	}
	return results
}

// Extract function declarations from AST
func extractFunctionDeclarations(funcDecl *ast.FuncDecl, pkg *packages.Package, fileInfo *FileInfo) []*FunctionInfo {
	var results []*FunctionInfo

	// Get line numbers for the function declaration
	startPos := pkg.Fset.Position(funcDecl.Pos())
	endPos := pkg.Fset.Position(funcDecl.End())

	rangeInfo := &Range{
		FileInfo:  fileInfo,
		StartLine: startPos.Line,
		EndLine:   endPos.Line,
	}

	// Determine receiver type (empty for functions, populated for methods)
	receiverType := ""
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
		// Extract receiver type
		if starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
			// Pointer receiver like *Service
			if ident, ok := starExpr.X.(*ast.Ident); ok {
				receiverType = "*" + ident.Name
			}
		} else if ident, ok := funcDecl.Recv.List[0].Type.(*ast.Ident); ok {
			// Value receiver like Service
			receiverType = ident.Name
		}
	}

	results = append(results, &FunctionInfo{
		Range:        rangeInfo,
		FuncDecl:     funcDecl,
		Name:         funcDecl.Name.Name,
		ReceiverType: receiverType,
	})

	return results
}

// ScanPackagesRecursively recursively scans all packages starting from the specified path
// and invokes the callback function for each package found. It uses afero.Fs for file system operations
// to enable easy testing with mocked file systems.
// Parameters:
//   - fs: The afero filesystem to use for file operations
//   - pkgPath: The relative package path to start scanning from (e.g., "pkg/utils")
//   - basePkgUrl: The base package URL/module path (e.g., "github.com/user/project")
//   - callback: Function called for each package, receives *PackageInfo and full package URL
//   - progressCallback: Optional callback for progress updates, receives ProgressInfo
func ScanPackagesRecursively(pkgPath, basePkgUrl string, callback func(*PackageInfo, string), progressCallback func(ProgressInfo)) error {
	// First, discover all packages to get accurate total count
	allPackages := findSubPackages(pkgPath)
	if pkgPath != "" || len(allPackages) == 0 {
		// Include the root package if we're scanning from a specific path or if no sub-packages found
		allPackages = append([]string{pkgPath}, allPackages...)
	}

	var completedWork int
	totalDiscovered := len(allPackages)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Helper function to report progress
	reportProgress := func(current string) {
		mu.Lock()
		defer mu.Unlock()

		var percentage float64
		if totalDiscovered > 0 {
			percentage = float64(completedWork) / float64(totalDiscovered) * 100.0
		}

		if progressCallback != nil {
			progressCallback(ProgressInfo{
				Completed:  completedWork,
				Total:      totalDiscovered,
				Current:    current,
				Percentage: percentage,
			})
		}
	}

	// Create a channel for work distribution
	workChan := make(chan string, len(allPackages))

	// Create error channel to collect errors from workers
	errChan := make(chan error, len(allPackages))

	// Determine number of workers (limited by CPU count)
	numWorkers := runtime.NumCPU()
	if numWorkers > len(allPackages) {
		numWorkers = len(allPackages)
	}

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for currentPkgPath := range workChan {
				// Report progress before processing
				reportProgress(currentPkgPath)

				// Scan the current package
				packageInfo, err := ScanPackage(currentPkgPath, basePkgUrl)
				if err != nil {
					errChan <- fmt.Errorf("failed to scan package %s: %w", currentPkgPath, err)
					continue
				}

				// Calculate the full package URL
				var fullPkgUrl string
				if currentPkgPath == "" {
					fullPkgUrl = basePkgUrl
				} else {
					fullPkgUrl = fmt.Sprintf("%s/%s", basePkgUrl, currentPkgPath)
				}

				// Invoke callback for current package (protect with mutex for thread safety)
				mu.Lock()
				callback(packageInfo, fullPkgUrl)
				completedWork++
				mu.Unlock()
			}
		}()
	}

	// Send all packages to work channel
	for _, pkg := range allPackages {
		workChan <- pkg
	}
	close(workChan)

	// Wait for all workers to complete
	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	// Report final 100% completion
	if progressCallback != nil {
		progressCallback(ProgressInfo{
			Completed:  totalDiscovered,
			Total:      totalDiscovered,
			Current:    "Completed",
			Percentage: 100.0,
		})
	}

	return nil
}

// findSubPackages discovers all sub-packages under the given package path
func findSubPackages(pkgPath string) []string {
	var dirPath string
	if pkgPath == "" {
		dirPath = "."
	} else {
		dirPath = pkgPath
	}

	// Use recursive helper function
	return findPackagesRecursively(dirPath, pkgPath)
}

// findPackagesRecursively recursively discovers all packages in directory structure
func findPackagesRecursively(dirPath, pkgPath string) []string {
	var subPackages []string

	entries, err := afero.ReadDir(sourceFs, dirPath)
	if err != nil {
		return subPackages
	}

	for _, entry := range entries {
		if !entry.IsDir() || shouldSkipDirectory(entry.Name()) {
			continue
		}

		subDirPath := filepath.Join(dirPath, entry.Name())

		// Construct sub-package path
		var subPkgPath string
		if pkgPath == "" {
			subPkgPath = entry.Name()
		} else {
			subPkgPath = fmt.Sprintf("%s/%s", pkgPath, entry.Name())
		}

		// FIXED: Add ALL directories to scan queue
		subPackages = append(subPackages, subPkgPath)

		// FIXED: Recursively search subdirectories
		nestedPackages := findPackagesRecursively(subDirPath, subPkgPath)
		subPackages = append(subPackages, nestedPackages...)
	}

	return subPackages
}

// shouldSkipDirectory determines if a directory should be skipped during package scanning
func shouldSkipDirectory(dirName string) bool {
	skipDirs := map[string]bool{
		"vendor":   true,
		".git":     true,
		".idea":    true,
		".vscode":  true,
		"testdata": true,
	}

	// Skip hidden directories (starting with .)
	if strings.HasPrefix(dirName, ".") && dirName != "." {
		return true
	}

	return skipDirs[dirName]
}

// ScanPackage is an alias for ScanSinglePackage for backward compatibility
var ScanPackage = ScanSinglePackage
