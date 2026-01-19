package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// RenameSoloXMPFile renames a solo XMP file from "photo.ext.xmp" to "photo.xmp"
func RenameSoloXMPFile(solo SoloXMPFile, dryRun bool) error {
	if dryRun {
		fmt.Printf("[DRY RUN] Would rename:\n")
		fmt.Printf("  %s -> %s\n", solo.Current, solo.Target)
		return nil
	}

	if err := os.Rename(solo.Current, solo.Target); err != nil {
		return fmt.Errorf("failed to rename %s to %s: %w", solo.Current, solo.Target, err)
	}

	color.Green("Renamed: %s -> %s\n", solo.Current, solo.Target)
	return nil
}
