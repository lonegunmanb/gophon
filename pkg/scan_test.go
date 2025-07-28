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
		if filepath.IsAbs(file.FileName) && filepath.Base(file.FileName) == "subjects.go" && file.Package == "github.com/lonegunmanb/gophon/pkg/testharness" {
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
	expectedPackage := "github.com/lonegunmanb/gophon/pkg/testharness"
	for _, file := range packageResult.Files {
		assert.Equal(t, expectedPackage, file.Package,
			"All files should belong to the exact package being scanned")
	}
}

func TestScanPackage_PackageNameDifferentFromDirectoryName(t *testing.T) {
	// Act - scan the mismatched_dir directory which contains package different_pkg
	packagePath := "test-harness/mismatched_dir"
	result, err := ScanPackage(packagePath, "github.com/lonegunmanb/gophon/pkg")

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify the TestVariable is found
	require.Len(t, result.Variables, 1, "Should find exactly one variable")
	assert.Equal(t, "TestVariable", result.Variables[0].Name, "Should find TestVariable")

	// Verify variable's package path - should use directory-based path for consistency
	expectedPackagePath := "github.com/lonegunmanb/gophon/pkg/test-harness/different_pkg"
	assert.Equal(t, expectedPackagePath, result.Variables[0].PackagePath,
		"Variable package path should use directory-based path for external accessibility")

	// Verify one file is included (example.go)
	require.Len(t, result.Files, 1, "Should find exactly one file")
	assert.Equal(t, "example.go", filepath.Base(result.Files[0].FileName), "Should find example.go")

	// Verify file's package path - should also use directory-based path
	assert.Equal(t, expectedPackagePath, result.Files[0].Package,
		"File package path should use directory-based path for external accessibility")
}

// scanHarnessPackage is a helper function that scans the test-harness package and returns the package result
func scanHarnessPackage(t *testing.T) *PackageInfo {
	packagePath := "test-harness"
	result, err := ScanPackage(packagePath, "github.com/lonegunmanb/gophon/pkg")

	require.NoError(t, err)
	require.NotNil(t, result)

	return result
}
