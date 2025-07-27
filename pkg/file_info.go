package pkg

import (
	"os"
	"sync"
)

// FileInfo contains information about a single Go file
type FileInfo struct {
	FileName string
	FilePath string
	Package  string
	Imports  string
	content  *string      // cached file content
	mu       sync.RWMutex // mutex for thread-safe cache access
}

// String reads and returns the content of the file
func (f *FileInfo) String() string {
	// Check cache with read lock
	f.mu.RLock()
	if f.content != nil {
		cached := *f.content
		f.mu.RUnlock()
		return cached
	}
	f.mu.RUnlock()

	// Acquire write lock for file reading and caching
	f.mu.Lock()
	defer f.mu.Unlock()

	// Double-check pattern: another goroutine might have cached it
	if f.content != nil {
		return *f.content
	}

	// Read file content
	contentBytes, err := os.ReadFile(f.FileName)
	if err != nil {
		return ""
	}

	// Cache the content
	contentStr := string(contentBytes)
	f.content = &contentStr

	return contentStr
}
