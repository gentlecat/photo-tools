package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestProcessFile(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := ioutil.TempDir("", "xmp-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		input    string
		expected string
		changed  bool
	}{
		{
			name:     "contains double quotes",
			input:    `<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF><rdf:Description xmp:Rating="0"/></rdf:RDF></x:xmpmeta>`,
			expected: `<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF><rdf:Description /></rdf:RDF></x:xmpmeta>`,
			changed:  true,
		},
		{
			name:     "contains single quotes",
			input:    `<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF><rdf:Description xmp:Rating='0'/></rdf:RDF></x:xmpmeta>`,
			expected: `<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF><rdf:Description /></rdf:RDF></x:xmpmeta>`,
			changed:  true,
		},
		{
			name:     "no rating 0",
			input:    `<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF><rdf:Description xmp:Rating="5"/></rdf:RDF></x:xmpmeta>`,
			expected: `<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF><rdf:Description xmp:Rating="5"/></rdf:RDF></x:xmpmeta>`,
			changed:  false,
		},
		{
			name:     "mixed ratings",
			input:    `xmp:Rating="0" xmp:Rating="1" xmp:Rating='0'`,
			expected: ` xmp:Rating="1" `,
			changed:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.name+".xmp")
			err := ioutil.WriteFile(path, []byte(tt.input), 0644)
			if err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			err = processFile(path, false)
			if err != nil {
				t.Fatalf("processFile failed: %v", err)
			}

			content, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read file after process: %v", err)
			}

			if string(content) != tt.expected {
				t.Errorf("expected content %q, got %q", tt.expected, string(content))
			}
		})
	}
}
