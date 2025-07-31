package pkg

import (
	"fmt"
	"go/ast"
)

// TypeInfo contains information about type declarations
type TypeInfo struct {
	*Range
	*ast.GenDecl
	Name string
}

// IndexFileName generates a predictable index file name for this type
// Returns a file name in the format: type.<TypeName>.goindex
func (t *TypeInfo) IndexFileName() string {
	return fmt.Sprintf("type.%s.goindex", t.Name)
}
