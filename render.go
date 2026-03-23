package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// treeNode is used internally by buildTree to reconstruct the hierarchy.
type treeNode struct {
	name     string
	isDir    bool
	children []*treeNode
}

// buildTree takes a flat slice of FileEntry and produces the tree-formatted
// directory structure string.
func buildTree(entries []FileEntry) string {
	root := &treeNode{name: "", isDir: true}

	for _, e := range entries {
		parts := strings.Split(e.RelPath, string(filepath.Separator))
		cur := root
		for i, part := range parts {
			isLast := i == len(parts)-1
			found := false
			for _, child := range cur.children {
				if child.name == part {
					cur = child
					found = true
					break
				}
			}
			if !found {
				node := &treeNode{
					name:  part,
					isDir: e.IsDir || !isLast,
				}
				cur.children = append(cur.children, node)
				cur = node
			}
		}
	}

	var sb strings.Builder
	renderTree(&sb, root.children, "")
	return sb.String()
}

func renderTree(sb *strings.Builder, nodes []*treeNode, indent string) {
	for i, node := range nodes {
		isLast := i == len(nodes)-1
		if node.isDir {
			sb.WriteString(fmt.Sprintf("%s%s\n", indent, node.name))
			childIndent := indent + "│   "
			renderTree(sb, node.children, childIndent)
		} else {
			if isLast {
				sb.WriteString(fmt.Sprintf("%s└── %s\n", indent, node.name))
			} else {
				sb.WriteString(fmt.Sprintf("%s├── %s\n", indent, node.name))
			}
		}
	}
}

// buildFileContents takes a flat slice of FileEntry and produces the
// concatenated file contents section.
func buildFileContents(entries []FileEntry) string {
	var sb strings.Builder
	for _, e := range entries {
		if e.IsDir {
			continue
		}
		sb.WriteString(fmt.Sprintf("\n// %s\n", e.RelPath))
		if e.ReadErr {
			sb.WriteString("(Error reading file)\n")
		} else if e.IsBin {
			sb.WriteString("(Binary file)\n")
		} else {
			sb.WriteString(string(e.Content))
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
