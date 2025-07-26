package pkg

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanPackage_ContainsSubjectsFile(t *testing.T) {
	// Arrange
	packagePath := "test-harness"

	// Act
	result, err := ScanPackage(packagePath, "github.com/lonegunmanb/gophon/pkg")

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)

	// Check that the result contains a file with absolute path ending in "subjects.go"
	found := false
	for _, file := range result.Files {

		if filepath.Base(file.FileName) == "subjects.go" && file.Package == "github.com/lonegunmanb/gophon/pkg/test-harness" {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected to find subjects.go in the scanned package files")
}
