package torrent

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldIgnoreFile(t *testing.T) {
	tests := []struct {
		path            string
		excludePatterns []string
		includePatterns []string
		expected        bool
	}{
		// Built-in patterns are case-insensitive and should always be ignored.
		{path: "example.torrent", expected: true},
		{path: ".DS_store", expected: true},
		{path: ".ds_store", expected: true},
		{path: "Thumbs.db", expected: true},
		{path: "thumbs.db", expected: true},
		{path: "desktop.ini", expected: true},
		{path: "zone.identifier", expected: true},

		// Built-in patterns take precedence over includePatterns.
		{"example.torrent", []string{}, []string{"*.torrent"}, true},
		{"Thumbs.db", []string{}, []string{"Thumbs.db"}, true},

		// If no includePatterns or excludePatterns are specified, paths are included by default.
		{"media.mkv", []string{}, []string{}, false},
		{"media2.mkv", []string{}, []string{}, false},

		// If includePatterns are given, a path is ignored if not included.
		{"media.mkv", []string{}, []string{"media.mkv"}, false},
		{"media.mkv", []string{}, []string{"media2.mkv"}, true},

		// includePatterns should be case-insensitive.
		{"Media.mkv", []string{}, []string{"media.mkv"}, false},
		{"media.mkv", []string{}, []string{"Media.mkv"}, false},

		// If excludePatterns are given, a path is ignored only if excluded.
		{"media.mkv", []string{"media.mkv"}, []string{}, true},
		{"media.mkv", []string{"media2.mkv"}, []string{}, false},

		// excludePatterns should be case-insensitive.
		{"Media.mkv", []string{"media.mkv"}, []string{}, true},
		{"media.mkv", []string{"Media.mkv"}, []string{}, true},

		// If both patterns are given, only includePatterns is used.
		{"media.mkv", []string{"media.mkv"}, []string{"media.mkv"}, false},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			ignored, err := shouldIgnoreFile(tc.path, tc.excludePatterns, tc.includePatterns)
			assert.Nil(t, err)
			if ignored != tc.expected {
				t.Errorf("shouldIgnoreFile(%v, excludePatterns = %v, includePatterns = %v) = %v, want %v", tc.path, tc.excludePatterns, tc.includePatterns, ignored, tc.expected)
			}
		})
	}
}

func TestShouldIgnoreFileWithPathPatterns(t *testing.T) {
	tests := []struct {
		path            string
		excludePatterns []string
		includePatterns []string
		expected        bool
	}{
		{"dir/file", []string{}, []string{filepath.Join("dir", "*")}, false},
		{"dir2/file", []string{}, []string{filepath.Join("dir", "*")}, true},

		{"dir/file", []string{filepath.Join("dir", "*")}, []string{}, true},
		{"dir2/file", []string{filepath.Join("dir", "*")}, []string{}, false},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			ignored, err := shouldIgnoreFile(tc.path, tc.excludePatterns, tc.includePatterns)
			assert.Nil(t, err)
			if ignored != tc.expected {
				t.Errorf("shouldIgnoreFile(%v, excludePatterns = %v, includePatterns = %v) = %v, want %v", tc.path, tc.excludePatterns, tc.includePatterns, ignored, tc.expected)
			}
		})
	}
}
