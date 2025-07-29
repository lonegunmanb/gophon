package pkg

import (
	"io/fs"
	"strings"
	"testing"

	"github.com/prashantv/gostub"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexSourceCode(t *testing.T) {
	mockFs := afero.NewMemMapFs()
	stub := gostub.Stub(&destFs, mockFs)
	defer stub.Reset()

	// Test the index file generator against real pkg/testharness
	require.NoError(t, IndexSourceCode("testharness", "github.com/lonegunmanb/gophon/pkg", "output", nil))

	// Verify that index files were created in the destination filesystem
	// Check for some expected files based on what's actually in pkg/testharness

	// Verify directory structure is maintained
	expectedDirs := []string{
		"output",
		"output/testharness",
	}

	for _, expectedDir := range expectedDirs {
		exists, err := afero.DirExists(destFs, expectedDir)
		require.NoError(t, err)
		assert.True(t, exists, "Expected directory to exist: %s", expectedDir)
	}

	// Check that at least some .goindex files were created
	indexFileCount := 0
	err := afero.Walk(destFs, "output", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".goindex") {
			indexFileCount++

			// Verify file content is not empty
			content, err := afero.ReadFile(destFs, path)
			require.NoError(t, err)
			assert.NotEmpty(t, content, "Expected index file to have content: %s", path)
		}
		return nil
	})
	require.NoError(t, err)

	// We should have at least some index files from the real testharness package
	assert.Greater(t, indexFileCount, 0, "Expected at least one index file to be created")
}

func TestIndexSourceCode_EmptyPackage(t *testing.T) {
	mockFs := afero.NewMemMapFs()

	stubs := gostub.Stub(&destFs, mockFs)
	defer stubs.Reset()

	// Mock ScanPackagesRecursively to return empty package
	stubs.Stub(&ScanPackage, func(pkgPath, basePkgUrl string) (*PackageInfo, error) {
		return &PackageInfo{
			Files:     []*FileInfo{},
			Constants: []*ConstantInfo{},
			Variables: []*VariableInfo{},
			Types:     []*TypeInfo{},
			Functions: []*FunctionInfo{},
		}, nil
	})

	// Test with empty package
	err := IndexSourceCode("empty", "github.com/example/test", "output", nil)
	require.NoError(t, err)

	// Verify output directory does not exists
	exists, err := afero.DirExists(destFs, "output")
	require.NoError(t, err)
	assert.False(t, exists)
}
