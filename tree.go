package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// BuildTree writes the directory tree to the writer
func BuildTree(writer *bufio.Writer, root string, config Config) error {
	fmt.Fprintln(writer, ".")
	return generateTree(writer, root, root, "", config)
}

// generateTree recursively prints the directory structure
func generateTree(writer *bufio.Writer, root, path, prefix string, config Config) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Filter and sort
	entries = filterEntries(entries, root, path, config)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for i, entry := range entries {
		isLast := i == len(entries)-1

		branch := "├── "
		nextPrefix := prefix + "│   "
		if isLast {
			branch = "└── "
			nextPrefix = prefix + "    "
		}

		fmt.Fprintf(writer, "%s%s%s\n", prefix, branch, entry.Name())

		if entry.IsDir() {
			// Skip symlinks to avoid infinite loops
			info, err := entry.Info()
			if err != nil {
				continue
			}
			if info.Mode()&os.ModeSymlink != 0 {
				continue
			}

			if err := generateTree(writer, root, filepath.Join(path, entry.Name()), nextPrefix, config); err != nil {
				return err
			}
		}
	}
	return nil
}
