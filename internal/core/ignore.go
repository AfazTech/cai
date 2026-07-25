package core

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	ignore "github.com/sabhiram/go-gitignore"
)

var (
	regexCache   = make(map[string]*regexp.Regexp)
	regexCacheMu sync.RWMutex

	gitignoreCache   = make(map[string]ignore.IgnoreParser)
	gitignoreCacheMu sync.RWMutex
)

// isRegexPattern checks if pattern starts with "regex:"
func isRegexPattern(pattern string) bool {
	return strings.HasPrefix(pattern, "regex:")
}

// trimRegexPrefix removes "regex:" prefix from pattern
func trimRegexPrefix(pattern string) string {
	if isRegexPattern(pattern) {
		return pattern[6:] // len("regex:")
	}
	return pattern
}

func compileRegexSafe(pattern string) *regexp.Regexp {
	regexCacheMu.RLock()
	if re, ok := regexCache[pattern]; ok {
		regexCacheMu.RUnlock()
		return re
	}
	regexCacheMu.RUnlock()

	re, err := regexp.Compile(pattern)
	if err != nil {
		os.Stderr.WriteString("Warning: invalid regex pattern '" + pattern + "': " + err.Error() + "\n")
		return nil
	}
	regexCacheMu.Lock()
	regexCache[pattern] = re
	regexCacheMu.Unlock()
	return re
}

// matchPattern checks if relPath matches a pattern (literal or regex with "regex:" prefix)
func matchPattern(relPath, pattern string) bool {
	if isRegexPattern(pattern) {
		re := compileRegexSafe(trimRegexPrefix(pattern))
		if re != nil && re.MatchString(relPath) {
			return true
		}
		return false
	}
	// literal match
	return matchesLiteral(relPath, pattern)
}

func ShouldIgnore(relPath, fullPath, root string, cfg Config) bool {
	relPath = filepath.ToSlash(relPath)

	if cfg.Gitignore {
		if isIgnoredByGitignore(relPath, fullPath, root) {
			return true
		}
	}

	for _, p := range cfg.Ignore {
		if matchPattern(relPath, p) {
			return true
		}
	}
	return false
}

func IsIncluded(relPath string, cfg Config) bool {
	if len(cfg.Include) == 0 {
		return true
	}
	// "*" literal matches everything
	for _, p := range cfg.Include {
		if p == "*" {
			return true
		}
	}
	for _, p := range cfg.Include {
		if matchPattern(relPath, p) {
			return true
		}
	}
	return false
}

func matchesLiteral(relPath, pattern string) bool {
	if relPath == pattern {
		return true
	}
	if strings.HasPrefix(relPath, pattern+"/") {
		return true
	}
	parts := strings.Split(relPath, "/")
	for _, part := range parts {
		if part == pattern {
			return true
		}
	}
	return false
}

func isIgnoredByGitignore(relPath, fullPath, root string) bool {
	parser := loadGitignoreParser(fullPath, root)
	if parser == nil {
		return false
	}
	return parser.MatchesPath(relPath)
}

func loadGitignoreParser(filePath, root string) ignore.IgnoreParser {
	dir := filepath.Dir(filePath)
	if dir == "" {
		dir = "."
	}
	dirAbs, _ := filepath.Abs(dir)

	gitignoreCacheMu.RLock()
	if p, ok := gitignoreCache[dirAbs]; ok {
		gitignoreCacheMu.RUnlock()
		return p
	}
	gitignoreCacheMu.RUnlock()

	// collect all .gitignore files from root to dir
	var allLines []string
	current := dirAbs
	for {
		gp := filepath.Join(current, ".gitignore")
		if data, err := os.ReadFile(gp); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "#") {
					allLines = append(allLines, line)
				}
			}
		}
		if current == root {
			break
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	parser := ignore.CompileIgnoreLines(allLines...)

	gitignoreCacheMu.Lock()
	gitignoreCache[dirAbs] = parser
	gitignoreCacheMu.Unlock()

	return parser
}

func filterEntries(entries []os.DirEntry, root, dirPath string, cfg Config) []os.DirEntry {
	result := []os.DirEntry{}
	for _, e := range entries {
		fullPath := filepath.Join(dirPath, e.Name())
		relPath, _ := filepath.Rel(root, fullPath)
		if !IsIncluded(relPath, cfg) {
			continue
		}
		if ShouldIgnore(relPath, fullPath, root, cfg) {
			continue
		}
		result = append(result, e)
	}
	return result
}
