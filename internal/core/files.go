package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const sniffLen = 8000

type LimitedWriter struct {
	w        io.Writer
	written  int64
	maxBytes int64
	exceeded bool
}

func NewLimitedWriter(w io.Writer, max int64) *LimitedWriter {
	return &LimitedWriter{w: w, maxBytes: max}
}

func (lw *LimitedWriter) Write(p []byte) (n int, err error) {
	if lw.exceeded {
		return 0, errors.New("max file size exceeded")
	}
	if lw.written+int64(len(p)) > lw.maxBytes {
		lw.exceeded = true
		return 0, errors.New("max file size exceeded")
	}
	n, err = lw.w.Write(p)
	lw.written += int64(n)
	return n, err
}

func BuildFiles(w io.Writer, root string, cfg Config) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		if info.Mode()&os.ModeSymlink != 0 {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !IsIncluded(rel, cfg) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if ShouldIgnore(rel, path, root, cfg) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}

		fmt.Fprintf(bw, FileDelimiter, filepath.ToSlash(rel))
		fmt.Fprintln(bw)
		fmt.Fprintln(bw)

		// Stream file content with limited copy
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		// check binary with first sniffLen bytes
		sniff := make([]byte, sniffLen)
		n, _ := f.Read(sniff)
		f.Seek(0, 0) // reset

		if n > 0 && bytesContainsNull(sniff[:n]) {
			fmt.Fprintln(bw, "[binary file omitted]")
		} else {
			// Copy with limited writer
			lw := NewLimitedWriter(bw, 1<<62) // no limit here (outer limit handles)
			if _, err := io.Copy(lw, f); err != nil {
				return err
			}
		}

		fmt.Fprintln(bw)
		fmt.Fprintln(bw)
		return nil
	})
}

func bytesContainsNull(b []byte) bool {
	for _, c := range b {
		if c == 0 {
			return true
		}
	}
	return false
}
