package pkg

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestScanPackage_ExtractsVariables(t *testing.T) {
	// Act
	packageResult := scanHarnessPackage(t)

	// Assert - check that variables are extracted
	assert.Len(t, packageResult.Variables, 2, "Should extract 2 global variables from test harness")

	// Find specific variables we expect
	globalCounterVar := findVariableByName(packageResult.Variables, "GlobalCounter")
	isDebugModeVar := findVariableByName(packageResult.Variables, "isDebugMode")

	// Verify GlobalCounter variable
	require.NotNil(t, globalCounterVar, "Should find GlobalCounter variable")
	assert.Equal(t, "GlobalCounter", globalCounterVar.Name)
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/testharness", globalCounterVar.PackagePath)
	assert.Contains(t, globalCounterVar.FileName, "subjects.go")
	assert.True(t, filepath.IsAbs(globalCounterVar.FileName), "FileName should be absolute path")

	// Assert String() method returns the exact source code from subjects.go
	assert.Equal(t, "\tGlobalCounter int64", globalCounterVar.String(), "GlobalCounter String() should return exact source code line")

	// Verify IsDebugMode variable
	require.NotNil(t, isDebugModeVar, "Should find IsDebugMode variable")
	assert.Equal(t, "isDebugMode", isDebugModeVar.Name)
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/testharness", isDebugModeVar.PackagePath)
	assert.Contains(t, isDebugModeVar.FileName, "subjects.go")
	assert.True(t, filepath.IsAbs(isDebugModeVar.FileName), "FileName should be absolute path")

	// Assert String() method returns the exact source code from subjects.go
	assert.Equal(t, "\tisDebugMode bool = false", isDebugModeVar.String(), "IsDebugMode String() should return exact source code line")
}

func TestScanPackage_VariablePackagePath(t *testing.T) {
	// Act
	packageResult := scanHarnessPackage(t)

	// Assert - all variables should have correct package path
	expectedPackagePath := "github.com/lonegunmanb/gophon/pkg/testharness"
	for _, variable := range packageResult.Variables {
		assert.Equal(t, expectedPackagePath, variable.PackagePath, "Variable %s should have correct package path", variable.Name)
	}
}

func TestVariableInfo_IndexFileName(t *testing.T) {
	// Arrange
	tests := []struct {
		name         string
		variableName string
		expected     string
	}{
		{
			name:         "Simple variable name",
			variableName: "GlobalCounter",
			expected:     "var.GlobalCounter.goindex",
		},
		{
			name:         "Variable with lowercase name",
			variableName: "isDebugMode",
			expected:     "var.isDebugMode.goindex",
		},
		{
			name:         "Variable with mixed case",
			variableName: "DefaultTimeout",
			expected:     "var.DefaultTimeout.goindex",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			variable := VariableInfo{
				Name:        tt.variableName,
				PackagePath: "github.com/example/pkg",
			}

			// Act
			result := variable.IndexFileName()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVariableInfo_ImplementsIndexableSymbol(t *testing.T) {
	// Arrange
	variable := VariableInfo{
		Name:        "TestVariable",
		PackagePath: "github.com/example/pkg",
	}

	// Act & Assert - this should compile without error, proving VariableInfo implements IndexableSymbol
	var _ IndexableSymbol = variable
	assert.Equal(t, "var.TestVariable.goindex", variable.IndexFileName())
}

// Helper function to find a variable by name
func findVariableByName(variables []VariableInfo, name string) VariableInfo {
	var variable VariableInfo
	for _, v := range variables {
		if v.Name == name {
			variable = v
		}
	}
	return variable
}
