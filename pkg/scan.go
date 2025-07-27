package pkg

import (
	"fmt"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/packages"
)

// ScanPackage scans the specified package and returns comprehensive information
func ScanPackage(pkgPath string, basePkgUrl string) (*PackageInfo, error) {
	pkgPath = fmt.Sprintf("%s/%s", basePkgUrl, pkgPath)
	cfg := &packages.Config{
		Mode: packages.NeedFiles | packages.NeedName | packages.NeedImports | packages.NeedTypes | packages.NeedSyntax,
	}

	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil {
		return nil, err
	}

	if len(pkgs) == 0 {
		return &PackageInfo{}, nil
	}

	pkg := pkgs[0]
	var files []FileInfo
	var constants []ConstantInfo
	var variables []VariableInfo
	var types []TypeInfo
	var functions []FunctionInfo

	// Extract file information
	for _, file := range pkg.GoFiles {
		files = append(files, FileInfo{
			FileName: file,
			Package:  pkgPath,
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
			Package:  pkgPath,
		}

		// Walk through declarations to find constants, variables, types, and functions
		for _, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok {
				if genDecl.Tok == token.CONST {
					constants = append(constants, extractDeclarations(pkgPath, genDecl, pkg, fileInfo, func(name string, pkgPath string, rangeInfo *Range) ConstantInfo {
						return ConstantInfo{
							Name:        name,
							PackagePath: pkgPath,
							Range:       rangeInfo,
						}
					})...)
				} else if genDecl.Tok == token.VAR {
					variables = append(variables, extractDeclarations(pkgPath, genDecl, pkg, fileInfo, func(name string, pkgPath string, rangeInfo *Range) VariableInfo {
						return VariableInfo{
							Name:        name,
							PackagePath: pkgPath,
							Range:       rangeInfo,
						}
					})...)
				} else if genDecl.Tok == token.TYPE {
					types = append(types, extractTypeDeclarations(pkgPath, genDecl, pkg, fileInfo)...)
				}
			} else if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				functions = append(functions, extractFunctionDeclarations(pkgPath, funcDecl, pkg, fileInfo)...)
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
func extractDeclarations[T any](pkgPath string, genDecl *ast.GenDecl, pkg *packages.Package, fileInfo *FileInfo, createFunc func(name string, pkgPath string, rangeInfo *Range) T) []T {
	var results []T
	for _, spec := range genDecl.Specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			for _, name := range valueSpec.Names {
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
func extractTypeDeclarations(pkgPath string, genDecl *ast.GenDecl, pkg *packages.Package, fileInfo *FileInfo) []TypeInfo {
	var results []TypeInfo
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

			results = append(results, TypeInfo{
				Name:        typeSpec.Name.Name,
				PackagePath: pkgPath,
				Range:       rangeInfo,
			})
		}
	}
	return results
}

// Extract function declarations from AST
func extractFunctionDeclarations(pkgPath string, funcDecl *ast.FuncDecl, pkg *packages.Package, fileInfo *FileInfo) []FunctionInfo {
	var results []FunctionInfo

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

	results = append(results, FunctionInfo{
		Range:        rangeInfo,
		Name:         funcDecl.Name.Name,
		ReceiverType: receiverType,
		PackagePath:  pkgPath,
	})

	return results
}
