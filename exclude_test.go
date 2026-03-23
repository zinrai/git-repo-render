package main

import (
	"path/filepath"
	"testing"
)

func TestShouldExclude(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		excludeList []string
		want        bool
	}{
		{
			name:        "exact match on basename",
			path:        ".git",
			excludeList: []string{".git"},
			want:        true,
		},
		{
			name:        "no match on partial name",
			path:        ".github",
			excludeList: []string{".git"},
			want:        false,
		},
		{
			name:        "nested path matches basename",
			path:        filepath.Join("src", ".git"),
			excludeList: []string{".git"},
			want:        true,
		},
		{
			name:        "nested path does not match parent containing substring",
			path:        filepath.Join(".github", "workflows"),
			excludeList: []string{".git"},
			want:        false,
		},
		{
			name:        "empty exclude list entry is skipped",
			path:        "README.md",
			excludeList: []string{""},
			want:        false,
		},
		{
			name:        "multiple excludes first matches",
			path:        "node_modules",
			excludeList: []string{".git", "node_modules", "vendor"},
			want:        true,
		},
		{
			name:        "multiple excludes last matches",
			path:        filepath.Join("src", "vendor"),
			excludeList: []string{".git", "node_modules", "vendor"},
			want:        true,
		},
		{
			name:        "multiple excludes none match",
			path:        filepath.Join("src", "main.go"),
			excludeList: []string{".git", "node_modules", "vendor"},
			want:        false,
		},
		{
			name:        "similar name not excluded",
			path:        "legit",
			excludeList: []string{".git"},
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldExclude(tt.path, tt.excludeList)
			if got != tt.want {
				t.Errorf("shouldExclude(%q, %v) = %v, want %v", tt.path, tt.excludeList, got, tt.want)
			}
		})
	}
}

func TestParseExcludeList(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "single entry",
			input: ".git",
			want:  []string{".git"},
		},
		{
			name:  "multiple entries",
			input: ".git,node_modules,vendor",
			want:  []string{".git", "node_modules", "vendor"},
		},
		{
			name:  "entries with spaces",
			input: ".git , node_modules , vendor",
			want:  []string{".git", "node_modules", "vendor"},
		},
		{
			name:  "empty string produces single empty entry",
			input: "",
			want:  []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseExcludeList(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("parseExcludeList(%q) returned %d items, want %d", tt.input, len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseExcludeList(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}
