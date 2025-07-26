package pkg

// FunctionInfo contains detailed information about functions and methods
type FunctionInfo struct {
	Name         string
	FileName     string
	StartLine    int
	EndLine      int
	ReceiverType string // for methods, empty for functions
}
