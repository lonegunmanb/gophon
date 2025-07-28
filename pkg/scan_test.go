package pkg

import (
	"fmt"
	"github.com/prashantv/gostub"
	"github.com/spf13/afero"
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
	assert.Equal(t, expectedPackagePath, result.Variables[0].PackagePath,
		"Variable package path should use directory-based path for external accessibility")

	// Verify one file is included (example.go)
	require.Len(t, result.Files, 1, "Should find exactly one file")
	assert.Equal(t, "example.go", filepath.Base(result.Files[0].FileName), "Should find example.go")

	// Verify file's package path - should also use directory-based path
	assert.Equal(t, expectedPackagePath, result.Files[0].Package,
		"File package path should use directory-based path for external accessibility")
}

func TestScanPackagesRecursively(t *testing.T) {
	// Create a mock file system using the helper function
	files := map[string]string{
		// Root package
		"main.go": `package main

func main() {
	fmt.Println("Hello, World!")
}
`,
		// Sub-package: utils
		"utils/helper.go": `package utils

func Helper() string {
	return "helper"
}
`,
		// Sub-package: models
		"models/user.go": `package models

type User struct {
	Name string
	Age  int
}
`,
		// Nested sub-package: models/db
		"models/db/connection.go": `package db

type Connection struct {
	URL string
}
`,
		// Directory without Go files (should be skipped)
		"docs/README.md": "# Documentation",
	}

	mockFs := afero.NewMemMapFs()
	setupMemoryFilesystem(mockFs, files)

	// Stub the filesystem variable to use our memory filesystem
	stubs := gostub.Stub(&fs, mockFs).Stub(&ScanPackage, func(pkgPath, basePkgUrl string) (*PackageInfo, error) {
		// Return mock package info based on the path
		var packageName string
		switch pkgPath {
		case "":
			packageName = "main"
		case "utils":
			packageName = "utils"
		case "models":
			packageName = "models"
		case "models/db":
			packageName = "db"
		default:
			packageName = "unknown"
		}

		return &PackageInfo{
			Files: []FileInfo{
				{
					FileName: fmt.Sprintf("%s/%s.go", pkgPath, packageName),
					Package:  fmt.Sprintf("%s/%s", basePkgUrl, packageName),
				},
			},
			Constants: []ConstantInfo{},
			Variables: []VariableInfo{},
			Types:     []TypeInfo{},
			Functions: []FunctionInfo{},
		}, nil
	})
	defer stubs.Reset()

	// Collect results from the callback
	var results []struct {
		PackageInfo *PackageInfo
		PkgUrl      string
	}

	callback := func(pkgInfo *PackageInfo, pkgUrl string) {
		results = append(results, struct {
			PackageInfo *PackageInfo
			PkgUrl      string
		}{
			PackageInfo: pkgInfo,
			PkgUrl:      pkgUrl,
		})
	}

	// Test the recursive scanner with memory filesystem
	basePkgUrl := "github.com/example/testproject"
	require.NoError(t, ScanPackagesRecursively("", basePkgUrl, callback))
	require.NotEmpty(t, results)
	require.Len(t, results, 4)

	// Verify specific packages were found
	foundPackages := make(map[string]bool)
	for _, result := range results {
		foundPackages[result.PkgUrl] = true
	}

	for _, expectedUrl := range []string{
		"github.com/example/testproject",
		"github.com/example/testproject/utils",
		"github.com/example/testproject/models",
		"github.com/example/testproject/models/db",
	} {
		assert.Contains(t, foundPackages, expectedUrl)
	}
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
