package pkg

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestScanPackage_ExtractsFunctions(t *testing.T) {
	// Act
	packageResult := scanHarnessPackage(t)

	// Assert - check that functions are extracted
	assert.Len(t, packageResult.Functions, 5, "Should extract 5 function declarations from test harness")

	// Find specific functions we expect
	newServiceFunc := findFunctionByName(packageResult.Functions, "NewService")
	validateEmailFunc := findFunctionByName(packageResult.Functions, "ValidateEmail")
	containsFunc := findFunctionByName(packageResult.Functions, "contains")

	// Find specific methods we expect
	createUserMethod := findMethodByNameAndReceiver(packageResult.Functions, "CreateUser", "*Service")
	getUserMethod := findMethodByNameAndReceiver(packageResult.Functions, "GetUser", "*Service")

	// Verify NewService function (standalone function)
	require.NotNil(t, newServiceFunc, "Should find NewService function")
	assert.Equal(t, "NewService", newServiceFunc.Name)
	assert.Equal(t, "", newServiceFunc.ReceiverType, "NewService should be a standalone function (no receiver)")
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/test-harness", newServiceFunc.PackagePath)
	assert.Contains(t, newServiceFunc.FileName, "subjects.go")
	assert.True(t, filepath.IsAbs(newServiceFunc.FileName), "FileName should be absolute path")

	// Assert String() method returns the exact source code from subjects.go
	expectedNewServiceSource := `func NewService(userService UserService) *Service {
	return &Service{
		userService: userService,
	}
}`
	assert.Equal(t, expectedNewServiceSource, newServiceFunc.String(), "NewService String() should return exact source code")

	// Verify ValidateEmail function (standalone function with parameters and return values)
	require.NotNil(t, validateEmailFunc, "Should find ValidateEmail function")
	assert.Equal(t, "ValidateEmail", validateEmailFunc.Name)
	assert.Equal(t, "", validateEmailFunc.ReceiverType, "ValidateEmail should be a standalone function (no receiver)")
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/test-harness", validateEmailFunc.PackagePath)
	assert.Contains(t, validateEmailFunc.FileName, "subjects.go")
	assert.True(t, filepath.IsAbs(validateEmailFunc.FileName), "FileName should be absolute path")

	// Assert String() method returns the exact source code from subjects.go
	expectedValidateEmailSource := `func ValidateEmail(email string) bool {
	return len(email) > 0 && contains(email, "@")
}`
	assert.Equal(t, expectedValidateEmailSource, validateEmailFunc.String(), "ValidateEmail String() should return exact source code")

	// Verify contains function (private/unexported function)
	require.NotNil(t, containsFunc, "Should find contains function")
	assert.Equal(t, "contains", containsFunc.Name)
	assert.Equal(t, "", containsFunc.ReceiverType, "contains should be a standalone function (no receiver)")
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/test-harness", containsFunc.PackagePath)
	assert.Contains(t, containsFunc.FileName, "subjects.go")
	assert.True(t, filepath.IsAbs(containsFunc.FileName), "FileName should be absolute path")

	// Assert String() method returns the exact source code from subjects.go
	expectedContainsSource := `func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}`
	assert.Equal(t, expectedContainsSource, containsFunc.String(), "contains String() should return exact source code")

	// Verify CreateUser method (method with receiver)
	require.NotNil(t, createUserMethod, "Should find CreateUser method")
	assert.Equal(t, "CreateUser", createUserMethod.Name)
	assert.Equal(t, "*Service", createUserMethod.ReceiverType, "CreateUser should be a method with *Service receiver")
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/test-harness", createUserMethod.PackagePath)
	assert.Contains(t, createUserMethod.FileName, "subjects.go")
	assert.True(t, filepath.IsAbs(createUserMethod.FileName), "FileName should be absolute path")

	// Assert String() method returns the exact source code from subjects.go
	expectedCreateUserSource := `func (s *Service) CreateUser(ctx context.Context, name, email string) (*User, error) {
	if !ValidateEmail(email) {
		return nil, fmt.Errorf("invalid email: %s", email)
	}

	user := &User{
		Name:  name,
		Email: email,
	}

	return user, s.userService.Create(ctx, user)
}`
	assert.Equal(t, expectedCreateUserSource, createUserMethod.String(), "CreateUser String() should return exact source code")

	// Verify GetUser method (method with different parameter types)
	require.NotNil(t, getUserMethod, "Should find GetUser method")
	assert.Equal(t, "GetUser", getUserMethod.Name)
	assert.Equal(t, "*Service", getUserMethod.ReceiverType, "GetUser should be a method with *Service receiver")
	assert.Equal(t, "github.com/lonegunmanb/gophon/pkg/test-harness", getUserMethod.PackagePath)
	assert.Contains(t, getUserMethod.FileName, "subjects.go")
	assert.True(t, filepath.IsAbs(getUserMethod.FileName), "FileName should be absolute path")

	// Assert String() method returns the exact source code from subjects.go
	expectedGetUserSource := `func (s *Service) GetUser(ctx context.Context, id int64) (*User, error) {
	return s.userService.GetByID(ctx, id)
}`
	assert.Equal(t, expectedGetUserSource, getUserMethod.String(), "GetUser String() should return exact source code")
}

func TestScanPackage_FunctionPackagePath(t *testing.T) {
	// Act
	packageResult := scanHarnessPackage(t)

	// Assert - all functions should have correct package path
	expectedPackagePath := "github.com/lonegunmanb/gophon/pkg/test-harness"
	for _, functionInfo := range packageResult.Functions {
		assert.Equal(t, expectedPackagePath, functionInfo.PackagePath, "Function %s should have correct package path", functionInfo.Name)
	}
}

// findFunctionByName is a helper function that finds a function by name in a slice of FunctionInfo
func findFunctionByName(functions []FunctionInfo, name string) *FunctionInfo {
	for i := range functions {
		if functions[i].Name == name && functions[i].ReceiverType == "" {
			return &functions[i]
		}
	}
	return nil
}

// findMethodByNameAndReceiver is a helper function that finds a method by name and receiver type in a slice of FunctionInfo
func findMethodByNameAndReceiver(functions []FunctionInfo, name string, receiverType string) *FunctionInfo {
	for i := range functions {
		if functions[i].Name == name && functions[i].ReceiverType == receiverType {
			return &functions[i]
		}
	}
	return nil
}
