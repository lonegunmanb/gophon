package pkg

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScanPackage_ExtractsConstants(t *testing.T) {
	// Act
	packageResult := scanHarnessPackage(t)

	// Assert - check that constants are extracted
	assert.Len(t, packageResult.Constants, 2)

	// Find specific constants we expect
	defaultTimeoutConst := findConstantByName(packageResult.Constants, "DefaultTimeout")
	maxRetriesConst := findConstantByName(packageResult.Constants, "maxRetries")

	// Verify DefaultTimeout constant
	require.NotNil(t, defaultTimeoutConst, "Should find DefaultTimeout constant")
	assert.Equal(t, "DefaultTimeout", defaultTimeoutConst.Name)
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/testharness", defaultTimeoutConst.PackagePath)

	// Assert String() method returns the exact source code from subjects.go
	assert.Equal(t, "\tDefaultTimeout = 30 * time.Second", defaultTimeoutConst.String(), "DefaultTimeout String() should return exact source code line")

	// Verify MaxRetries constant
	require.NotNil(t, maxRetriesConst, "Should find MaxRetries constant")
	assert.Equal(t, "maxRetries", maxRetriesConst.Name)
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/testharness", maxRetriesConst.PackagePath)
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/testharness", maxRetriesConst.PackagePath)

	// Assert String() method returns the exact source code from subjects.go
	assert.Equal(t, "const maxRetries = 3", maxRetriesConst.String(), "maxRetries String() should return exact source code line")
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
