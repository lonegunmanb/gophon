package pkg

// ConstantInfo contains information about constant declarations
type ConstantInfo struct {
	*Range
	Name        string
	PackagePath string
}
