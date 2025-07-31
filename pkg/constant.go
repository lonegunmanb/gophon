package pkg

import (
	"fmt"
	"go/ast"
)

// ConstantInfo contains information about constant declarations
type ConstantInfo struct {
	*Range
	*ast.GenDecl
	Name string
}

// IndexFileName generates a predictable index file name for this constant
// Returns a file name in the format: var.<ConstantName>.goindex
// Uses 'var' prefix to simplify AI agent lookups (same as variables)
func (c *ConstantInfo) IndexFileName() string {
	return fmt.Sprintf("var.%s.goindex", c.Name)
}
