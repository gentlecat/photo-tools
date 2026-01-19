package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRenameSoloXMPFile_DryRun(t *testing.T) {
	tmpDir := t.TempDir()

	current := filepath.Join(tmpDir, "photo.jpg.xmp")
	target := filepath.Join(tmpDir, "photo.xmp")

	xmpContent := []byte(`<?xml version="1.0"?><x:xmpmeta xmlns:x="adobe:ns:meta/"></x:xmpmeta>`)
	if err := os.WriteFile(current, xmpContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	solo := SoloXMPFile{
		Current: current,
		Target:  target,
	}

	err := RenameSoloXMPFile(solo, true)
	if err != nil {
		t.Fatalf("RenameSoloXMPFile failed: %v", err)
	}

	// Original file should still exist in dry run
	if _, err := os.Stat(current); err != nil {
		t.Errorf("original file should still exist after dry run")
	}

	// Target file should not exist
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Errorf("target file should not exist after dry run")
	}
}

func TestRenameSoloXMPFile_ActualRename(t *testing.T) {
	tmpDir := t.TempDir()

	current := filepath.Join(tmpDir, "photo.jpg.xmp")
	target := filepath.Join(tmpDir, "photo.xmp")

	xmpContent := []byte(`<?xml version="1.0"?><x:xmpmeta xmlns:x="adobe:ns:meta/"></x:xmpmeta>`)
	if err := os.WriteFile(current, xmpContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	solo := SoloXMPFile{
		Current: current,
		Target:  target,
	}

	err := RenameSoloXMPFile(solo, false)
	if err != nil {
		t.Fatalf("RenameSoloXMPFile failed: %v", err)
	}

	// Original file should not exist
	if _, err := os.Stat(current); !os.IsNotExist(err) {
		t.Errorf("original file should not exist after rename")
	}

	// Target file should exist
	if _, err := os.Stat(target); err != nil {
		t.Errorf("target file should exist after rename: %v", err)
	}

	// Verify content is preserved
	newContent, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("failed to read renamed file: %v", err)
	}
	if string(newContent) != string(xmpContent) {
		t.Errorf("file content was not preserved after rename")
	}
}
