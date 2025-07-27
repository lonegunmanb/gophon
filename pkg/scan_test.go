package pkg

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanPackage_ContainsSubjectsFile(t *testing.T) {
	// Act
	packageResult := findSubjectsFile(t)

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

func TestScanPackage_ExtractsConstants(t *testing.T) {
	// Act
	packageResult := findSubjectsFile(t)

	// Assert - check that constants are extracted
	assert.Len(t, packageResult.Constants, 2)

	// Find specific constants we expect
	defaultTimeoutConst := findConstantByName(packageResult.Constants, "DefaultTimeout")
	maxRetriesConst := findConstantByName(packageResult.Constants, "maxRetries")

	// Verify DefaultTimeout constant
	require.NotNil(t, defaultTimeoutConst, "Should find DefaultTimeout constant")
	assert.Equal(t, "DefaultTimeout", defaultTimeoutConst.Name)
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/test-harness", defaultTimeoutConst.PackagePath)

	// Assert String() method returns the exact source code from subjects.go
	assert.Equal(t, "\tDefaultTimeout = 30 * time.Second", defaultTimeoutConst.String(), "DefaultTimeout String() should return exact source code line")

	// Verify MaxRetries constant
	require.NotNil(t, maxRetriesConst, "Should find MaxRetries constant")
	assert.Equal(t, "maxRetries", maxRetriesConst.Name)
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/test-harness", maxRetriesConst.PackagePath)

	// Assert String() method returns the exact source code from subjects.go
	assert.Equal(t, "const maxRetries = 3", maxRetriesConst.String(), "maxRetries String() should return exact source code line")
}

// findSubjectsFile is a helper function that scans the test-harness package and returns the package result
func findSubjectsFile(t *testing.T) *PackageInfo {
	packagePath := "test-harness"
	result, err := ScanPackage(packagePath, "github.com/lonegunmanb/gophon/pkg")

	require.NoError(t, err)
	require.NotNil(t, result)

	return result
}

// findConstantByName is a helper function that finds a constant by name in a slice of ConstantInfo
func findConstantByName(constants []ConstantInfo, name string) *ConstantInfo {
	for i := range constants {
		if constants[i].Name == name {
			return &constants[i]
		}
	}
	return nil
}
