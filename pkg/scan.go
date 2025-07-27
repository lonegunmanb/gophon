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

	// Extract file information
	for _, file := range pkg.GoFiles {
		files = append(files, FileInfo{
			FileName: file,
			Package:  pkgPath,
		})
	}

	// Extract constants from AST
	for _, file := range pkg.Syntax {
		if file == nil {
			continue
		}

		fileName := pkg.Fset.Position(file.Pos()).Filename
		fileInfo := &FileInfo{
			FileName: fileName,
			Package:  pkgPath,
		}

		// Walk through declarations to find constants
		for _, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
				for _, spec := range genDecl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						for _, name := range valueSpec.Names {
							// Get line numbers for the constant
							startPos := pkg.Fset.Position(spec.Pos())
							endPos := pkg.Fset.Position(spec.End())

							constantInfo := ConstantInfo{
								Name:        name.Name,
								PackagePath: pkgPath,
								Range: &Range{
									FileInfo:  fileInfo,
									StartLine: startPos.Line,
									EndLine:   endPos.Line,
								},
							}
							constants = append(constants, constantInfo)
						}
					}
				}
			}
		}
	}

	return &PackageInfo{
		Files:     files,
		Constants: constants,
	}, nil
}
