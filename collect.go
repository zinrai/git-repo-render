package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileEntry represents a file or directory discovered during traversal.
type FileEntry struct {
	RelPath string
	IsDir   bool
	Content []byte
	IsBin   bool
	ReadErr bool
}

// collectEntries walks the directory tree rooted at dirPath, returning a flat
// slice of FileEntry values in the order they would appear in a tree listing.
// Directories appear as entries with IsDir=true before their children.
func collectEntries(dirPath, relPath string, excludeList []string) ([]FileEntry, error) {
	osEntries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", dirPath, err)
	}

	var result []FileEntry

	// Filter excluded entries
	var filtered []os.DirEntry
	for _, e := range osEntries {
		currentRelPath := filepath.Join(relPath, e.Name())
		if !shouldExclude(currentRelPath, excludeList) {
			filtered = append(filtered, e)
		}
	}

	for _, e := range filtered {
		absPath := filepath.Join(dirPath, e.Name())
		currentRelPath := filepath.Join(relPath, e.Name())

		if e.IsDir() {
			result = append(result, FileEntry{
				RelPath: currentRelPath,
				IsDir:   true,
			})
			children, err := collectEntries(absPath, currentRelPath, excludeList)
			if err != nil {
				return nil, err
			}
			result = append(result, children...)
		} else {
			entry := FileEntry{RelPath: currentRelPath}
			data, err := os.ReadFile(absPath)
			if err != nil {
				entry.ReadErr = true
			} else {
				entry.IsBin = isBinaryContent(data)
				entry.Content = data
			}
			result = append(result, entry)
		}
	}

	return result, nil
}
