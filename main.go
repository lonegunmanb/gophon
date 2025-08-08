package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/lonegunmanb/gophon/pkg"
)

func main() {
	var (
		pkgPath    = flag.String("pkg", "", "Package path to scan (e.g., 'testharness' or '' for root)")
		basePkgUrl = flag.String("base", "", "Base package URL (e.g., 'github.com/lonegunmanb/gophon/pkg')")
		destDir    = flag.String("dest", "./index", "Destination directory for generated index files")
		help       = flag.Bool("help", false, "Show help message")
	)

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "gophon - Go Project Code Indexing Tool\n\n")
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		_, _ = fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		_, _ = fmt.Fprintf(os.Stderr, "  GOPHON_CPU_LIMIT    Limit CPU usage percentage (1-100, default: 100)\n")
		_, _ = fmt.Fprintf(os.Stderr, "                      Lower values reduce CPU usage but increase processing time\n")
		_, _ = fmt.Fprintf(os.Stderr, "                      Useful for CI/CD environments to avoid timeouts\n")
		_, _ = fmt.Fprintf(os.Stderr, "\nExamples:\n")
		_, _ = fmt.Fprintf(os.Stderr, "  # Index the entire project\n")
		_, _ = fmt.Fprintf(os.Stderr, "  %s -base=github.com/lonegunmanb/gophon/pkg -dest=./output\n\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "  # Index a specific package\n")
		_, _ = fmt.Fprintf(os.Stderr, "  %s -pkg=testharness -base=github.com/lonegunmanb/gophon/pkg -dest=./output\n\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "  # Index with CPU throttling (50%% CPU usage)\n")
		_, _ = fmt.Fprintf(os.Stderr, "  GOPHON_CPU_LIMIT=50 %s -base=github.com/lonegunmanb/gophon/pkg -dest=./output\n\n", os.Args[0])
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *basePkgUrl == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Error: -base flag is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Convert destination path to absolute path
	absDestDir, err := filepath.Abs(*destDir)
	if err != nil {
		log.Fatalf("Failed to resolve destination directory: %v", err)
	}

	fmt.Printf("Gophon Code Indexer\n")
	fmt.Printf("===================\n")
	fmt.Printf("Package path: %s\n", *pkgPath)
	fmt.Printf("Base URL: %s\n", *basePkgUrl)
	fmt.Printf("Destination: %s\n", absDestDir)
	
	// Show CPU throttling status
	if cpuLimit := os.Getenv("GOPHON_CPU_LIMIT"); cpuLimit != "" {
		fmt.Printf("CPU Limit: %s%%\n", cpuLimit)
	} else {
		fmt.Printf("CPU Limit: 100%% (no throttling)\n")
	}
	
	fmt.Printf("\nGenerating index files...\n")

	// Track start time for elapsed time and ETA calculations
	startTime := time.Now()
	
	// Create progress callback with rich visual feedback
	progressCallback := func(progress pkg.ProgressInfo) {
		elapsed := time.Since(startTime)
		
		// Calculate ETA
		var eta time.Duration
		if progress.Percentage > 0 {
			totalEstimated := elapsed * time.Duration(100.0/progress.Percentage)
			eta = totalEstimated - elapsed
		}

		// Progress bar visualization
		barWidth := 50
		filled := int(float64(barWidth) * progress.Percentage / 100.0)
		bar := ""
		for i := 0; i < barWidth; i++ {
			if i < filled {
				bar += "█"
			} else {
				bar += "░"
			}
		}

		// Truncate current package path if too long
		current := progress.Current
		if len(current) > 30 {
			current = "..." + current[len(current)-27:]
		}

		// Calculate processing rate
		var rate float64
		if elapsed.Seconds() > 0 {
			rate = float64(progress.Completed) / elapsed.Seconds()
		}

		// Display progress with rich formatting
		fmt.Printf("\r[%s] %.1f%% (%d/%d) | ⏱️ %.1fs", 
			bar, progress.Percentage, progress.Completed, progress.Total, elapsed.Seconds())
		
		if progress.Percentage > 0 && progress.Percentage < 100 {
			fmt.Printf(" | 🔮 ETA: %.1fs", eta.Seconds())
		}
		
		if current != "Completed" {
			fmt.Printf(" | 📦 %s", current)
		}
		
		if rate > 0 {
			fmt.Printf(" | ⚡ %.1f pkg/s", rate)
		}
		
		if progress.Percentage >= 100 {
			fmt.Printf("\n")
		}
	}

	// Call IndexSourceCode with progress callback
	err = pkg.IndexSourceCode(*pkgPath, *basePkgUrl, absDestDir, progressCallback)
	if err != nil {
		log.Fatalf("Failed to generate index files: %v", err)
	}

	// Calculate final statistics
	elapsed := time.Since(startTime)
	
	fmt.Printf("✅ Index generation completed successfully!\n")
	fmt.Printf("📊 Summary:\n")
	fmt.Printf("   • Total time: %.1fs\n", elapsed.Seconds())
	fmt.Printf("   • Output directory: %s\n", absDestDir)
}
