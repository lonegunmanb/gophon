package pkg

// PackageInfo holds comprehensive information about a scanned package
type PackageInfo struct {
	PackageName string
	PackageID   string
	Files       []FileInfo
	Functions   []FunctionInfo
	Types       []TypeInfo
	Variables   []VariableInfo
	Constants   []ConstantInfo
}

// FileInfo contains information about a single Go file
type FileInfo struct {
	FileName string
	FilePath string
	Package  string
	Imports  string
}

// FunctionInfo contains detailed information about functions and methods
type FunctionInfo struct {
	Name         string
	FileName     string
	StartLine    int
	EndLine      int
	ReceiverType string // for methods, empty for functions
}

// TypeInfo contains information about type declarations
type TypeInfo struct {
	Name      string
	FileName  string
	StartLine int
	EndLine   int
}

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

// ScanPackage scans the specified package and returns comprehensive information
func ScanPackage(pkgPath string) (*PackageInfo, error) {
	panic("implement me")
}
