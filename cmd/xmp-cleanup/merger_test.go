package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/trimmer-io/go-xmp/xmp"
)

func TestMergeXMPFiles_DryRun(t *testing.T) {
	tmpDir := t.TempDir()

	withExt := filepath.Join(tmpDir, "test.jpg.xmp")
	withoutExt := filepath.Join(tmpDir, "test.xmp")
	baseFile := filepath.Join(tmpDir, "test.jpg")

	xmpContent := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:xmp="http://ns.adobe.com/xap/1.0/" rdf:about="">
      <xmp:Rating>5</xmp:Rating>
    </rdf:Description>
  </rdf:RDF>
</x:xmpmeta>`)

	if err := os.WriteFile(withExt, xmpContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if err := os.WriteFile(withoutExt, xmpContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if err := os.WriteFile(baseFile, []byte("fake"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	pair := XMPPair{
		WithExt:    withExt,
		WithoutExt: withoutExt,
		BaseFile:   baseFile,
	}

	// Test dry run - files should not be modified
	err := MergeXMPFiles(pair, true)
	if err != nil {
		t.Fatalf("MergeXMPFiles failed: %v", err)
	}

	// Both files should still exist
	if _, err := os.Stat(withExt); err != nil {
		t.Errorf("WithExt file should still exist after dry run")
	}
	if _, err := os.Stat(withoutExt); err != nil {
		t.Errorf("WithoutExt file should still exist after dry run")
	}
}

func TestMergeXMPFiles_NewerWithExt(t *testing.T) {
	tmpDir := t.TempDir()

	withExt := filepath.Join(tmpDir, "test.jpg.xmp")
	withoutExt := filepath.Join(tmpDir, "test.xmp")
	baseFile := filepath.Join(tmpDir, "test.jpg")

	newerContent := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:xmp="http://ns.adobe.com/xap/1.0/" rdf:about="">
      <xmp:Rating>5</xmp:Rating>
    </rdf:Description>
  </rdf:RDF>
</x:xmpmeta>`)

	olderContent := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:xmp="http://ns.adobe.com/xap/1.0/" rdf:about="">
      <xmp:Rating>3</xmp:Rating>
    </rdf:Description>
  </rdf:RDF>
</x:xmpmeta>`)

	// Create older file first
	if err := os.WriteFile(withoutExt, olderContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	time.Sleep(1100 * time.Millisecond) // Ensure different timestamp (filesystem resolution can be 1s)

	// Create newer file
	if err := os.WriteFile(withExt, newerContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if err := os.WriteFile(baseFile, []byte("fake"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	pair := XMPPair{
		WithExt:    withExt,
		WithoutExt: withoutExt,
		BaseFile:   baseFile,
	}

	// Perform rename operation
	err := MergeXMPFiles(pair, false)
	if err != nil {
		t.Fatalf("MergeXMPFiles failed: %v", err)
	}

	// Original withExt should not exist (renamed to withoutExt)
	if _, err := os.Stat(withExt); !os.IsNotExist(err) {
		t.Errorf("WithExt file should not exist after processing")
	}

	// WithoutExt should exist with newer content
	if _, err := os.Stat(withoutExt); err != nil {
		t.Fatalf("WithoutExt file should exist after processing: %v", err)
	}

	// Old file should be renamed with _old suffix
	oldFile := filepath.Join(tmpDir, "test.xmp_old")
	if _, err := os.Stat(oldFile); err != nil {
		t.Fatalf("Old file should exist with _old suffix: %v", err)
	}

	// Verify content is from newer file
	keptData, err := os.ReadFile(withoutExt)
	if err != nil {
		t.Fatalf("failed to read kept file: %v", err)
	}

	doc, err := xmp.Read(bytes.NewReader(keptData))
	if err != nil {
		t.Fatalf("failed to parse kept XMP: %v", err)
	}
	defer doc.Close()

	rating, err := doc.GetPath(xmp.NewPath("xmp", "Rating"))
	if err != nil {
		t.Fatalf("failed to get rating: %v", err)
	}
	if rating != "5" {
		t.Errorf("expected rating 5 from newer file, got %s", rating)
	}

	// Verify old file has old content
	oldData, err := os.ReadFile(oldFile)
	if err != nil {
		t.Fatalf("failed to read old file: %v", err)
	}

	oldDoc, err := xmp.Read(bytes.NewReader(oldData))
	if err != nil {
		t.Fatalf("failed to parse old XMP: %v", err)
	}
	defer oldDoc.Close()

	oldRating, err := oldDoc.GetPath(xmp.NewPath("xmp", "Rating"))
	if err != nil {
		t.Fatalf("failed to get old rating: %v", err)
	}
	if oldRating != "3" {
		t.Errorf("expected rating 3 in old file, got %s", oldRating)
	}
}

func TestMergeXMPFiles_NewerWithoutExt(t *testing.T) {
	tmpDir := t.TempDir()

	withExt := filepath.Join(tmpDir, "test.jpg.xmp")
	withoutExt := filepath.Join(tmpDir, "test.xmp")
	baseFile := filepath.Join(tmpDir, "test.jpg")

	olderContent := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:xmp="http://ns.adobe.com/xap/1.0/" rdf:about="">
      <xmp:Label>Old Label</xmp:Label>
    </rdf:Description>
  </rdf:RDF>
</x:xmpmeta>`)

	newerContent := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description xmlns:xmp="http://ns.adobe.com/xap/1.0/" rdf:about="">
      <xmp:Label>New Label</xmp:Label>
    </rdf:Description>
  </rdf:RDF>
</x:xmpmeta>`)

	// Create older file with extension first
	if err := os.WriteFile(withExt, olderContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	time.Sleep(1100 * time.Millisecond) // Ensure different timestamp

	// Create newer file without extension
	if err := os.WriteFile(withoutExt, newerContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if err := os.WriteFile(baseFile, []byte("fake"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	pair := XMPPair{
		WithExt:    withExt,
		WithoutExt: withoutExt,
		BaseFile:   baseFile,
	}

	err := MergeXMPFiles(pair, false)
	if err != nil {
		t.Fatalf("MergeXMPFiles failed: %v", err)
	}

	// WithoutExt should still exist (it was newer)
	if _, err := os.Stat(withoutExt); err != nil {
		t.Fatalf("WithoutExt file should exist: %v", err)
	}

	// WithExt should be renamed to _old
	oldFile := filepath.Join(tmpDir, "test.jpg.xmp_old")
	if _, err := os.Stat(oldFile); err != nil {
		t.Fatalf("Old file should exist with _old suffix: %v", err)
	}

	// Verify content is from newer file
	keptData, err := os.ReadFile(withoutExt)
	if err != nil {
		t.Fatalf("failed to read kept file: %v", err)
	}

	doc, err := xmp.Read(bytes.NewReader(keptData))
	if err != nil {
		t.Fatalf("failed to parse kept XMP: %v", err)
	}
	defer doc.Close()

	label, err := doc.GetPath(xmp.NewPath("xmp", "Label"))
	if err != nil {
		t.Fatalf("failed to get label: %v", err)
	}
	if label != "New Label" {
		t.Errorf("expected 'New Label' from newer file, got %s", label)
	}
}
