package pkg

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileInfo_Imports(t *testing.T) {
	testCases := []struct {
		name            string
		goFileContent   string
		expectedImports string
	}{
		{
			name: "multiple imports",
			goFileContent: `package main

import (
	"fmt"
	"os"
	"strings"

    "github.com/lonegunmanb/gophon/pkg"
)

func main() {
	fmt.Println("Hello, World!")
}
`,
			expectedImports: `import (
	"fmt"
	"os"
	"strings"

    "github.com/lonegunmanb/gophon/pkg"
)`,
		},
		{
			name: "single import",
			goFileContent: `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}
`,
			expectedImports: `import "fmt"`,
		},
		{
			name: "no imports",
			goFileContent: `package main

func main() {
	println("Hello")
}
`,
			expectedImports: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create FileInfo instance
			fileInfo := &FileInfo{
				FileName: "test.go",
				FilePath: "/path/to/test.go",
				Package:  "main",
				content:  &tc.goFileContent,
			}

			// Call Imports() method and assert the result
			actualImports := fileInfo.Imports()

			assert.Equal(t, tc.expectedImports, actualImports, "Imports() should return the expected import statements")
		})
	}
}

func TestFileInfo_String(t *testing.T) {
	// Get current working directory to build absolute path
	currentDir, err := os.Getwd()
	require.NoError(t, err)

	// Create absolute path to subjects.go
	subjectsPath := filepath.Join(currentDir, "test-harness", "subjects.go")

	// Create a new FileInfo directly with the absolute file path
	directFileInfo := &FileInfo{
		FileName: subjectsPath,
	}

	// Read file content directly
	expectedContent, err := os.ReadFile(subjectsPath)
	require.NoError(t, err)

	// Assert - compare String() result with direct file reading
	content := directFileInfo.String()
	assert.Equal(t, string(expectedContent), content, "String() method should return the same content as direct file reading")
}
