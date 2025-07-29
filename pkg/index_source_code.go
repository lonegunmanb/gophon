package pkg

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

var destFs = afero.NewOsFs()

// IndexSourceCode recursively scans packages and generates index files for all indexable symbols.
// It uses ScanPackagesRecursively to discover packages and generates individual .goindex files
// for each symbol (constants, variables, types, functions, methods) in the destination filesystem.
//
// Parameters:
//   - pkgPath: The relative package path to start scanning from (e.g., "testharness")
//   - basePkgUrl: The base package URL/module path (e.g., "github.com/lonegunmanb/gophon/pkg")
//   - destFolder: The destination folder path where index files will be organized
//   - progressCallback: Optional function called with progress updates, can be nil
func IndexSourceCode(pkgPath, basePkgUrl string, destFolder string, progressCallback func(ProgressInfo)) error {
	// Define the callback function that will be called for each package
	callback := func(pkgInfo *PackageInfo, pkgUrl string) {
		// Extract the relative package path from the full package URL
		relativePkgPath := strings.TrimPrefix(pkgUrl, basePkgUrl)
		relativePkgPath = strings.TrimPrefix(relativePkgPath, "/")

		// Create the destination directory for this package
		pkgDestDir := filepath.Join(destFolder, relativePkgPath)

		// Process all indexable symbols in this package
		saveIndexes(pkgDestDir, pkgInfo.Constants)
		saveIndexes(pkgDestDir, pkgInfo.Variables)
		saveIndexes(pkgDestDir, pkgInfo.Types)
		saveIndexes(pkgDestDir, pkgInfo.Functions)
	}

	// Call ScanPackagesRecursively with our callback and progress tracking
	return ScanPackagesRecursively(pkgPath, basePkgUrl, callback, progressCallback)
}

// IndexSourceCodeWithoutProgress is a backward-compatible wrapper that calls
// IndexSourceCode without progress tracking
func IndexSourceCodeWithoutProgress(pkgPath, basePkgUrl string, destFolder string) error {
	return IndexSourceCode(pkgPath, basePkgUrl, destFolder, nil)
}

// saveConstants creates index files for all constants in the package
func saveIndexes[T IndexableSymbol](pkgDestDir string, indexes []T) {
	for _, index := range indexes {
		// Get the index filename using the IndexableSymbol interface
		filename := index.IndexFileName()

		// Create the full file path
		filePath := filepath.Join(pkgDestDir, filename)

		// Ensure the directory exists
		dir := filepath.Dir(filePath)
		if err := destFs.MkdirAll(dir, 0700); err != nil {
			// Log error but continue processing other files
			fmt.Printf("Warning: Failed to create directory %s: %v\n", dir, err)
			return
		}

		// Generate the index file content
		content := generateIndexContent(index)

		// Write the index file
		if err := afero.WriteFile(destFs, filePath, []byte(content), 0600); err != nil {
			// Log error but continue processing other files
			fmt.Printf("Warning: Failed to write index file %s: %v\n", filePath, err)
		}
	}
}

// generateIndexContent generates the content for an index file
func generateIndexContent(symbol IndexableSymbol) string {
	return fmt.Sprintf(`package %s
%s
%s
`, symbol.PackagePath(), symbol.Imports(), symbol.String())
}
