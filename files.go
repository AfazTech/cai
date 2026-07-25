package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// binarySniffLen is the number of leading bytes inspected to detect binary content.
const binarySniffLen = 8000

// BuildFiles writes the content of all non-excluded files to the writer
func BuildFiles(writer *bufio.Writer, root string, config Config) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}

		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		// Skip symlinks
		if info.Mode()&os.ModeSymlink != 0 {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Apply include and ignore filters
		if !IsIncluded(relativePath, config) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if ShouldIgnore(relativePath, path, root, config) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		fmt.Fprintf(writer, FileDelimiter, filepath.ToSlash(relativePath))
		fmt.Fprintln(writer)
		fmt.Fprintln(writer)

		if isBinaryContent(content) {
			fmt.Fprintln(writer, "[binary file omitted]")
		} else {
			writer.Write(content)
		}

		fmt.Fprintln(writer)
		fmt.Fprintln(writer)

		return nil
	})
}

// isBinaryContent does a quick heuristic check for binary data by looking
// for a NUL byte in the first bytes of the file, similar to how git detects
// binary files.
func isBinaryContent(content []byte) bool {
	sniff := content
	if len(sniff) > binarySniffLen {
		sniff = sniff[:binarySniffLen]
	}
	return bytes.IndexByte(sniff, 0) != -1
}
