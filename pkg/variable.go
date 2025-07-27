package pkg

// VariableInfo contains information about variable declarations
type VariableInfo struct {
	*Range
	Name        string
	PackagePath string
}
