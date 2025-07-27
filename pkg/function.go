package pkg

// FunctionInfo contains detailed information about functions and methods
type FunctionInfo struct {
	*Range
	Name         string
	ReceiverType string // for methods, empty for functions
	PackagePath  string
}
