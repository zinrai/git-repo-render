package main

import (
	"flag"
	"fmt"
	"os"
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

	excludeList := parseExcludeList(*excludePaths)

	entries, err := collectEntries(*reposPath, "", excludeList)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	structure := buildTree(entries)
	files := buildFileContents(entries)
	content := structure + "\n" + files

	if err := os.WriteFile(*outputFile, []byte(content), 0644); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Directory structure and file contents have been written to %s\n", *outputFile)
}
