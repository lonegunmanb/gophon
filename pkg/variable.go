package pkg

import (
	"fmt"
	"go/ast"
)

// VariableInfo contains information about variable declarations
type VariableInfo struct {
	*Range
	*ast.GenDecl
	Name string
}

// IndexFileName generates a predictable index file name for this variable
// Returns a file name in the format: var.<VariableName>.goindex
func (v *VariableInfo) IndexFileName() string {
	return fmt.Sprintf("var.%s.goindex", v.Name)
}
