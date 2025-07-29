package pkg

import (
	"fmt"
	"github.com/spf13/afero"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

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

// ScanPackagesRecursively recursively scans all packages starting from the specified path
// and invokes the callback function for each package found. It uses afero.Fs for file system operations
// to enable easy testing with mocked file systems.
// Parameters:
//   - fs: The afero filesystem to use for file operations
//   - pkgPath: The relative package path to start scanning from (e.g., "pkg/utils")
//   - basePkgUrl: The base package URL/module path (e.g., "github.com/user/project")
//   - callback: Function called for each package, receives *PackageInfo and full package URL
func ScanPackagesRecursively(pkgPath, basePkgUrl string, callback func(*PackageInfo, string)) error {
	// Scan the current package
	packageInfo, err := ScanPackage(pkgPath, basePkgUrl)
	if err != nil {
		return fmt.Errorf("failed to scan package %s: %w", pkgPath, err)
	}

	// Calculate the full package URL
	var fullPkgUrl string
	if pkgPath == "" {
		fullPkgUrl = basePkgUrl
	} else {
		fullPkgUrl = fmt.Sprintf("%s/%s", basePkgUrl, pkgPath)
	}

	// Invoke callback for current package
	callback(packageInfo, fullPkgUrl)

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
		// If we can't read the directory, just return without error
		// This handles cases where the package path doesn't correspond to a physical directory
		return nil
	}

	// Recursively scan subdirectories that contain Go files
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

			// Recursively scan the sub-package
			if err = ScanPackagesRecursively(subPkgPath, basePkgUrl, callback); err != nil {
				return fmt.Errorf("failed to scan sub-package %s: %w", subPkgPath, err)
			}
		}
	}

	return nil
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
