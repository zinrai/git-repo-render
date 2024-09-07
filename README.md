# git-repo-render

`git-repo-render` is a command-line tool that generates a comprehensive view of a Git repository's structure and file contents. It's designed to provide a quick and easy way to understand the layout and content of a project.

## Features

- Displays the directory structure of a Git repository
- Shows the content of each file
- Ignores the `.git` directory to focus on the actual project files
- Outputs the result to a file for easy viewing and sharing

## Installation

To install `git-repo-render`, you need to have Go installed on your system. Then, follow these steps:

```
$ go build
```

## Usage

To use `git-repo-render`, run the following command:

```
$ git-repo-render <repository_path> [output_file]
```

- `<repository_path>`: The path to the Git repository you want to render
- `[output_file]`: (Optional) The name of the file to write the output to. If not specified, it defaults to "output.txt"

Example:

```
$ git-repo-render /path/to/your/repo my_repo_structure.txt
```

This will create a file named `my_repo_structure.txt` containing the structure and contents of the repository.

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
