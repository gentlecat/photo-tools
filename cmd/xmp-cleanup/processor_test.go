package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProcessDirectory_NoPairs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create only standalone files
	xmpContent := []byte(`<?xml version="1.0"?><x:xmpmeta xmlns:x="adobe:ns:meta/"></x:xmpmeta>`)
	if err := os.WriteFile(filepath.Join(tmpDir, "standalone.xmp"), xmpContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err := ProcessDirectory(tmpDir, false)
	if err != nil {
		t.Fatalf("ProcessDirectory failed: %v", err)
	}

	// Standalone file should still exist
	if _, err := os.Stat(filepath.Join(tmpDir, "standalone.xmp")); err != nil {
		t.Errorf("standalone.xmp should still exist")
	}
}

func TestProcessDirectory_WithPairs(t *testing.T) {
	tmpDir := t.TempDir()

	withExt := filepath.Join(tmpDir, "photo.jpg.xmp")
	withoutExt := filepath.Join(tmpDir, "photo.xmp")
	baseFile := filepath.Join(tmpDir, "photo.jpg")

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

	// Create files with different timestamps
	if err := os.WriteFile(withoutExt, olderContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	time.Sleep(10 * time.Millisecond)
	if err := os.WriteFile(withExt, newerContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if err := os.WriteFile(baseFile, []byte("fake"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err := ProcessDirectory(tmpDir, false)
	if err != nil {
		t.Fatalf("ProcessDirectory failed: %v", err)
	}

	// WithExt should be removed
	if _, err := os.Stat(withExt); !os.IsNotExist(err) {
		t.Errorf("photo.jpg.xmp should be removed after processing")
	}

	// WithoutExt should exist
	if _, err := os.Stat(withoutExt); err != nil {
		t.Errorf("photo.xmp should exist after processing: %v", err)
	}

	// Base file should still exist
	if _, err := os.Stat(baseFile); err != nil {
		t.Errorf("photo.jpg should still exist after processing: %v", err)
	}
}

func TestProcessDirectory_DryRun(t *testing.T) {
	tmpDir := t.TempDir()

	withExt := filepath.Join(tmpDir, "photo.jpg.xmp")
	withoutExt := filepath.Join(tmpDir, "photo.xmp")
	baseFile := filepath.Join(tmpDir, "photo.jpg")

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

	err := ProcessDirectory(tmpDir, true)
	if err != nil {
		t.Fatalf("ProcessDirectory failed: %v", err)
	}

	// Both files should still exist in dry run mode
	if _, err := os.Stat(withExt); err != nil {
		t.Errorf("photo.jpg.xmp should still exist after dry run")
	}
	if _, err := os.Stat(withoutExt); err != nil {
		t.Errorf("photo.xmp should still exist after dry run")
	}
}

func TestProcessDirectory_MultiplePairs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple pairs
	for i := 1; i <= 3; i++ {
		baseName := filepath.Join(tmpDir, "photo"+string(rune('0'+i)))
		withExt := baseName + ".jpg.xmp"
		withoutExt := baseName + ".xmp"
		baseFile := baseName + ".jpg"

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
	}

	err := ProcessDirectory(tmpDir, false)
	if err != nil {
		t.Fatalf("ProcessDirectory failed: %v", err)
	}

	// Check all withExt files are removed
	for i := 1; i <= 3; i++ {
		baseName := filepath.Join(tmpDir, "photo"+string(rune('0'+i)))
		withExt := baseName + ".jpg.xmp"
		withoutExt := baseName + ".xmp"

		if _, err := os.Stat(withExt); !os.IsNotExist(err) {
			t.Errorf("%s should be removed", withExt)
		}
		if _, err := os.Stat(withoutExt); err != nil {
			t.Errorf("%s should exist: %v", withoutExt, err)
		}
	}
}

func TestProcessDirectory_NonExistentDir(t *testing.T) {
	err := ProcessDirectory("/nonexistent/directory", false)
	if err == nil {
		t.Errorf("ProcessDirectory should fail for non-existent directory")
	}
}
