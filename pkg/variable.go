package pkg

// VariableInfo contains information about variable declarations
type VariableInfo struct {
	Name       string
	Type       string
	FileName   string
	StartLine  int
	EndLine    int
	IsExported bool
	Comments   string
}
