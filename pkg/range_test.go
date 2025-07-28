package pkg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRange_String_WithCachedContent(t *testing.T) {
	testContent := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
	fmt.Println("Line 7")
	fmt.Println("Line 8")
	fmt.Println("Line 9")
	fmt.Println("Line 10")
}`

	testCases := []struct {
		name      string
		startLine int
		endLine   int
		expected  string
	}{
		{
			name:      "single line - first line",
			startLine: 1,
			endLine:   1,
			expected:  "package main",
		},
		{
			name:      "single line - middle line",
			startLine: 5,
			endLine:   5,
			expected:  "func main() {",
		},
		{
			name:      "multiple lines - import section",
			startLine: 2,
			endLine:   4,
			expected: `
import "fmt"
`,
		},
		{
			name:      "multiple lines - function body",
			startLine: 6,
			endLine:   8,
			expected: `	fmt.Println("Hello, World!")
	fmt.Println("Line 7")
	fmt.Println("Line 8")`,
		},
		{
			name:      "entire content",
			startLine: 1,
			endLine:   11,
			expected:  testContent,
		},
		{
			name:      "last line only",
			startLine: 11,
			endLine:   11,
			expected:  "}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create FileInfo with cached content
			fileInfo := &FileInfo{
				FileName: "test.go",
				FilePath: "/path/to/test.go",
				Package:  "main",
				content:  &testContent,
			}

			// Create Range instance
			range_ := &Range{
				FileInfo:  fileInfo,
				StartLine: tc.startLine,
				EndLine:   tc.endLine,
			}

			// Test String() method
			result := range_.String()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRange_String_WithActualFile(t *testing.T) {
	// Get current working directory to build absolute path
	currentDir, err := os.Getwd()
	require.NoError(t, err)

	// Create absolute path to subjects.go
	subjectsPath := filepath.Join(currentDir, "testharness", "subjects.go")

	// Verify the file exists
	_, err = os.Stat(subjectsPath)
	require.NoError(t, err, "subjects.go should exist for this test")

	// Create FileInfo with actual file path
	fileInfo := &FileInfo{
		FileName: subjectsPath,
	}

	// Read the actual file content to determine expected results
	contentBytes, err := os.ReadFile(subjectsPath)
	require.NoError(t, err)
	content := string(contentBytes)

	// Split content into lines using the same logic as Range.String()
	lines := strings.Split(content, "\n")

	if len(lines) >= 3 {
		testCases := []struct {
			name      string
			startLine int
			endLine   int
		}{
			{
				name:      "first line of actual file",
				startLine: 1,
				endLine:   1,
			},
			{
				name:      "first three lines",
				startLine: 1,
				endLine:   3,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				range_ := &Range{
					FileInfo:  fileInfo,
					StartLine: tc.startLine,
					EndLine:   tc.endLine,
				}

				result := range_.String()
				assert.NotEmpty(t, result, "Should return non-empty content for valid line range")

				// Verify that the result contains the expected number of lines
				resultLines := strings.Split(result, "\n")
				expectedLineCount := tc.endLine - tc.startLine + 1
				assert.Equal(t, expectedLineCount, len(resultLines), "Should return correct number of lines")

				// Verify the actual content matches the expected lines from the file
				expectedLines := lines[tc.startLine-1 : tc.endLine]
				expectedContent := strings.Join(expectedLines, "\n")
				assert.Equal(t, expectedContent, result, "Should return the exact lines from the file")
			})
		}
	}
}

func TestRange_String_EdgeCases(t *testing.T) {
	testContent := `line 1
line 2
line 3`

	testCases := []struct {
		name      string
		fileInfo  *FileInfo
		startLine int
		endLine   int
		expected  string
	}{
		{
			name:      "nil FileInfo",
			fileInfo:  nil,
			startLine: 1,
			endLine:   1,
			expected:  "",
		},
		{
			name: "empty file content",
			fileInfo: &FileInfo{
				FileName: "empty.go",
				content:  new(string), // empty string
			},
			startLine: 1,
			endLine:   1,
			expected:  "",
		},
		{
			name: "start line greater than end line",
			fileInfo: &FileInfo{
				FileName: "test.go",
				content:  &testContent,
			},
			startLine: 3,
			endLine:   1,
			expected:  "",
		},
		{
			name: "start line out of bounds (too high)",
			fileInfo: &FileInfo{
				FileName: "test.go",
				content:  &testContent,
			},
			startLine: 10,
			endLine:   10,
			expected:  "",
		},
		{
			name: "end line out of bounds (too high)",
			fileInfo: &FileInfo{
				FileName: "test.go",
				content:  &testContent,
			},
			startLine: 2,
			endLine:   10,
			expected:  "",
		},
		{
			name: "negative start line",
			fileInfo: &FileInfo{
				FileName: "test.go",
				content:  &testContent,
			},
			startLine: -1,
			endLine:   1,
			expected:  "",
		},
		{
			name: "zero start line",
			fileInfo: &FileInfo{
				FileName: "test.go",
				content:  &testContent,
			},
			startLine: 0,
			endLine:   1,
			expected:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			range_ := &Range{
				FileInfo:  tc.fileInfo,
				StartLine: tc.startLine,
				EndLine:   tc.endLine,
			}

			result := range_.String()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRange_EmbeddedFileInfo(t *testing.T) {
	// Test that Range properly embeds FileInfo and can access its fields and methods
	testContent := "package test\n\nfunc TestFunc() {}"

	fileInfo := &FileInfo{
		FileName: "embedded_test.go",
		FilePath: "/test/embedded_test.go",
		Package:  "testpkg",
		content:  &testContent,
	}

	range_ := &Range{
		FileInfo:  fileInfo,
		StartLine: 1,
		EndLine:   2,
	}

	// Test that we can access FileInfo fields through the embedded struct
	assert.Equal(t, "embedded_test.go", range_.FileName)
	assert.Equal(t, "/test/embedded_test.go", range_.FilePath)
	assert.Equal(t, "testpkg", range_.Package)

	// Test that we can call FileInfo methods through the embedded struct
	assert.NotPanics(t, func() {
		_ = range_.Imports() // Should not panic
	})

	// Test the Range's String method returns correct lines
	expected := "package test\n"
	assert.Equal(t, expected, range_.String())
}
