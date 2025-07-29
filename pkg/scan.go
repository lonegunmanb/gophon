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

// ScanSinglePackage scans the specified package and returns comprehensive information
func ScanSinglePackage(pkgPath, basePkgUrl string) (*PackageInfo, error) {
	loadPath := fmt.Sprintf("%s/%s", basePkgUrl, pkgPath)
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
		parts := strings.Split(pkgPath, "/")
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

	// Extract file information
	for _, file := range pkg.GoFiles {
		files = append(files, &FileInfo{
			FileName: file,
			Package:  actualPkgPath,
		})
	}

	// Extract constants, variables, and types from AST
	for _, file := range pkg.Syntax {
		if file == nil {
			continue
		}

		fileName := pkg.Fset.Position(file.Pos()).Filename
		fileInfo := &FileInfo{
			FileName: fileName,
			Package:  actualPkgPath,
		}

		// Walk through declarations to find constants, variables, types, and functions
		for _, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok {
				switch genDecl.Tok {
				case token.CONST:
					constants = append(constants, extractDeclarations(actualPkgPath, genDecl, pkg, fileInfo, func(name string, pkgPath string, rangeInfo *Range) *ConstantInfo {
						return &ConstantInfo{
							Name:  name,
							Range: rangeInfo,
						}
					})...)
				case token.VAR:
					variables = append(variables, extractDeclarations(actualPkgPath, genDecl, pkg, fileInfo, func(name string, pkgPath string, rangeInfo *Range) *VariableInfo {
						return &VariableInfo{
							Name:  name,
							Range: rangeInfo,
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
				Name:  typeSpec.Name.Name,
				Range: rangeInfo,
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
		Name:         funcDecl.Name.Name,
		ReceiverType: receiverType,
	})

	return results
}

// ProgressInfo represents progress information during package scanning
type ProgressInfo struct {
	Completed   int     // Number of packages completed
	Total       int     // Total number of packages discovered so far
	Current     string  // Currently processing package path
	Percentage  float64 // Completion percentage (completed/total * 100)
}

// ScanPackagesRecursively recursively scans all packages starting from the specified path
// and invokes the callback function for each package found. It uses a worker pool pattern
// with goroutines limited to the number of processors for optimal performance.
// Parameters:
//   - pkgPath: The relative package path to start scanning from (e.g., "pkg/utils")
//   - basePkgUrl: The base package URL/module path (e.g., "github.com/user/project")
//   - callback: Function called for each package, receives *PackageInfo and full package URL
//   - progressCallback: Optional function called with progress updates, can be nil
func ScanPackagesRecursively(pkgPath, basePkgUrl string, callback func(*PackageInfo, string), progressCallback func(ProgressInfo)) error {
	// Use a worker pool with goroutines limited to the number of processors
	numWorkers := runtime.NumCPU()
	
	// Channel for work items (package paths to scan)
	workChan := make(chan scanWork, 100)
	
	// Channel for results
	resultChan := make(chan scanResult, 100)
	
	// WaitGroup to track workers
	var wg sync.WaitGroup
	
	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			scanWorker(workChan, resultChan)
		}()
	}
	
	// Start a goroutine to close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	// Send initial work item
	workChan <- scanWork{pkgPath: pkgPath, basePkgUrl: basePkgUrl}
	
	// Keep track of pending work and completed work for progress tracking
	pendingWork := 1
	completedWork := 0
	totalDiscovered := 1
	var mu sync.Mutex
	
	// Helper function to report progress
	reportProgress := func(currentPkg string) {
		if progressCallback != nil {
			mu.Lock()
			completed := completedWork
			total := totalDiscovered
			mu.Unlock()
			
			percentage := 0.0
			if total > 0 {
				percentage = float64(completed) / float64(total) * 100.0
			}
			
			progressCallback(ProgressInfo{
				Completed:  completed,
				Total:      total,
				Current:    currentPkg,
				Percentage: percentage,
			})
		}
	}
	
	// Process results and generate more work
	for result := range resultChan {
		if result.err != nil {
			close(workChan)
			return fmt.Errorf("failed to scan package %s: %w", result.work.pkgPath, result.err)
		}
		
		// Calculate the full package URL
		var fullPkgUrl string
		if result.work.pkgPath == "" {
			fullPkgUrl = result.work.basePkgUrl
		} else {
			fullPkgUrl = fmt.Sprintf("%s/%s", result.work.basePkgUrl, result.work.pkgPath)
		}
		
		// Report progress before processing
		reportProgress(fullPkgUrl)
		
		// Invoke callback for current package
		callback(result.packageInfo, fullPkgUrl)
		
		// Find subdirectories and add them as new work items
		subPackages := findSubPackages(result.work.pkgPath)
		
		mu.Lock()
		for _, subPkg := range subPackages {
			workChan <- scanWork{pkgPath: subPkg, basePkgUrl: result.work.basePkgUrl}
			pendingWork++
			totalDiscovered++
		}
		pendingWork--
		completedWork++
		
		// Close work channel when no more work is pending
		if pendingWork == 0 {
			close(workChan)
		}
		mu.Unlock()
	}
	
	// Report final progress (100% completion)
	if progressCallback != nil {
		mu.Lock()
		completed := completedWork
		total := totalDiscovered
		mu.Unlock()
		
		progressCallback(ProgressInfo{
			Completed:  completed,
			Total:      total,
			Current:    "Completed",
			Percentage: 100.0,
		})
	}
	
	return nil
}

// scanWork represents a work item for the worker pool
type scanWork struct {
	pkgPath    string
	basePkgUrl string
}

// scanResult represents the result of scanning a package
type scanResult struct {
	work        scanWork
	packageInfo *PackageInfo
	err         error
}

// scanWorker is a worker function that processes scan work items
func scanWorker(workChan <-chan scanWork, resultChan chan<- scanResult) {
	for work := range workChan {
		packageInfo, err := ScanPackage(work.pkgPath, work.basePkgUrl)
		resultChan <- scanResult{
			work:        work,
			packageInfo: packageInfo,
			err:         err,
		}
	}
}

// findSubPackages finds all sub-packages within the given package path
func findSubPackages(pkgPath string) []string {
	var subPackages []string
	
	// Determine the physical directory path to scan for subdirectories
	var dirPath string
	if pkgPath == "" {
		dirPath = "."
	} else {
		dirPath = pkgPath
	}

	// Read directory contents using afero filesystem
	entries, err := afero.ReadDir(sourceFs, dirPath)
	if err != nil {
		// If we can't read the directory, return empty slice
		return subPackages
	}

	// Find subdirectories that contain Go files
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Skip common non-package directories
		if shouldSkipDirectory(entry.Name()) {
			continue
		}

		subDirPath := filepath.Join(dirPath, entry.Name())

		// Check if the subdirectory contains any .go files
		if hasGoFiles(subDirPath) {
			// Construct the sub-package path
			var subPkgPath string
			if pkgPath == "" {
				subPkgPath = entry.Name()
			} else {
				subPkgPath = fmt.Sprintf("%s/%s", pkgPath, entry.Name())
			}

			subPackages = append(subPackages, subPkgPath)
		}
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

// hasGoFiles checks if a directory contains any .go files
func hasGoFiles(dirPath string) bool {
	entries, err := afero.ReadDir(sourceFs, dirPath)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			// Skip test files for package detection
			if !strings.HasSuffix(entry.Name(), "_test.go") {
				return true
			}
		}
	}

	return false
}

// ScanPackage is an alias for ScanSinglePackage for backward compatibility
var ScanPackage = ScanSinglePackage

// ScanPackagesRecursivelyWithoutProgress is a backward-compatible wrapper that calls
// ScanPackagesRecursively without progress tracking
func ScanPackagesRecursivelyWithoutProgress(pkgPath, basePkgUrl string, callback func(*PackageInfo, string)) error {
	return ScanPackagesRecursively(pkgPath, basePkgUrl, callback, nil)
}
