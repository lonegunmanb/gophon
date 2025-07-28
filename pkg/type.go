package pkg

import "fmt"

// TypeInfo contains information about type declarations
type TypeInfo struct {
	*Range
	Name        string
	PackagePath string
}

// IndexFileName generates a predictable index file name for this type
// Returns a file name in the format: type.<TypeName>.goindex
func (t TypeInfo) IndexFileName() string {
	return fmt.Sprintf("type.%s.goindex", t.Name)
}
