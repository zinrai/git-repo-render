package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: git-repo-render <repository_path> [output_file]")
		os.Exit(1)
	}

	dirPath := os.Args[1]
	outputFile := "output.txt"
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	}

	structure, files, err := exploreDirectory(dirPath, "", true, "")
	if err != nil {
		fmt.Printf("Error exploring directory: %v\n", err)
		os.Exit(1)
	}

	content := structure + "\n" + files

	err = os.WriteFile(outputFile, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Directory structure and file contents have been written to %s\n", outputFile)
}

func exploreDirectory(dirPath, indent string, isRoot bool, relPath string) (string, string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", "", err
	}

	var structure strings.Builder
	var files strings.Builder

	if !isRoot {
		structure.WriteString(fmt.Sprintf("%s%s\n", indent, filepath.Base(dirPath)))
		indent += "│   "
	}

	for i, entry := range entries {
		if entry.Name() == ".git" {
			continue // Skip .git directory
		}

		path := filepath.Join(dirPath, entry.Name())
		currentRelPath := filepath.Join(relPath, entry.Name())

		if entry.IsDir() {
			subStructure, subFiles, err := exploreDirectory(path, indent, false, currentRelPath)
			if err != nil {
				return "", "", err
			}
			structure.WriteString(subStructure)
			files.WriteString(subFiles)
		} else {
			if i == len(entries)-1 {
				structure.WriteString(fmt.Sprintf("%s└── %s\n", indent, entry.Name()))
			} else {
				structure.WriteString(fmt.Sprintf("%s├── %s\n", indent, entry.Name()))
			}

			// Add file content
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return "", "", err
			}
			files.WriteString(fmt.Sprintf("\n// %s\n%s\n", currentRelPath, string(fileContent)))
		}
	}

	return structure.String(), files.String(), nil
}
