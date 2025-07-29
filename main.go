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
		_, _ = fmt.Fprintf(os.Stderr, "\nExamples:\n")
		_, _ = fmt.Fprintf(os.Stderr, "  # Index the entire project\n")
		_, _ = fmt.Fprintf(os.Stderr, "  %s -base=github.com/lonegunmanb/gophon/pkg -dest=./output\n\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "  # Index a specific package\n")
		_, _ = fmt.Fprintf(os.Stderr, "  %s -pkg=testharness -base=github.com/lonegunmanb/gophon/pkg -dest=./output\n\n", os.Args[0])
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
	fmt.Printf("\n🚀 Starting index generation...\n")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	startTime := time.Now()
	packageCount := 0

	// Progress callback with detailed information
	progressCallback := func(progress pkg.ProgressInfo) {
		elapsed := time.Since(startTime)
		
		// Calculate estimated time remaining
		var eta time.Duration
		if progress.Percentage > 0 {
			totalEstimated := elapsed * time.Duration(100.0/progress.Percentage)
			eta = totalEstimated - elapsed
		}

		// Progress bar
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

		// Clear line and show detailed progress
		fmt.Printf("\r[%s] %.1f%% (%d/%d)", 
			bar, progress.Percentage, progress.Completed, progress.Total)
		
		if progress.Percentage < 100.0 {
			fmt.Printf(" | ⏱️  %v", elapsed.Round(time.Millisecond))
			if eta > 0 && eta < time.Hour {
				fmt.Printf(" | 🔮 ETA: %v", eta.Round(time.Second))
			}
			
			// Truncate package name if too long
			currentPkg := progress.Current
			if len(currentPkg) > 50 {
				currentPkg = "..." + currentPkg[len(currentPkg)-47:]
			}
			fmt.Printf(" | 📦 %s", currentPkg)
		}

		// Show processing rate
		if progress.Completed > 0 && elapsed > 0 {
			rate := float64(progress.Completed) / elapsed.Seconds()
			fmt.Printf(" | ⚡ %.1f pkg/s", rate)
		}

		// Move to next line when complete
		if progress.Percentage >= 100.0 {
			fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
			packageCount = progress.Completed
		}
	}

	// Call IndexSourceCode with progress tracking
	err = pkg.IndexSourceCode(*pkgPath, *basePkgUrl, absDestDir, progressCallback)
	if err != nil {
		log.Fatalf("❌ Failed to generate index files: %v", err)
	}

	// Final summary
	totalTime := time.Since(startTime)
	avgRate := float64(packageCount) / totalTime.Seconds()

	fmt.Printf("✅ Index generation completed successfully!\n")
	fmt.Printf("📊 Summary:\n")
	fmt.Printf("   • Total packages indexed: %d\n", packageCount)
	fmt.Printf("   • Total time: %v\n", totalTime.Round(time.Millisecond))
	fmt.Printf("   • Average rate: %.2f packages/second\n", avgRate)
	fmt.Printf("   • Output directory: %s\n", absDestDir)
}
