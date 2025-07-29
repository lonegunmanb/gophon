package pkg

import "fmt"

// VariableInfo contains information about variable declarations
type VariableInfo struct {
	*Range
	Name string
}

// IndexFileName generates a predictable index file name for this variable
// Returns a file name in the format: var.<VariableName>.goindex
func (v VariableInfo) IndexFileName() string {
	return fmt.Sprintf("var.%s.goindex", v.Name)
}
