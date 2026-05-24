package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	dirPath := flag.String("dir", "", "Path to the directory containing XMP files")
	dryRun := flag.Bool("dry-run", false, "Show what would be done without making changes")
	workers := flag.Int("workers", 4, "Number of parallel workers")
	flag.Parse()

	if *dirPath == "" {
		fmt.Println("Directory path must be specified")
		flag.Usage()
		os.Exit(1)
	}

	if *dryRun {
		log.Printf("DRY RUN mode - no files will be modified")
	}

	log.Printf("Cleaning up XMP ratings in %v (workers: %d)", *dirPath, *workers)

	paths := make(chan string, *workers)
	errs := make(chan error, 1)

	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range paths {
				if err := processFile(path, *dryRun); err != nil {
					select {
					case errs <- err:
					default:
					}
				}
			}
		}()
	}

	walkErr := filepath.Walk(*dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".xmp") {
			return nil
		}
		paths <- path
		return nil
	})

	close(paths)
	wg.Wait()

	if walkErr != nil {
		log.Fatal(walkErr)
	}

	select {
	case err := <-errs:
		log.Fatal(err)
	default:
	}

	fmt.Println("\nDone!")
}

func processFile(path string, dryRun bool) error {
	content, err := os.ReadFile(path)
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
			err = os.WriteFile(path, modifiedContent, 0644)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %v", path, err)
			}
		}
	}

	return nil
}
