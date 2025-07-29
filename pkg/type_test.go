package pkg

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestScanPackage_ExtractsTypes(t *testing.T) {
	// Act
	packageResult := scanHarnessPackage(t)

	// Assert - check that types are extracted
	assert.Len(t, packageResult.Types, 5, "Should extract 5 type declarations from test harness")

	// Find specific types we expect
	stringAType := findTypeByName(packageResult.Types, "StringA")
	stringBType := findTypeByName(packageResult.Types, "StringB")
	userType := findTypeByName(packageResult.Types, "User")
	userServiceType := findTypeByName(packageResult.Types, "UserService")
	serviceType := findTypeByName(packageResult.Types, "Service")

	// Verify StringA type (type alias)
	require.NotNil(t, stringAType, "Should find StringA type")
	assert.Equal(t, "StringA", stringAType.Name)
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/testharness", stringAType.PackagePath())
	assert.Contains(t, stringAType.FileName, "subjects.go")
	assert.True(t, filepath.IsAbs(stringAType.FileName), "FileName should be absolute path")

	// Assert String() method returns the exact source code from subjects.go
	assert.Equal(t, "type StringA string", stringAType.String(), "StringA String() should return exact source code line")

	// Verify StringB type (type alias with =)
	require.NotNil(t, stringBType, "Should find StringB type")
	assert.Equal(t, "StringB", stringBType.Name)
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/testharness", stringBType.PackagePath())
	assert.Contains(t, stringBType.FileName, "subjects.go")
	assert.True(t, filepath.IsAbs(stringBType.FileName), "FileName should be absolute path")

	// Assert String() method returns the exact source code from subjects.go
	assert.Equal(t, "type StringB = string", stringBType.String(), "StringB String() should return exact source code line")

	// Verify User struct type
	require.NotNil(t, userType, "Should find User type")
	assert.Equal(t, "User", userType.Name)
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/testharness", userType.PackagePath())
	assert.Contains(t, userType.FileName, "subjects.go")
	assert.True(t, filepath.IsAbs(userType.FileName), "FileName should be absolute path")

	// Assert String() method returns the exact source code from subjects.go for User struct
	expectedUserSource := `type User struct {
	ID    int64  ` + "`json:\"id\" db:\"user_id\"`" + `
	Name  string ` + "`json:\"name\" db:\"full_name\"`" + `
	Email string ` + "`json:\"email\" db:\"email\"`" + `
}`
	assert.Equal(t, expectedUserSource, userType.String(), "User String() should return exact source code")

	// Verify UserService interface type
	require.NotNil(t, userServiceType, "Should find UserService type")
	assert.Equal(t, "UserService", userServiceType.Name)
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/testharness", userServiceType.PackagePath())
	assert.Contains(t, userServiceType.FileName, "subjects.go")
	assert.True(t, filepath.IsAbs(userServiceType.FileName), "FileName should be absolute path")

	// Assert String() method returns the exact source code from subjects.go for UserService interface
	expectedUserServiceSource := `type UserService interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, user *User) error
}`
	assert.Equal(t, expectedUserServiceSource, userServiceType.String(), "UserService String() should return exact source code")

	// Verify Service struct type
	require.NotNil(t, serviceType, "Should find Service type")
	assert.Equal(t, "Service", serviceType.Name)
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/testharness", serviceType.PackagePath())
	assert.Contains(t, serviceType.FileName, "subjects.go")
	assert.True(t, filepath.IsAbs(serviceType.FileName), "FileName should be absolute path")

	// Assert String() method returns the exact source code from subjects.go for Service struct
	expectedServiceSource := `type Service struct {
	userService UserService
}`
	assert.Equal(t, expectedServiceSource, serviceType.String(), "Service String() should return exact source code")
}

func TestScanPackage_TypePackagePath(t *testing.T) {
	// Act
	packageResult := scanHarnessPackage(t)

	// Assert - all types should have correct package path
	expectedPackagePath := "github.com/lonegunmanb/gophon/pkg/testharness"
	for _, typeInfo := range packageResult.Types {
		assert.Equal(t, expectedPackagePath, typeInfo.PackagePath(), "Type %s should have correct package path", typeInfo.Name)
	}
}

// findTypeByName is a helper function that finds a type by name in a slice of TypeInfo
func findTypeByName(types []*TypeInfo, name string) *TypeInfo {
	for i := range types {
		if types[i].Name == name {
			return types[i]
		}
	}
	return nil
}

func TestTypeInfo_IndexFileName(t *testing.T) {
	// Arrange
	tests := []struct {
		name     string
		typeName string
		expected string
	}{
		{
			name:     "Simple struct type",
			typeName: "User",
			expected: "type.User.goindex",
		},
		{
			name:     "Interface type",
			typeName: "UserService",
			expected: "type.UserService.goindex",
		},
		{
			name:     "Type alias",
			typeName: "StringA",
			expected: "type.StringA.goindex",
		},
		{
			name:     "Type alias with equals",
			typeName: "StringB",
			expected: "type.StringB.goindex",
		},
		{
			name:     "Type with mixed case",
			typeName: "HTTPClient",
			expected: "type.HTTPClient.goindex",
		},
		{
			name:     "Type with underscores",
			typeName: "API_Config",
			expected: "type.API_Config.goindex",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			typeInfo := &TypeInfo{
				Name: tt.typeName,
				Range: &Range{
					FileInfo: &FileInfo{
						Package: "github.com/example/pkg",
					},
				},
			}

			// Act
			result := typeInfo.IndexFileName()

			// Assert - verify the type implements IndexableSymbol interface
			var _ IndexableSymbol = typeInfo
			assert.Equal(t, tt.expected, result)
		})
	}
}
