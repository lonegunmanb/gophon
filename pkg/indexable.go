package pkg

// IndexableSymbol represents a Go symbol that can generate a predictable index file name
// for AI agents to easily guess and access when reading Go source code.
type IndexableSymbol interface {
	// IndexFileName generates a predictable, unique file name for this symbol
	// that follows the pattern: <type>_<PackageName>_<SymbolName>.go
	// Examples:
	//   - type_mypackage_User.go
	//   - func_mypackage_ValidateEmail.go
	//   - method_mypackage_Service_CreateUser.go
	//   - var_mypackage_GlobalCounter.go
	//   - const_mypackage_DefaultTimeout.go
	IndexFileName() string
}
