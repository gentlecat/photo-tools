// This script cleans up XMP files by removing zero ratings. Some photo management applications write xmp:Rating="0"
// to XMP sidecar files when no rating has been assigned, which can interfere with other tools. This script scans a
// directory recursively, finds all .xmp files containing xmp:Rating="0", and removes that attribute from them.
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
	"sync/atomic"
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

	paths := make(chan string, *workers*10)
	errs := make(chan error, 1)

	var total, cleaned atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range paths {
				total.Add(1)
				changed, err := processFile(path, *dryRun)
				if err != nil {
					select {
					case errs <- err:
					default:
						log.Printf("error processing file: %v", err)
					}
				} else if changed {
					cleaned.Add(1)
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

	fmt.Printf("\nDone! Cleaned %d of %d file(s).\n", cleaned.Load(), total.Load())
}

func processFile(path string, dryRun bool) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to stat file %s: %w", path, err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("failed to read file %s: %w", path, err)
	}

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
			err = os.WriteFile(path, modifiedContent, info.Mode())
			if err != nil {
				return false, fmt.Errorf("failed to write file %s: %w", path, err)
			}
		}
	}

	return changed, nil
}
