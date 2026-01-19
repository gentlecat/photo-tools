package main

import (
	"os"
	"path/filepath"
	"strings"
)

// XMPPair represents a pair of XMP files that need merging
type XMPPair struct {
	WithExt    string // e.g., "photo.jpg.xmp"
	WithoutExt string // e.g., "photo.xmp"
	BaseFile   string // e.g., "photo.jpg"
}

// SoloXMPFile represents a single XMP file with extension that should be renamed
type SoloXMPFile struct {
	Current string // e.g., "photo.jpg.xmp"
	Target  string // e.g., "photo.xmp"
}

// ScanForXMPPairs finds XMP file pairs that need merging, recursively scanning subdirectories
func ScanForXMPPairs(dirPath string) ([]XMPPair, error) {
	var allPairs []XMPPair

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		// Scan current directory
		pairs, err := scanDirectoryForPairs(path)
		if err != nil {
			return err
		}
		allPairs = append(allPairs, pairs...)
		return nil
	})

	return allPairs, err
}

func scanDirectoryForPairs(dirPath string) ([]XMPPair, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	xmpFiles := make(map[string]string) // base name -> full path
	var pairs []XMPPair

	// First pass: collect all .xmp files
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".xmp") {
			continue
		}

		fullPath := filepath.Join(dirPath, name)
		baseName := strings.TrimSuffix(name, filepath.Ext(name))
		xmpFiles[baseName] = fullPath
	}

	// Second pass: find pairs
	for baseName, xmpPath := range xmpFiles {
		// Check if this is "photo.ext.xmp" format
		if ext := filepath.Ext(baseName); ext != "" {
			// This is "photo.ext.xmp", check for "photo.xmp"
			baseWithoutExt := strings.TrimSuffix(baseName, ext)
			if altXmpPath, exists := xmpFiles[baseWithoutExt]; exists {
				// Found a pair
				baseFilePath := filepath.Join(dirPath, baseName)
				if _, err := os.Stat(baseFilePath); err == nil {
					pairs = append(pairs, XMPPair{
						WithExt:    xmpPath,
						WithoutExt: altXmpPath,
						BaseFile:   baseFilePath,
					})
				}
			}
		}
	}

	return pairs, nil
}

// ScanForSoloXMPFiles finds XMP files with extensions that have no corresponding file without extension, recursively scanning subdirectories
func ScanForSoloXMPFiles(dirPath string) ([]SoloXMPFile, error) {
	var allSoloFiles []SoloXMPFile

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		// Scan current directory
		soloFiles, err := scanDirectoryForSoloFiles(path)
		if err != nil {
			return err
		}
		allSoloFiles = append(allSoloFiles, soloFiles...)
		return nil
	})

	return allSoloFiles, err
}

func scanDirectoryForSoloFiles(dirPath string) ([]SoloXMPFile, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	xmpFiles := make(map[string]string) // base name -> full path
	var soloFiles []SoloXMPFile

	// First pass: collect all .xmp files
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".xmp") {
			continue
		}

		fullPath := filepath.Join(dirPath, name)
		baseName := strings.TrimSuffix(name, filepath.Ext(name))
		xmpFiles[baseName] = fullPath
	}

	// Second pass: find solo files with extension
	for baseName, xmpPath := range xmpFiles {
		// Check if this is "photo.ext.xmp" format
		if ext := filepath.Ext(baseName); ext != "" {
			// This is "photo.ext.xmp", check if "photo.xmp" exists
			baseWithoutExt := strings.TrimSuffix(baseName, ext)
			if _, exists := xmpFiles[baseWithoutExt]; !exists {
				// No corresponding file without extension, this is a solo file
				targetPath := filepath.Join(dirPath, baseWithoutExt+".xmp")
				soloFiles = append(soloFiles, SoloXMPFile{
					Current: xmpPath,
					Target:  targetPath,
				})
			}
		}
	}

	return soloFiles, nil
}
