package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanForSoloXMPFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a solo file with extension
	xmpContent := []byte(`<?xml version="1.0"?><x:xmpmeta xmlns:x="adobe:ns:meta/"></x:xmpmeta>`)
	if err := os.WriteFile(filepath.Join(tmpDir, "photo1.jpg.xmp"), xmpContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Create a pair (should not be in solo results)
	if err := os.WriteFile(filepath.Join(tmpDir, "photo2.cr2.xmp"), xmpContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "photo2.xmp"), xmpContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Create standalone file without extension (should not be in solo results)
	if err := os.WriteFile(filepath.Join(tmpDir, "photo3.xmp"), xmpContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	soloFiles, err := ScanForSoloXMPFiles(tmpDir)
	if err != nil {
		t.Fatalf("ScanForSoloXMPFiles failed: %v", err)
	}

	if len(soloFiles) != 1 {
		t.Errorf("expected 1 solo file, got %d", len(soloFiles))
	}

	if len(soloFiles) > 0 {
		expected := filepath.Join(tmpDir, "photo1.jpg.xmp")
		if soloFiles[0].Current != expected {
			t.Errorf("expected current path %s, got %s", expected, soloFiles[0].Current)
		}

		expectedTarget := filepath.Join(tmpDir, "photo1.xmp")
		if soloFiles[0].Target != expectedTarget {
			t.Errorf("expected target path %s, got %s", expectedTarget, soloFiles[0].Target)
		}
	}
}

func TestScanForSoloXMPFiles_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	soloFiles, err := ScanForSoloXMPFiles(tmpDir)
	if err != nil {
		t.Fatalf("ScanForSoloXMPFiles failed: %v", err)
	}

	if len(soloFiles) != 0 {
		t.Errorf("expected 0 solo files in empty directory, got %d", len(soloFiles))
	}
}

func TestScanForSoloXMPFiles_MultipleSoloFiles(t *testing.T) {
	tmpDir := t.TempDir()

	xmpContent := []byte(`<?xml version="1.0"?><x:xmpmeta xmlns:x="adobe:ns:meta/"></x:xmpmeta>`)

	// Create multiple solo files
	for i := 1; i <= 3; i++ {
		filename := filepath.Join(tmpDir, "photo"+string(rune('0'+i))+".jpg.xmp")
		if err := os.WriteFile(filename, xmpContent, 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	soloFiles, err := ScanForSoloXMPFiles(tmpDir)
	if err != nil {
		t.Fatalf("ScanForSoloXMPFiles failed: %v", err)
	}

	if len(soloFiles) != 3 {
		t.Errorf("expected 3 solo files, got %d", len(soloFiles))
	}
}
