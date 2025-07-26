package pkg

import (
	"fmt"
	"golang.org/x/tools/go/packages"
)

// ScanPackage scans the specified package and returns comprehensive information
func ScanPackage(pkgPath string, basePkgUrl string) (*PackageInfo, error) {
	pkgPath = fmt.Sprintf("%s/%s", basePkgUrl, pkgPath)
	cfg := &packages.Config{
		Mode: packages.NeedFiles | packages.NeedName | packages.NeedImports | packages.NeedTypes,
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

	for _, file := range pkg.GoFiles {
		files = append(files, FileInfo{
			FileName: file,
			Package:  pkgPath,
		})
	}

	return &PackageInfo{
		Files: files,
	}, nil
}
