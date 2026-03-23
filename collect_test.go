package main

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestDir creates a temporary directory with the given structure.
// dirs is a list of relative directory paths.
// files is a map of relative file path to content.
func setupTestDir(t *testing.T, dirs []string, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(root, d), 0755); err != nil {
			t.Fatal(err)
		}
	}
	for path, content := range files {
		full := filepath.Join(root, path)
		dir := filepath.Dir(full)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestCollectEntries(t *testing.T) {
	root := setupTestDir(t,
		[]string{".git", ".github/workflows", "src"},
		map[string]string{
			".git/config":                     "[core]",
			".github/workflows/ci.yml":        "name: CI",
			"src/main.go":                     "package main",
			"README.md":                       "# Hello",
			filepath.Join("src", "image.png"): string([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00}),
		},
	)

	entries, err := collectEntries(root, "", []string{".git"})
	if err != nil {
		t.Fatal(err)
	}

	relPaths := make(map[string]bool)
	for _, e := range entries {
		relPaths[e.RelPath] = true
	}

	// .git should be excluded, .github should remain
	if relPaths[".git"] || relPaths[filepath.Join(".git", "config")] {
		t.Error(".git should be excluded")
	}
	if !relPaths[".github"] {
		t.Error(".github should NOT be excluded")
	}
	if !relPaths[filepath.Join(".github", "workflows")] {
		t.Error(".github/workflows should NOT be excluded")
	}
	if !relPaths["README.md"] {
		t.Error("README.md should be present")
	}

	// Check binary detection through collectEntries
	for _, e := range entries {
		if e.RelPath == filepath.Join("src", "image.png") {
			if !e.IsBin {
				t.Error("image.png should be detected as binary")
			}
		}
		if e.RelPath == filepath.Join("src", "main.go") {
			if e.IsBin {
				t.Error("main.go should not be detected as binary")
			}
		}
	}
}

func TestCollectEntriesEmptyDir(t *testing.T) {
	root := t.TempDir()

	entries, err := collectEntries(root, "", []string{".git"})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for empty dir, got %d", len(entries))
	}
}
