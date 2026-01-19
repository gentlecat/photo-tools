package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanForXMPPairs(t *testing.T) {
	tests := []struct {
		name          string
		expectedPairs int
	}{
		{
			name:          "finds pairs in testdata",
			expectedPairs: 2, // photo1 and photo2 have pairs
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pairs, err := ScanForXMPPairs("testdata")
			if err != nil {
				t.Fatalf("ScanForXMPPairs failed: %v", err)
			}

			if len(pairs) != tt.expectedPairs {
				t.Errorf("expected %d pairs, got %d", tt.expectedPairs, len(pairs))
			}

			// Verify each pair has the correct structure
			for _, pair := range pairs {
				if pair.WithExt == "" || pair.WithoutExt == "" || pair.BaseFile == "" {
					t.Errorf("pair has empty fields: %+v", pair)
				}

				// Verify files exist
				if _, err := os.Stat(pair.WithExt); err != nil {
					t.Errorf("WithExt file doesn't exist: %s", pair.WithExt)
				}
				if _, err := os.Stat(pair.WithoutExt); err != nil {
					t.Errorf("WithoutExt file doesn't exist: %s", pair.WithoutExt)
				}
				if _, err := os.Stat(pair.BaseFile); err != nil {
					t.Errorf("BaseFile doesn't exist: %s", pair.BaseFile)
				}
			}
		})
	}
}

func TestScanForXMPPairs_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	pairs, err := ScanForXMPPairs(tmpDir)
	if err != nil {
		t.Fatalf("ScanForXMPPairs failed: %v", err)
	}
	if len(pairs) != 0 {
		t.Errorf("expected 0 pairs in empty directory, got %d", len(pairs))
	}
}

func TestScanForXMPPairs_OnlyStandalone(t *testing.T) {
	tmpDir := t.TempDir()

	// Create only standalone XMP file
	xmpContent := []byte(`<?xml version="1.0"?><x:xmpmeta xmlns:x="adobe:ns:meta/"></x:xmpmeta>`)
	if err := os.WriteFile(filepath.Join(tmpDir, "standalone.xmp"), xmpContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	pairs, err := ScanForXMPPairs(tmpDir)
	if err != nil {
		t.Fatalf("ScanForXMPPairs failed: %v", err)
	}
	if len(pairs) != 0 {
		t.Errorf("expected 0 pairs with only standalone XMP, got %d", len(pairs))
	}
}

func TestScanForXMPPairs_MissingBaseFile(t *testing.T) {
	tmpDir := t.TempDir()

	xmpContent := []byte(`<?xml version="1.0"?><x:xmpmeta xmlns:x="adobe:ns:meta/"></x:xmpmeta>`)
	// Create both XMP files but no base file
	if err := os.WriteFile(filepath.Join(tmpDir, "photo.jpg.xmp"), xmpContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "photo.xmp"), xmpContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	pairs, err := ScanForXMPPairs(tmpDir)
	if err != nil {
		t.Fatalf("ScanForXMPPairs failed: %v", err)
	}
	if len(pairs) != 0 {
		t.Errorf("expected 0 pairs when base file is missing, got %d", len(pairs))
	}
}
