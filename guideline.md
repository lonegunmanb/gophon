# Go Project Remote Code Indexing Tool - Development Plan

## Project Overview
Build a comprehensive Go project indexing system that provides remote access to code definitions, types, methods, and functions through MCP (Model Context Protocol). This tool aims to reduce token usage and context noise when AI agents need to access specific code information from remote repositories.

## Phase 1: Core Indexing Engine Enhancement

### 1.1 Extend AST Analysis Capabilities
- [x] **Test Harness Design**: Create test subjects with examples of all Go constructs for indexing validation
- [x] **Data Structure Design**: Define core data structures for storing indexed Go code information
- [x] **Range Type Implementation**: Added Range type with embedded *FileInfo and String() method for line extraction
- [x] **Constant Declarations**: Parse and extract package-level constants with proper Range information
- [x] **Variable Declarations**: Parse global and package-level variables (constants completed)
- [x] **Type Definitions**: Extract all type definitions including structs (with field names, types, and tags), interfaces (with method signatures), type aliases, and custom types
- [x] **Function Definitions**: Extract standalone functions (non-method functions)
- [x] **IndexableSymbol Interface**: Design and implement a new interface with `IndexFileName()` method that `FunctionInfo`, `TypeInfo`, `VariableInfo`, and `ConstantInfo` implement. This interface generates predictable, unique file names that AI agents can easily guess when reading Go source code (e.g., `type.TypeName.goindex`, `func.FunctionName.goindex`, `method.ReceiverType.MethodName.goindex`, `var.VariableName.goindex`). Constants use `var.` prefix to unify with variables for simplified AI agent lookups.
- [x] **Recursive Package Scanner**: Implement `ScanPackagesRecursively` function that accepts `pkgPath`, `basePkgUrl` parameters and a callback function. The callback receives `*PackageInfo` and `pkgUrl` string parameters. This function will scan the specified package using `ScanSinglePackage`, invoke the callback with results, then recursively scan all subdirectories containing Go packages within the current module scope.
- [x] **Index File Generator**: Implement `IndexSourceCode` function that accepts `pkgPath`, `basePkgUrl`, and `destFolder` parameters. Uses `ScanPackagesRecursively` with a callback to write all indexable symbols to individual `.goindex` files in the destination folder, maintaining the package directory structure with secure file permissions (0700/0600).
- [ ] **Command Line Tool**: Modify `main.go` to wrap gophon as a command line tool with flags for source directory, destination directory, and package URL. Provide user-friendly interface for generating index files from Go projects.

## Implementation Priority
1. **Phase 1**: Essential for basic functionality
2. **Phase 2**: Required for practical usage
3. **Phase 3**: Advanced features for enhanced user experience

## Technical Stack
- **Language**: Go (for performance and Go-native AST parsing)
- **Storage**: File system with structured Go code files
- **Dependencies**: golang.org/x/tools/go/packages, go/ast, go/token
- **Testing**: Comprehensive test suite with real-world Go projects
- **Documentation**: Detailed API documentation and usage examples

## Success Criteria
- [ ] Reduce token usage by 70% for AI code queries
- [ ] Index processing time under 30 seconds for medium-sized projects
- [ ] Generate individual symbol files in original project structure
