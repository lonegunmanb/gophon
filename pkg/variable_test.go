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

	// Assert - check that variables are extracted (should be 2, not 3, because blank identifier should be skipped)
	assert.Len(t, packageResult.Variables, 2, "Should extract 2 global variables from test harness (blank identifier should be skipped)")

	// Find specific variables we expect
	globalCounterVar := findVariableByName(packageResult.Variables, "GlobalCounter")
	isDebugModeVar := findVariableByName(packageResult.Variables, "isDebugMode")

	// Verify GlobalCounter variable
	require.NotNil(t, globalCounterVar, "Should find GlobalCounter variable")
	assert.Equal(t, "GlobalCounter", globalCounterVar.Name)
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/testharness", globalCounterVar.PackagePath())
	assert.Contains(t, globalCounterVar.FileName, "subjects.go")
	assert.True(t, filepath.IsAbs(globalCounterVar.FileName), "FileName should be absolute path")
	assert.NotNil(t, globalCounterVar.GenDecl)

	// Assert String() method returns the exact source code from subjects.go
	assert.Equal(t, "\tGlobalCounter int64", globalCounterVar.String(), "GlobalCounter String() should return exact source code line")

	// Verify IsDebugMode variable
	require.NotNil(t, isDebugModeVar, "Should find IsDebugMode variable")
	assert.Equal(t, "isDebugMode", isDebugModeVar.Name)
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/testharness", isDebugModeVar.PackagePath())
	assert.Contains(t, isDebugModeVar.FileName, "subjects.go")
	assert.True(t, filepath.IsAbs(isDebugModeVar.FileName), "FileName should be absolute path")
	assert.NotNil(t, isDebugModeVar.GenDecl)

	// Assert String() method returns the exact source code from subjects.go
	assert.Equal(t, "\tisDebugMode bool = false", isDebugModeVar.String(), "IsDebugMode String() should return exact source code line")
}

func TestScanPackage_SkipsBlankIdentifierVariables(t *testing.T) {
	// Act
	packageResult := scanHarnessPackage(t)

	// Assert - verify that no variable with blank identifier "_" is extracted
	for _, variable := range packageResult.Variables {
		assert.NotEqual(t, "_", variable.Name, "Blank identifier variables should be skipped during scanning")
	}

	// Also verify we don't have any variables that would be from the blank identifier line
	blankVar := findVariableByName(packageResult.Variables, "_")
	assert.Nil(t, blankVar, "Should not find any variable with blank identifier name")
}

func TestScanPackage_VariablePackagePath(t *testing.T) {
	// Act
	packageResult := scanHarnessPackage(t)

	// Assert - all variables should have correct package path
	expectedPackagePath := "github.com/lonegunmanb/gophon/pkg/testharness"
	for _, variable := range packageResult.Variables {
		assert.Equal(t, expectedPackagePath, variable.PackagePath(), "Variable %s should have correct package path", variable.Name)
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
			variable := &VariableInfo{
				Name: tt.variableName,
				Range: &Range{
					FileInfo: &FileInfo{
						Package: "github.com/example/pkg",
					},
				},
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
	variable := &VariableInfo{
		Name: "TestVariable",
		Range: &Range{
			FileInfo: &FileInfo{
				Package: "github.com/example/pkg",
			},
		},
	}

	// Act & Assert - this should compile without error, proving VariableInfo implements IndexableSymbol
	var _ IndexableSymbol = variable
	assert.Equal(t, "var.TestVariable.goindex", variable.IndexFileName())
}

// Helper function to find a variable by name
func findVariableByName(variables []*VariableInfo, name string) *VariableInfo {
	var variable *VariableInfo
	for _, v := range variables {
		if v.Name == name {
			variable = v
		}
	}
	return variable
}
