package pkg

// ConstantInfo contains information about constant declarations
type ConstantInfo struct {
	Name       string
	Type       string
	Value      string
	FileName   string
	StartLine  int
	EndLine    int
	IsExported bool
	Comments   string
}
