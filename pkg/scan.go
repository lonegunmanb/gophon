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

	// Extract file information
	for _, file := range pkg.GoFiles {
		files = append(files, FileInfo{
			FileName: file,
			Package:  pkgPath,
		})
	}

	// Extract constants and variables from AST
	for _, file := range pkg.Syntax {
		if file == nil {
			continue
		}

		fileName := pkg.Fset.Position(file.Pos()).Filename
		fileInfo := &FileInfo{
			FileName: fileName,
			Package:  pkgPath,
		}

		// Walk through declarations to find constants and variables
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
				}
			}
		}
	}

	return &PackageInfo{
		Files:     files,
		Constants: constants,
		Variables: variables,
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
