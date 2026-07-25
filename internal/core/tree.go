package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func BuildTree(w *bufio.Writer, root string, cfg Config) error {
	fmt.Fprintln(w, ".")
	return generateTree(w, root, root, "", cfg)
}

func generateTree(w *bufio.Writer, root, dir, prefix string, cfg Config) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	entries = filterEntries(entries, root, dir, cfg)
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })

	for i, e := range entries {
		isLast := i == len(entries)-1
		branch := "├── "
		next := prefix + "│   "
		if isLast {
			branch = "└── "
			next = prefix + "    "
		}
		fmt.Fprintf(w, "%s%s%s\n", prefix, branch, e.Name())
		if e.IsDir() {
			info, _ := e.Info()
			if info != nil && info.Mode()&os.ModeSymlink != 0 {
				continue
			}
			if err := generateTree(w, root, filepath.Join(dir, e.Name()), next, cfg); err != nil {
				return err
			}
		}
	}
	return nil
}
