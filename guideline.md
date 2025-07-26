# Go Project Remote Code Indexing Tool - Development Plan

## Project Overview
Build a comprehensive Go project indexing system that provides remote access to code definitions, types, methods, and functions through MCP (Model Context Protocol). This tool aims to reduce token usage and context noise when AI agents need to access specific code information from remote repositories.

## Phase 1: Core Indexing Engine Enhancement

### 1.1 Extend AST Analysis Capabilities
- [ ] **Function Definitions**: Extract standalone functions (non-method functions)
- [ ] **Method Definitions**: Extract methods with receivers (already implemented in current code)
- [ ] **Variable/Constant Declarations**: Parse global and package-level variables and constants
- [ ] **Interface Definitions**: Extract interface types and their method signatures
- [ ] **Struct Field Details**: Include struct field names, types, and tags
- [ ] **Import Analysis**: Track package dependencies and import paths
- [ ] **Documentation Extraction**: Parse and store Go doc comments

### 1.2 Reverse Index File Generation
- [ ] **Individual Symbol Files**: Generate standalone files for each method, type, function, etc.
- [ ] **AI-Friendly Format**: Structure files with package path, imports, and source code in easily parsable format
- [ ] **Smart File Naming**: Use descriptive names like `method_PackageName_TypeName_MethodName.go`, `type_PackageName_TypeName.go`
- [ ] **Complete Context**: Include necessary imports and type definitions for each symbol's context
- [ ] **Dependency Resolution**: Automatically include referenced types and their definitions
- [ ] **Cross-Reference Links**: Add metadata about related symbols and their file locations

### 1.3 Recursive Package Discovery
- [ ] **Directory Traversal**: Recursively scan all packages in the current module only
- [ ] **Go Module Support**: Handle go.mod files and respect module boundaries
- [ ] **Current Module Scope**: Only index packages within the current module, exclude external dependencies
- [ ] **Build Constraint Awareness**: Respect build tags and constraints within the current module

## Phase 2: Data Storage and Indexing

### 2.1 Index Data Structure Design
- [ ] **Go Code Format**: Design index files as valid Go code instead of JSON to minimize markup tokens
- [ ] **Hierarchical Organization**: Package -> File -> Symbol structure using Go struct definitions
- [ ] **Metadata Storage**: Git commit hash, timestamp, Go version compatibility as Go constants/variables
- [ ] **Search Optimization**: Maintain original project directory structure for index files to enable direct GitHub raw URL access

## Implementation Priority
1. **Phase 1**: Essential for basic functionality
2. **Phase 2**: Required for practical usage

## Technical Stack
- **Language**: Go (for performance and Go-native AST parsing)
- **Storage**: File system with structured Go code files
- **Dependencies**: golang.org/x/tools/go/packages, go/ast, go/token

## Success Criteria
- [ ] Reduce token usage by 70% for AI code queries
- [ ] Index processing time under 30 seconds for medium-sized projects
- [ ] Generate individual symbol files in original project structure
- [ ] Files readable as valid Go code without JSON markup overhead
