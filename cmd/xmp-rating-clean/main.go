package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dirPath := flag.String("dir", "", "Path to the directory containing XMP files")
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

	log.Printf("Cleaning up XMP ratings in %v", *dirPath)

	err := filepath.Walk(*dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".xmp") {
			return nil
		}

		return processFile(path, *dryRun)
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nDone!")
}

func processFile(path string, dryRun bool) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", path, err)
	}

	// We look for xmp:Rating="0" and xmp:Rating='0'
	// The user specifically mentioned xmp:Rating="0"
	targets := []string{`xmp:Rating="0"`, `xmp:Rating='0'`}

	modifiedContent := content
	changed := false

	for _, target := range targets {
		if bytes.Contains(modifiedContent, []byte(target)) {
			modifiedContent = bytes.ReplaceAll(modifiedContent, []byte(target), []byte(""))
			changed = true
		}
	}

	if changed {
		if dryRun {
			log.Printf("[DRY RUN] Would clean up %s", path)
		} else {
			log.Printf("Cleaning up %s", path)
			err = ioutil.WriteFile(path, modifiedContent, 0644)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %v", path, err)
			}
		}
	}

	return nil
}
