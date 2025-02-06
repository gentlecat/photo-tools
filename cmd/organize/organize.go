// This script organizes photos into subdirectories based on original timestamps from their metadata.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/evanoberholster/imagemeta"
)

func main() {
	dirPath := flag.String("dir", ".", "Path to the directory containing image files")
	flag.Parse()

	if *dirPath == "" {
		fmt.Println("Directory path must be specified")
		flag.Usage()
		os.Exit(1)
	}

	err := organizeImages(*dirPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nDone!")
}

func organizeImages(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		err := processFile(filepath.Join(dirPath, entry.Name()), entry, dirPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func processFile(path string, info os.DirEntry, dirPath string) error {
	if info.IsDir() {
		return nil // Skip directories
	}

	if !isSupportedFormat(path) {
		return nil // Skip non-image files
	}

	date, err := extractCaptureDate(path)
	if err != nil {
		log.Printf("Error extracting date for %s: %v", path, err)
		return nil
	}

	if date.IsZero() {
		log.Printf("Date is missing for %s", path)
		return nil
	}

	year := date.Year()
	month := int(date.Month())
	day := date.Day()

	newDir := filepath.Join(
		dirPath,
		fmt.Sprintf("%d", year),
		fmt.Sprintf("%d-%02d", year, month),
		fmt.Sprintf("%d-%02d-%02d", year, month, day),
	)

	err = os.MkdirAll(newDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	newPath := filepath.Join(newDir, filepath.Base(path))
	err = os.Rename(path, newPath)
	log.Printf("%s -> %s", path, newPath)
	if err != nil {
		return err
	}

	return nil
}

func extractCaptureDate(path string) (time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		return time.Time{}, err
	}
	defer f.Close()

	metadata, err := imagemeta.Decode(f)
	if err != nil {
		return time.Time{}, err
	}

	return metadata.DateTimeOriginal(), nil
}

func isSupportedFormat(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".heic"
}
