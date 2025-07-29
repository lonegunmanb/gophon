package pkg

import "strings"

// Range represents a range of lines within a file, indicating start and end line numbers.
// It embeds *FileInfo to provide access to file information.
type Range struct {
	*FileInfo
	StartLine int // 1-based line number (inclusive)
	EndLine   int // 1-based line number (inclusive)
}

// String returns the lines from StartLine to EndLine (inclusive, 1-based) from the file.
// If the range is invalid or the file cannot be read, returns an empty string.
func (r *Range) String() string {
	if r.FileInfo == nil {
		return ""
	}

	// Get the complete file content
	content := r.FileInfo.String()
	if content == "" {
		return ""
	}

	// Split content into lines
	lines := strings.Split(content, "\n")

	// Validate line range (convert to 0-based indexing)
	startIdx := r.StartLine - 1
	endIdx := r.EndLine - 1

	if startIdx < 0 || endIdx < 0 || startIdx >= len(lines) || endIdx >= len(lines) || startIdx > endIdx {
		return ""
	}

	// Extract the specified lines
	selectedLines := lines[startIdx : endIdx+1]

	return strings.ReplaceAll(strings.Join(selectedLines, "\n"), "\r", "")
}
