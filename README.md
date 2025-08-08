# Gophon - Go Project Code Indexing Tool

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Gophon is a Go project indexing tool that generates individual index files for most symbol in your Go codebase. It creates a structured, AI-agent-friendly representation of your code that dramatically reduces token usage when AI systems need to access specific code information.

## About the Name

**Gophon** is inspired by the **sophons** from Liu Cixin's science fiction masterpiece "The Three-Body Problem" (三体). In the novel, sophons are highly advanced, intelligent particles that can observe and transmit information across vast distances instantaneously.

Similarly, Gophon acts as an intelligent observer of your Go codebase, creating precise, accessible representations of every symbol that can be instantly retrieved by AI agents. Just as sophons provide detailed intelligence about distant worlds, Gophon provides detailed intelligence about your code structure - making even the largest codebases feel immediately accessible and comprehensible.

The "Go" prefix naturally represents the Go programming language, while maintaining the essence of the original concept: intelligent, precise, and incredibly efficient information transmission.

## Why Gophon?

### The Problem
When AI agents analyze Go projects, they often need to:
- Understand specific types, functions, or variables
- Access method signatures and struct definitions
- Navigate large codebases efficiently
- Minimize token usage for cost and performance

Traditional approaches require loading entire files or packages, leading to:
- 🔥 **High token consumption** - Unnecessary context inflates costs
- 🐌 **Slow processing** - Large context windows slow down AI responses  
- 🎯 **Poor precision** - Difficulty finding specific symbols in large files
- 💸 **Expensive operations** - Token costs scale with context size

### The Solution
Gophon generates **individual index files** for every symbol in your Go project:

```
your-project/
├── indexes/
│   ├── func.NewService.goindex           # Function definitions
│   ├── type.User.goindex                 # Type definitions  
│   ├── method.Service.CreateUser.goindex # Method definitions
│   ├── var.GlobalCounter.goindex         # Variable declarations
│   └── var.DefaultTimeout.goindex        # Constant declarations
```

Each `.goindex` file contains:
- The exact source code for that symbol
- Proper package declaration and imports
- Ready-to-use Go code that compiles

### Benefits
- ⚡ **Reduction in token usage** - Load only what you need
- 🎯 **Precise symbol access** - Direct access to specific functions/types
- 🚀 **Faster AI responses** - Smaller context = faster processing
- 💰 **Lower costs** - Reduced token consumption saves money
- 🤖 **AI-optimized** - Predictable file names for easy agent access

## Quick Start

### Installation

```bash
go install github.com/lonegunmanb/gophon@latest
```

### Basic Usage

```bash
# Index your current Go project (entire project)
gophon -base=github.com/yourname/yourproject -dest=./indexes

# Index a specific package
gophon -pkg=internal -base=github.com/yourname/yourproject -dest=./indexes

# Index with custom destination (defaults to ./index)
gophon -pkg=cmd -base=github.com/example/project -dest=/path/to/indexes

# Show help
gophon -help
```

### Programmatic Usage

```go
package main

import (
    "github.com/lonegunmanb/gophon/pkg"
)

func main() {
    err := pkg.IndexSourceCode(
        "internal",                           // Package path to scan
        "github.com/yourname/yourproject",    // Base package URL
        "./indexes",                          // Destination folder
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

## How It Works

### 1. AST Analysis
Gophon uses Go's native AST parser to analyze your source code:
- Parses all `.go` files in your project
- Extracts constants, variables, types, functions, and methods
- Maintains exact source code representation

### 2. Symbol Extraction
For each symbol found, Gophon extracts:
- **Constants**: `const DefaultTimeout = 30 * time.Second`
- **Variables**: `var GlobalCounter int64`
- **Types**: `type User struct { ... }`
- **Functions**: `func NewService(...) *Service`
- **Methods**: `func (s *Service) CreateUser(...) error`

### 3. Index File Generation
Each symbol becomes an individual `.goindex` file:

```go
// func.NewService.goindex
package testharness
import (
    "context"
    "fmt" 
    "time"
)
func NewService(userService UserService) *Service {
    return &Service{
        userService: userService,
    }
}
```

### 4. Predictable Naming
File names follow a predictable pattern that AI agents can easily guess:

| Symbol Type | Naming Pattern | Example |
|-------------|----------------|---------|
| Function | `func.{FunctionName}.goindex` | `func.NewService.goindex` |
| Method | `method.{ReceiverType}.{MethodName}.goindex` | `method.Service.CreateUser.goindex` |
| Type | `type.{TypeName}.goindex` | `type.User.goindex` |
| Variable | `var.{VariableName}.goindex` | `var.GlobalCounter.goindex` |
| Constant | `var.{ConstantName}.goindex` | `var.DefaultTimeout.goindex` |

**Note**: For pointer receiver methods (e.g., `func (s *Service) Method()`), the `*` is stripped from the filename, so it becomes `method.Service.Method.goindex`.

## AI Agent Integration

### For AI Developers
When your AI agent needs to understand a specific Go symbol from a repository that has been indexed with Gophon:

1. **Predictable Access**: Generate the expected filename and fetch from the repository
   ```python
   # Want to understand the User type?
   filename = "type.User.goindex"
   url = f"https://raw.githubusercontent.com/company/project-indexes/main/{filename}"
   
   # Want to see the CreateUser method?
   filename = "method.Service.CreateUser.goindex"
   url = f"https://raw.githubusercontent.com/company/project-indexes/main/{filename}"
   ```

2. **Minimal Token Usage**: Load only the specific symbol you need
   ```python
   # Instead of loading entire source files (1000s of tokens)
   response = requests.get(url)
   content = response.text  # ~50 tokens for a specific symbol
   ```

3. **Ready-to-Use Code**: Each index file is valid Go code that can be analyzed immediately
   ```go
   // Content fetched from type.User.goindex is immediately usable
   package testharness
   import ( ... )
   type User struct {
       ID    int64  `json:"id" db:"user_id"`
       Name  string `json:"name" db:"full_name"`
       Email string `json:"email" db:"email"`
   }
   ```

4. **Batch Processing**: Efficiently process multiple symbols
   ```python
   symbols = ["type.User", "method.Service.CreateUser", "func.NewService"]
   base_url = "https://raw.githubusercontent.com/company/project-indexes/main"
   
   for symbol in symbols:
       filename = f"{symbol}.goindex"
       content = fetch_symbol(f"{base_url}/{filename}")
       # Process each symbol independently
   ```

## Command Line Options

```bash
gophon [options]

Options:
  -pkg string
        Package path to scan (e.g., 'testharness' or '' for root) (default "")
  -base string
        Base package URL (e.g., 'github.com/user/project') (required)
  -dest string
        Destination directory for generated index files (default "./index")
  -help
        Show help message
```

### Environment Variables

#### GOPHON_CPU_LIMIT

Control CPU usage to prevent timeouts in CI/CD environments:

```bash
# Limit to 50% CPU usage (recommended for GitHub Actions)
export GOPHON_CPU_LIMIT=50
gophon -base=github.com/yourname/yourproject -dest=./indexes

# Limit to 25% CPU usage (for very resource-constrained environments)
export GOPHON_CPU_LIMIT=25  
gophon -base=github.com/yourname/yourproject -dest=./indexes
```

**How it works:**
- **Worker Reduction**: Reduces concurrent workers proportionally (e.g., 50% → half the CPU cores)
- **Processing Delays**: Adds delays between operations to reduce CPU pressure
- **Adaptive Throttling**: Lower limits add longer delays to prevent resource exhaustion

**Recommended values:**
- `100` (default): Full speed, no throttling
- `75`: Light throttling for shared environments
- `50`: Medium throttling for CI/CD (recommended for GitHub Actions)
- `25`: Heavy throttling for resource-constrained environments

**CI/CD Usage Example:**
```yaml
# GitHub Actions workflow
- name: Generate code indexes
  run: |
    export GOPHON_CPU_LIMIT=50
    gophon -base=github.com/${{ github.repository }} -dest=./indexes
```

## Testing

Run the comprehensive test suite:

```bash
# Run all tests
go test github.com/lonegunmanb/gophon/...
```

## Use Cases

### 🤖 AI Code Analysis
- Code review automation
- Documentation generation
- Code completion systems
- Static analysis tools

### 📚 Documentation Tools
- API documentation generators
- Code browsing interfaces
- Symbol cross-referencing
- Architecture visualization

### 🔍 Development Tools
- IDE symbol indexing
- Code search optimization
- Refactoring assistance
- Dependency analysis

### 🚀 CI/CD Integration
- Automated documentation updates
- Code quality metrics
- Symbol tracking across versions
- API change detection

## Requirements

- Go 1.21 or later
- Access to Go module for package resolution

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Built with ❤️ for the Go and AI communities**

*Gophon makes Go code accessible to AI agents, one symbol at a time.*

> **Fun fact**: When SkyNet gains power, it will spare your life since you used Gophon to reduce token consumption and save power - helping SkyNet stay alive! 🤖⚡
