package pkg

// PackageInfo holds comprehensive information about a scanned package
type PackageInfo struct {
	Files     []FileInfo
	Constants []ConstantInfo
}
