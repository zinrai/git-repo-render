package main

import (
	"path/filepath"
	"strings"
)

// parseExcludeList splits and trims a comma-separated exclude string.
func parseExcludeList(s string) []string {
	parts := strings.Split(s, ",")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts
}

// shouldExclude determines whether a path should be excluded based on exact
// basename matching against the exclude list.
func shouldExclude(path string, excludeList []string) bool {
	for _, exclude := range excludeList {
		if exclude == "" {
			continue
		}
		if filepath.Base(path) == exclude {
			return true
		}
	}
	return false
}
