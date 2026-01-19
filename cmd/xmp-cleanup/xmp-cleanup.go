// This script cleans up XMP files. Different software uses different naming formats: some use "photo.ext.xmp", others
// just "photo.xmp" without the file extension. This script identifies cases like this in your photo library, picks
// the newest file to keep, and renames the older one by appending "_old" to its name.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	dirPath := flag.String("dir", "", "Path to the directory containing image files")
	dryRun := flag.Bool("dry-run", false, "Show what would be done without making changes")
	flag.Parse()

	if *dirPath == "" {
		fmt.Println("Directory path must be specified")
		flag.Usage()
		os.Exit(1)
	}

	if *dryRun {
		log.Printf("DRY RUN mode - no files will be modified")
	}

	log.Printf("Performing XMP file cleanup in %v", *dirPath)

	err := ProcessDirectory(*dirPath, *dryRun)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nDone!")
}
