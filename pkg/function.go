package pkg

import (
	"fmt"
	"strings"
)

// FunctionInfo contains detailed information about functions and methods
type FunctionInfo struct {
	*Range
	Name         string
	ReceiverType string
}

// IndexFileName generates a predictable index file name for this function or method
// For functions: func.<FunctionName>.goindex
// For methods: method.<ReceiverType>.<MethodName>.goindex (strips * from pointer receivers)
func (f *FunctionInfo) IndexFileName() string {
	if f.ReceiverType == "" {
		// Regular function
		return fmt.Sprintf("func.%s.goindex", f.Name)
	} else {
		// Method - clean up receiver type by removing * prefix
		receiverType := strings.TrimPrefix(f.ReceiverType, "*")
		return fmt.Sprintf("method.%s.%s.goindex", receiverType, f.Name)
	}
}
