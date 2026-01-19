package main

import (
	"fmt"
)

// ProcessDirectory scans and processes XMP files in the given directory
func ProcessDirectory(dirPath string, dryRun bool) error {
	pairs, err := ScanForXMPPairs(dirPath)
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	soloFiles, err := ScanForSoloXMPFiles(dirPath)
	if err != nil {
		return fmt.Errorf("failed to scan for solo files: %w", err)
	}

	if len(pairs) == 0 && len(soloFiles) == 0 {
		fmt.Println("No XMP files found that need processing")
		return nil
	}

	if len(pairs) > 0 {
		fmt.Printf("Found %d XMP file pair(s) to process\n", len(pairs))
		for _, pair := range pairs {
			if err := MergeXMPFiles(pair, dryRun); err != nil {
				return fmt.Errorf("failed to process XMP files: %w", err)
			}
		}
		fmt.Println()
	}

	if len(soloFiles) > 0 {
		fmt.Printf("Found %d solo XMP file(s) to rename\n", len(soloFiles))
		for _, solo := range soloFiles {
			if err := RenameSoloXMPFile(solo, dryRun); err != nil {
				return fmt.Errorf("failed to rename XMP file: %w", err)
			}
		}
	}

	return nil
}
