package pkg

import (
	"golang.org/x/tools/go/packages"
)

// ScanPackage scans the specified package and returns comprehensive information
func ScanPackage(pkgPath string) (*PackageInfo, error) {
	cfg := &packages.Config{
		Mode: packages.NeedFiles,
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
		})
	}

	return &PackageInfo{
		Files: files,
	}, nil
}
