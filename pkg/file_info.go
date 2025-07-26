package pkg

// FileInfo contains information about a single Go file
type FileInfo struct {
	FileName string
	FilePath string
	Package  string
	Imports  string
}
