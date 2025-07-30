package pkg

import (
	"github.com/prashantv/gostub"
	"github.com/spf13/afero"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanPackage_longPackagePath(t *testing.T) {
	pkgPath := filepath.Join("testharness", "dir_without_go_file", "dir_without_go_file", "dir")
	pkg, err := ScanSinglePackage(pkgPath, "github.com/lonegunmanb/gophon/pkg")
	require.NoError(t, err)
	assert.Equal(t, pkg.Variables[0].Package, "github.com/lonegunmanb/gophon/pkg/testharness/dir_without_go_file/dir_without_go_file/dir")
}

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
		assert.NotNil(t, file.File)
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
			"should_not_appear.go from sub_pkg should not be included when scanning only testharness package")

		// Check that no files from sub_pkg directory are included
		assert.NotContains(t, filePath, "sub_pkg",
			"Files from sub_pkg directory should not be included when scanning only testharness package")
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
	packagePath := "testharness/mismatched_dir"
	result, err := ScanSinglePackage(packagePath, "github.com/lonegunmanb/gophon/pkg")

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify the TestVariable is found
	require.Len(t, result.Variables, 1, "Should find exactly one variable")
	assert.Equal(t, "TestVariable", result.Variables[0].Name, "Should find TestVariable")

	// Verify variable's package path - should use directory-based path for consistency
	expectedPackagePath := "github.com/lonegunmanb/gophon/pkg/testharness/different_pkg"
	assert.Equal(t, expectedPackagePath, result.Variables[0].PackagePath(),
		"Variable package path should use directory-based path for external accessibility")

	// Verify one file is included (example.go)
	require.Len(t, result.Files, 1, "Should find exactly one file")
	assert.Equal(t, "example.go", filepath.Base(result.Files[0].FileName), "Should find example.go")

	// Verify file's package path - should also use directory-based path
	assert.Equal(t, expectedPackagePath, result.Files[0].Package,
		"File package path should use directory-based path for external accessibility")
}

func TestScanPackagesRecursively(t *testing.T) {
	var found = false
	require.NoError(t, ScanPackagesRecursively("testharness/dir_without_go_file", "github.com/lonegunmanb/gophon/pkg", func(info *PackageInfo, s string) {
		if info == nil {
			return
		}
		for _, v := range info.Variables {
			if v.Name == "ShouldBeIncluded" {
				found = true
			}
		}
	}, nil))
	assert.True(t, found)
}

func TestFindPackagesRecursively_EmptyMiddleFolderShouldNotBeSkipped(t *testing.T) {
	// Setup test filesystem with empty middle directories
	files := map[string]string{
		"internal/services/compute/compute.go":       `package compute`,
		"internal/services/compute/worker/worker.go": `package worker`,
		// Note: internal/ and internal/services/ contain no Go files
	}

	mockFs := afero.NewMemMapFs()
	setupMemoryFilesystem(mockFs, files)

	// Stub the filesystem variable to use our memory filesystem
	stub := gostub.Stub(&sourceFs, mockFs)
	defer stub.Reset()

	// Test package discovery from root
	packages := findPackagesRecursively(".", "")

	expectedPackages := []string{
		"internal",
		"internal/services",
		"internal/services/compute",
		"internal/services/compute/worker",
	}

	for _, expected := range expectedPackages {
		assert.Contains(t, packages, expected, "Should find package: %s", expected)
	}

	// Verify we found at least the expected packages
	assert.GreaterOrEqual(t, len(packages), len(expectedPackages),
		"Should find at least %d packages, found %d: %v", len(expectedPackages), len(packages), packages)
}

// scanHarnessPackage is a helper function that scans the testharness package and returns the package result
func scanHarnessPackage(t *testing.T) *PackageInfo {
	packagePath := "testharness"
	result, err := ScanSinglePackage(packagePath, "github.com/lonegunmanb/gophon/pkg")

	require.NoError(t, err)
	require.NotNil(t, result)

	return result
}

// setupMemoryFilesystem creates a memory filesystem with the given files
// files map contains filepath -> file content
func setupMemoryFilesystem(fs afero.Fs, files map[string]string) {

	for filePath, content := range files {
		exists, _ := afero.Exists(fs, filepath.Base(filePath))
		if !exists {
			// Write the file
			_ = afero.WriteFile(fs, filePath, []byte(content), 0644)
		}
	}
}
