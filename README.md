# git-repo-render

`git-repo-render` is a command-line tool that generates a comprehensive view of a Git repository's structure and file contents. It's designed to provide a quick and easy way to understand the layout and content of a project.

## Features

- Displays the directory structure of a Git repository
- Shows the content of each file
- Ignores specified directories and files (by default, the `.git` directory)
- Outputs the result to a file for easy viewing and sharing
- Supports customizable exclusion patterns for skipping unwanted files and directories

## Installation

To install `git-repo-render`, you need to have Go installed on your system.

```
$ go build
```

## Usage

To use `git-repo-render`, run the following command:

```
$ git-repo-render [options]
```

### Options

- `-repos <repository_path>`: The path to the Git repository you want to render (defaults to current directory)
- `-out <output_file>`: The name of the file to write the output to (defaults to "output.txt")
- `-exclude <exclude_paths>`: Comma-separated list of paths to exclude (defaults to ".git")
- `-h`: Display help information

### Examples

Render the current directory:

```
$ git-repo-render
```

Render a specific repository:

```
$ git-repo-render -repos /path/to/your/repo
```

Specify a custom output file:

```
$ git-repo-render -repos /path/to/your/repo -out my_repo_structure.txt
```

Exclude multiple directories:

```
$ git-repo-render -repos /path/to/your/repo -exclude .git,node_modules,vendor
```

## Output Format

The output file will contain two main sections:

1. Directory Structure: A tree-like representation of the repository's file and directory structure.
2. File Contents: The contents of each file, preceded by its relative path within the repository.

Example output:

```
.github
│   workflows
│   │   └── ghcr-publish.yml
├── .gitignore
├── Dockerfile
├── main.py
└── requirements.txt

// .github/workflows/ghcr-publish.yml
name: Build and Publish

on:
  push:
    tags:
      - 'v*'

...

// .gitignore
venv

// Dockerfile
FROM python:3.11

...

// main.py
import falcon
import json

...

// requirements.txt
falcon==3.1.3
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
