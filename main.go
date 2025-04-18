package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	reposPath := flag.String("repos", ".", "Path to the Git repository (defaults to current directory)")
	outputFile := flag.String("out", "output.txt", "Output file name")
	excludePaths := flag.String("exclude", ".git", "Comma-separated list of paths to exclude")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -repos <repository_path> [-out output_file] [-exclude exclude_paths]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -repos /path/to/repo\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -repos /path/to/repo -out my_output.txt\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -repos /path/to/repo -exclude .git,node_modules,vendor\n", os.Args[0])
	}

	flag.Parse()

	if _, err := os.Stat(*outputFile); err == nil {
		fmt.Printf("Error: %s already exists in the current directory\n", *outputFile)
		os.Exit(1)
	}

	excludeList := strings.Split(*excludePaths, ",")
	for i, exclude := range excludeList {
		excludeList[i] = strings.TrimSpace(exclude)
	}

	structure, files := exploreDirectory(*reposPath, "", true, "", excludeList)
	content := structure + "\n" + files

	err := os.WriteFile(*outputFile, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Directory structure and file contents have been written to %s\n", *outputFile)
}

// Determine whether a path should be excluded based on the exclude list
func shouldExclude(path string, excludeList []string) bool {
	for _, exclude := range excludeList {
		if exclude == "" {
			continue
		}

		// Check if filename or directory name matches the exclude list
		if filepath.Base(path) == exclude {
			return true
		}

		// Check if any part of the path matches the exclude list
		if strings.Contains(path, exclude) {
			return true
		}
	}
	return false
}

// Process a file and return its formatted content
func processFile(filePath, relPath string) string {
	var fileOutput strings.Builder

	fileOutput.WriteString(fmt.Sprintf("\n// %s\n", relPath))

	// Check if the file is binary
	if isBinary(filePath) {
		fileOutput.WriteString("(Binary file)\n")
		return fileOutput.String()
	}

	// Add file content for non-binary files
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		fileOutput.WriteString("(Error reading file)\n")
		return fileOutput.String()
	}

	fileOutput.WriteString(string(fileContent))
	fileOutput.WriteString("\n")

	return fileOutput.String()
}

// Add a file entry to the directory structure
func addFileToStructure(indent string, fileName string, isLast bool) string {
	if isLast {
		return fmt.Sprintf("%s└── %s\n", indent, fileName)
	}
	return fmt.Sprintf("%s├── %s\n", indent, fileName)
}

func exploreDirectory(dirPath, indent string, isRoot bool, relPath string, excludeList []string) (string, string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", dirPath, err)
		return "", ""
	}

	var structure strings.Builder
	var files strings.Builder

	if !isRoot {
		structure.WriteString(fmt.Sprintf("%s%s\n", indent, filepath.Base(dirPath)))
		indent += "│   "
	}

	// Filter entries to only those not excluded
	var filteredEntries []os.DirEntry
	for _, entry := range entries {
		currentRelPath := filepath.Join(relPath, entry.Name())
		if !shouldExclude(currentRelPath, excludeList) {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	// Process each entry
	for i, entry := range filteredEntries {
		path := filepath.Join(dirPath, entry.Name())
		currentRelPath := filepath.Join(relPath, entry.Name())
		isLast := i == len(filteredEntries)-1

		if entry.IsDir() {
			// Handle directory
			subStructure, subFiles := exploreDirectory(path, indent, false, currentRelPath, excludeList)
			structure.WriteString(subStructure)
			files.WriteString(subFiles)
		} else {
			// Handle file
			structure.WriteString(addFileToStructure(indent, entry.Name(), isLast))
			files.WriteString(processFile(path, currentRelPath))
		}
	}

	return structure.String(), files.String()
}

func isBinary(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Read the first 512 bytes
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return false
	}

	// Use the DetectContentType function to determine if it's binary
	contentType := http.DetectContentType(buffer[:n])
	return !strings.HasPrefix(contentType, "text/") && contentType != "application/json"
}
