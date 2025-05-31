// Copyright © 2025 Luka Ivanović
// This code is licensed under the terms of the MIT licence (see LICENCE for details)

package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

type FileInfo struct {
	Path     string
	Lines    int
	IsDir    bool
	Children []*FileInfo
}

var (
	silent = flag.Bool("silent", false, "print only the number of lines of code")
)

func main() {
	flag.Parse()
	args := flag.Args()
	rootPath := "."
	if len(args) != 0 {
		rootPath = args[0]
	}

	root := countLinesInDirectory(rootPath)
	if root != nil {
		if *silent {
			fmt.Println(root.Lines)
		} else {
			printResults(root, 0, "", true)
		}
	}
}

func isValid(filename string) bool {
	filename = path.Base(filename)
	if strings.HasPrefix(filename, ".") {
		return false
	}
	if slices.Contains([]string{".gz", ".json", ".zip", ".toml", ".jpg", ".jpeg", ".png", ".1", ".md", ".rst", ".sh", ".fish", ".txt", ".ttf", ".css", ".ico", ".zsh", ".ps1", ".gap", ".py", ".just", ".snap"}, filepath.Ext(filename)) {
		return false
	}
	if strings.Contains(filename, "testdata") {
		return false
	}
	if strings.Contains(filename, "changelog") {
		return false
	}
	if strings.Contains(filename, "azure") {
		return false
	}
	if strings.Contains(filename, "_test") {
		return false
	}
	if strings.Contains(filename, "tests") {
		return false
	}
	if strings.Contains(filename, "/doc/") {
		return false
	}
	if strings.Contains(filename, "Makefile") {
		return false
	}
	if strings.Contains(filename, "LICENSE") || strings.Contains(filename, "LICENCE") {
		return false
	}
	if strings.Contains(filename, "VERSION") {
		return false
	}
	if strings.Contains(filename, "docker") {
		return false
	}
	return true
}

func countLinesInDirectory(root string) *FileInfo {
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) == 0 {
		return nil
	}
	res := &FileInfo{
		Path:     root,
		Lines:    0,
		IsDir:    true,
		Children: make([]*FileInfo, 0),
	}

	for _, entry := range entries {
		filename := path.Join(root, entry.Name())
		if !isValid(filename) {
			continue
		}
		if entry.IsDir() {
			tmp := countLinesInDirectory(filename)
			if tmp != nil && tmp.Lines > 0 {
				res.Children = append(res.Children, tmp)
				res.Lines += tmp.Lines
			}
		} else {
			linesInFile := countLinesInFile(filename)
			if linesInFile > 0 {
				tmp := &FileInfo{
					Path:     filename,
					Lines:    linesInFile,
					IsDir:    false,
					Children: nil,
				}
				res.Children = append(res.Children, tmp)
				res.Lines += linesInFile
			}
		}
	}
	return res
}

func countLinesInFile(filename string) int {
	content, _ := os.ReadFile(filename)
	if len(content) == 0 {
		return 0
	}
	nl := 1
	for _, b := range content {
		if b == '\n' {
			nl += 1
		}
	}
	return nl
}

func printResults(fi *FileInfo, level int, prefix string, isLast bool) {
	symbol := ""
	if level == 0 {
		symbol = ""
	} else if isLast {
		symbol = "└── "
	} else {
		symbol = "├── "
	}

	displayPath := fi.Path
	if level > 0 {
		displayPath = path.Base(fi.Path)
	}
	if fi.IsDir {
		fmt.Printf("%s%s%s/: %d\n", prefix, symbol, displayPath, fi.Lines)
	} else {
		fmt.Printf("%s%s%s: %d\n", prefix, symbol, displayPath, fi.Lines)
	}

	newPrefix := ""
	if level == 0 {
		newPrefix = ""
	} else if isLast {
		newPrefix = prefix + "    "
	} else {
		newPrefix = prefix + "│   "
	}

	l := len(fi.Children) - 1
	for i, child := range fi.Children {
		if i == l {
			printResults(child, level+1, newPrefix, true)
		} else {
			printResults(child, level+1, newPrefix, false)
		}
	}
}
