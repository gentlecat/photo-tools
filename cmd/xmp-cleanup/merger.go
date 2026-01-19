package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

// MergeXMPFiles picks the newest XMP file and renames the older one with "_old" suffix
func MergeXMPFiles(pair XMPPair, dryRun bool) error {
	withExtStat, err := os.Stat(pair.WithExt)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", pair.WithExt, err)
	}

	withoutExtStat, err := os.Stat(pair.WithoutExt)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", pair.WithoutExt, err)
	}

	withExtTime := withExtStat.ModTime()
	withoutExtTime := withoutExtStat.ModTime()

	// Determine which file is newer
	var newerFile, olderFile string
	var newerTime, olderTime time.Time

	if withExtTime.After(withoutExtTime) {
		newerFile = pair.WithExt
		olderFile = pair.WithoutExt
		newerTime = withExtTime
		olderTime = withoutExtTime
	} else {
		newerFile = pair.WithoutExt
		olderFile = pair.WithExt
		newerTime = withoutExtTime
		olderTime = withExtTime
	}

	// Generate the "_old" filename for the older file
	oldFileRenamed := generateOldFileName(olderFile)

	if dryRun {
		fmt.Printf("[DRY RUN] Would process:\n")
		fmt.Printf("  Newer file (keep): %s (modified %s)\n",
			newerFile,
			newerTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Older file (rename): %s (modified %s)\n",
			olderFile,
			olderTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Rename to: %s\n", oldFileRenamed)
		return nil
	}

	// If the newer file is the one with extension, we need to:
	// 1. Rename the older file (without extension) to _old
	// 2. Rename the newer file (with extension) to without extension
	if newerFile == pair.WithExt {
		// Rename older file first
		if err := os.Rename(olderFile, oldFileRenamed); err != nil {
			return fmt.Errorf("failed to rename %s to %s: %w", olderFile, oldFileRenamed, err)
		}

		// Now rename the newer file to take the place
		targetFile := pair.WithoutExt
		if err := os.Rename(newerFile, targetFile); err != nil {
			// Try to restore the older file
			os.Rename(oldFileRenamed, olderFile)
			return fmt.Errorf("failed to rename %s to %s: %w", newerFile, targetFile, err)
		}

		fmt.Printf("Processed: %s, %s\n",
			color.GreenString("%s (kept as %s)", newerFile, targetFile),
			color.RedString("%s (renamed to %s)", olderFile, oldFileRenamed))
	} else {
		// Newer file is already without extension, just rename the older one
		if err := os.Rename(olderFile, oldFileRenamed); err != nil {
			return fmt.Errorf("failed to rename %s to %s: %w", olderFile, oldFileRenamed, err)
		}

		fmt.Printf("Processed: %s, %s\n",
			color.GreenString("%s (kept)", newerFile),
			color.RedString("%s (renamed to %s)", olderFile, oldFileRenamed))
	}

	return nil
}

// generateOldFileName creates a filename with "_old" suffix after the .xmp extension
func generateOldFileName(filepath string) string {
	// Handle case where file already has "_old" (shouldn't happen, but be safe)
	if strings.HasSuffix(filepath, ".xmp_old") {
		return filepath + "_old"
	}

	return filepath + "_old"
}
