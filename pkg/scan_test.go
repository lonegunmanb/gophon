package pkg

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanPackage_ContainsSubjectsFile(t *testing.T) {
	// Act
	packageResult := scanHarnessPackage(t)

	// Assert - find the subjects.go file and check it has absolute path
	found := false
	for _, file := range packageResult.Files {
		if filepath.IsAbs(file.FileName) && filepath.Base(file.FileName) == "subjects.go" && file.Package == "github.com/lonegunmanb/gophon/pkg/test-harness" {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected to find subjects.go with absolute path in the scanned package files")
}

func TestScanPackage_ExcludesSubPackageFiles(t *testing.T) {
	// Act
	packageResult := scanHarnessPackage(t)

	// Assert - ensure no files from sub_pkg are included
	for _, file := range packageResult.Files {
		fileName := filepath.Base(file.FileName)
		filePath := file.FileName
		
		// Check that should_not_appear.go is not included
		assert.NotEqual(t, "should_not_appear.go", fileName, 
			"should_not_appear.go from sub_pkg should not be included when scanning only test-harness package")
		
		// Check that no files from sub_pkg directory are included
		assert.NotContains(t, filePath, "sub_pkg", 
			"Files from sub_pkg directory should not be included when scanning only test-harness package")
	}
	
	// Additional verification: ensure we only have files from the exact package
	expectedPackage := "github.com/lonegunmanb/gophon/pkg/test-harness"
	for _, file := range packageResult.Files {
		assert.Equal(t, expectedPackage, file.Package, 
			"All files should belong to the exact package being scanned")
	}
}

// scanHarnessPackage is a helper function that scans the test-harness package and returns the package result
func scanHarnessPackage(t *testing.T) *PackageInfo {
	packagePath := "test-harness"
	result, err := ScanPackage(packagePath, "github.com/lonegunmanb/gophon/pkg")

	require.NoError(t, err)
	require.NotNil(t, result)

	return result
}
