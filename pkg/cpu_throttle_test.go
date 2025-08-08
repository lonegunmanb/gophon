package pkg

import (
	"os"
	"runtime"
	"testing"
	"time"
)

func TestGetCPUThrottleConfig(t *testing.T) {
	// Save original environment
	originalCPULimit := os.Getenv("GOPHON_CPU_LIMIT")
	defer func() {
		if originalCPULimit == "" {
			os.Unsetenv("GOPHON_CPU_LIMIT")
		} else {
			os.Setenv("GOPHON_CPU_LIMIT", originalCPULimit)
		}
	}()

	tests := []struct {
		name            string
		envValue        string
		expectedPercent int
		expectedWorkers int
		hasDelay        bool
	}{
		{
			name:            "no environment variable",
			envValue:        "",
			expectedPercent: 100,
			expectedWorkers: runtime.NumCPU(),
			hasDelay:        false,
		},
		{
			name:            "100% CPU limit",
			envValue:        "100",
			expectedPercent: 100,
			expectedWorkers: runtime.NumCPU(),
			hasDelay:        false,
		},
		{
			name:            "50% CPU limit",
			envValue:        "50",
			expectedPercent: 50,
			expectedWorkers: max(1, (runtime.NumCPU()*50)/100),
			hasDelay:        true,
		},
		{
			name:            "10% CPU limit",
			envValue:        "10",
			expectedPercent: 10,
			expectedWorkers: max(1, (runtime.NumCPU()*10)/100),
			hasDelay:        true,
		},
		{
			name:            "invalid value",
			envValue:        "invalid",
			expectedPercent: 100,
			expectedWorkers: runtime.NumCPU(),
			hasDelay:        false,
		},
		{
			name:            "out of range value",
			envValue:        "150",
			expectedPercent: 100,
			expectedWorkers: runtime.NumCPU(),
			hasDelay:        false,
		},
		{
			name:            "zero value",
			envValue:        "0",
			expectedPercent: 100,
			expectedWorkers: runtime.NumCPU(),
			hasDelay:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.envValue == "" {
				os.Unsetenv("GOPHON_CPU_LIMIT")
			} else {
				os.Setenv("GOPHON_CPU_LIMIT", tt.envValue)
			}

			// Get configuration
			config := getCPUThrottleConfig()

			// Verify expected values
			if config.CPULimitPercent != tt.expectedPercent {
				t.Errorf("expected CPU limit %d%%, got %d%%", tt.expectedPercent, config.CPULimitPercent)
			}

			if config.MaxWorkers != tt.expectedWorkers {
				t.Errorf("expected %d workers, got %d", tt.expectedWorkers, config.MaxWorkers)
			}

			if tt.hasDelay && config.WorkerDelay == 0 {
				t.Error("expected worker delay to be set, but it was 0")
			}

			if !tt.hasDelay && config.WorkerDelay != 0 {
				t.Errorf("expected no worker delay, but got %v", config.WorkerDelay)
			}

			// Verify delay calculation for throttled scenarios
			if tt.expectedPercent < 100 {
				expectedDelayMs := (100 - tt.expectedPercent) * 2
				expectedDelay := time.Duration(expectedDelayMs) * time.Millisecond
				if config.WorkerDelay != expectedDelay {
					t.Errorf("expected delay %v, got %v", expectedDelay, config.WorkerDelay)
				}
			}
		})
	}
}

func TestMaxFunction(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{1, 2, 2},
		{5, 3, 5},
		{0, 0, 0},
		{-1, 1, 1},
		{10, 10, 10},
	}

	for _, tt := range tests {
		result := max(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("max(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expected)
		}
	}
}
