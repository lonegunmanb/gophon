package pkg

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// findSubjectsFile is a helper function that scans the test-harness package and returns the package result
func findSubjectsFile(t *testing.T) *PackageInfo {
	packagePath := "test-harness"
	result, err := ScanPackage(packagePath, "github.com/lonegunmanb/gophon/pkg")

	require.NoError(t, err)
	require.NotNil(t, result)

	return result
}

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
