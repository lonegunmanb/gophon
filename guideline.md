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
- [ ] **Package Index Function**: Create empty function signature for package indexing (implementation pending)
- [ ] **Import Analysis**: Track package dependencies and import paths

### 1.2 Reverse Index File Generation
- [ ] **Individual Symbol Files**: Generate standalone files for each method, type, function, etc.
- [ ] **AI-Friendly Format**: Structure files with package path, imports, and source code in easily parsable format
- [ ] **Smart File Naming**: Use descriptive names like `method_PackageName_TypeName_MethodName.go`, `type_PackageName_TypeName.go`
- [ ] **Complete Context**: Include necessary imports and type definitions for each symbol's context
- [ ] **Dependency Resolution**: Automatically include referenced types and their definitions
- [ ] **Cross-Reference Links**: Add metadata about related symbols and their file locations
- [ ] **Usage Examples**: Include code examples showing how to use each symbol
- [ ] **Performance Metrics**: Track file sizes and generation times for optimization

### 1.3 Recursive Package Discovery
- [ ] **Directory Traversal**: Recursively scan all packages in the current module only
- [ ] **Go Module Support**: Handle go.mod files and respect module boundaries
- [ ] **Current Module Scope**: Only index packages within the current module, exclude external dependencies
- [ ] **Build Constraint Awareness**: Respect build tags and constraints within the current module
- [ ] **Vendor Directory Handling**: Properly handle vendor directories and their exclusion
- [ ] **Symlink Resolution**: Handle symbolic links in package directories

## Phase 2: Data Storage and Indexing

### 2.1 Index Data Structure Design
- [ ] **Go Code Format**: Design index files as valid Go code instead of JSON to minimize markup tokens
- [ ] **Hierarchical Organization**: Package -> File -> Symbol structure using Go struct definitions
- [ ] **Metadata Storage**: Git commit hash, timestamp, Go version compatibility as Go constants/variables
- [ ] **Search Optimization**: Maintain original project directory structure for index files to enable direct GitHub raw URL access
- [ ] **Incremental Updates**: Support for updating only changed files since last indexing
- [ ] **Compression Strategy**: Implement efficient storage for large codebases

### 2.2 Query and Retrieval System
- [ ] **Symbol Search**: Fast lookup by symbol name, type, or package
- [ ] **Fuzzy Matching**: Support approximate symbol name matching
- [ ] **Dependency Graph**: Build and maintain symbol dependency relationships
- [ ] **Usage Tracking**: Track where symbols are used across the codebase
- [ ] **API Compatibility**: Detect breaking changes between versions

## Phase 3: Advanced Features

### 3.1 Code Analysis and Insights
- [ ] **Complexity Metrics**: Calculate cyclomatic complexity for functions
- [ ] **Code Coverage Integration**: Link with test coverage data
- [ ] **Performance Profiling**: Integrate with pprof data for hot paths
- [ ] **Security Analysis**: Identify potential security issues and patterns
- [ ] **Code Quality Metrics**: Track technical debt and code smells

### 3.2 Integration and Export
- [ ] **MCP Server Implementation**: Provide remote access through Model Context Protocol
- [ ] **REST API**: HTTP endpoints for programmatic access
- [ ] **CLI Tools**: Command-line utilities for local indexing and querying
- [ ] **IDE Integration**: Plugins for popular Go IDEs
- [ ] **CI/CD Integration**: Automated indexing in build pipelines

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
