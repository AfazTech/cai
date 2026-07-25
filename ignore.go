package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// regexCache caches compiled regex patterns for performance
var regexCache = make(map[string]*regexp.Regexp)

// gitignoreCache stores parsed gitignore rules per directory
var gitignoreCache = make(map[string][]gitignoreRule)

type gitignoreRule struct {
	pattern string
	negate  bool // true if pattern starts with "!"
	dirOnly bool // true if pattern ends with "/"
	regex   *regexp.Regexp
}

// regexSpecialChars lists the characters (besides ".") whose presence in an
// ignore pattern indicates the user intended it as a regular expression.
// "." is intentionally excluded so plain names such as ".git" or ".env"
// are treated as literal path segments instead of "any character" regex.
const regexSpecialChars = `\+*?()|[]{}^$`

// isRegexPattern reports whether a pattern should be treated as a regex
// rather than a literal path segment.
func isRegexPattern(pattern string) bool {
	return strings.ContainsAny(pattern, regexSpecialChars)
}

// filterEntries removes entries that match any ignore pattern or gitignore,
// and also applies the include filter.
func filterEntries(entries []os.DirEntry, root, dirPath string, config Config) []os.DirEntry {
	result := make([]os.DirEntry, 0, len(entries))
	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		relPath, err := filepath.Rel(root, fullPath)
		if err != nil {
			relPath = entry.Name()
		}
		// 1. Apply include filter
		if !IsIncluded(relPath, config) {
			continue
		}
		// 2. Apply ignore filter
		if ShouldIgnore(relPath, fullPath, root, config) {
			continue
		}
		result = append(result, entry)
	}
	return result
}

// ShouldIgnore checks if a given path should be ignored.
// It applies:
//   1. .gitignore rules (if Gitignore is true)
//   2. Custom ignore patterns (regex or literal path segments)
func ShouldIgnore(relPath, fullPath, root string, config Config) bool {
	relPath = filepath.ToSlash(relPath)

	if root != "" {
		fullRel, err := filepath.Rel(root, fullPath)
		if err == nil {
			relPath = filepath.ToSlash(fullRel)
		}
	}

	// 1. Apply .gitignore if enabled
	if config.Gitignore {
		if isIgnoredByGitignore(relPath, fullPath, root) {
			return true
		}
	}

	// 2. Apply custom ignore patterns (literal segment match or regex)
	for _, pattern := range config.Ignore {
		pattern = filepath.ToSlash(pattern)

		if !isRegexPattern(pattern) {
			if matchesLiteralPattern(relPath, pattern) {
				return true
			}
			continue
		}

		re, ok := regexCache[pattern]
		if !ok {
			var err error
			re, err = regexp.Compile(pattern)
			if err != nil {
				re = nil
			}
			regexCache[pattern] = re
		}

		if re != nil && re.MatchString(relPath) {
			return true
		}
	}
	return false
}

// IsIncluded checks if a given path matches any include pattern.
// If Include list is empty or contains "*", all paths are included.
// Otherwise, the path must match at least one pattern (literal or regex).
func IsIncluded(relPath string, config Config) bool {
	if len(config.Include) == 0 {
		return true
	}
	for _, pattern := range config.Include {
		if pattern == "*" {
			return true
		}
		if !isRegexPattern(pattern) {
			if matchesLiteralPattern(relPath, pattern) {
				return true
			}
			continue
		}
		re, ok := regexCache[pattern]
		if !ok {
			var err error
			re, err = regexp.Compile(pattern)
			if err != nil {
				re = nil
			}
			regexCache[pattern] = re
		}
		if re != nil && re.MatchString(relPath) {
			return true
		}
	}
	return false
}

// matchesLiteralPattern checks an exact path or path-segment match for
// non-regex patterns, so a name such as "dist" only matches the
// "dist" segment itself and not paths like "mydist" or "distribution".
func matchesLiteralPattern(relPath, pattern string) bool {
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

// isIgnoredByGitignore checks if the given relative path matches any .gitignore rule.
// It loads .gitignore files from the directory hierarchy.
func isIgnoredByGitignore(relPath, fullPath, root string) bool {
	rules := loadGitignoreRules(fullPath, root)

	// Check each rule in order (last matching rule wins)
	match := false
	for _, rule := range rules {
		if rule.regex.MatchString(relPath) {
			match = !rule.negate
		}
	}
	return match
}

// loadGitignoreRules loads .gitignore rules from the current directory and parent directories.
// It returns a list of rules in order of precedence (parent first, child last).
func loadGitignoreRules(filePath, root string) []gitignoreRule {
	dir := filepath.Dir(filePath)
	if dir == "" {
		dir = "."
	}

	dirAbs, _ := filepath.Abs(dir)
	if rules, ok := gitignoreCache[dirAbs]; ok {
		return rules
	}

	var allRules []gitignoreRule

	relDir, err := filepath.Rel(root, dir)
	if err != nil {
		relDir = "."
	}
	parts := strings.Split(filepath.ToSlash(relDir), "/")
	if len(parts) == 0 || (len(parts) == 1 && parts[0] == ".") {
		parts = []string{}
	}

	pathParts := []string{root}
	for i := 0; i <= len(parts); i++ {
		var curDir string
		if i == 0 {
			curDir = root
		} else {
			pathParts = append(pathParts[:1], parts[:i]...)
			curDir = filepath.Join(pathParts...)
		}

		gitignorePath := filepath.Join(curDir, ".gitignore")
		if data, err := os.ReadFile(gitignorePath); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				negate := strings.HasPrefix(line, "!")
				if negate {
					line = line[1:]
				}
				dirOnly := strings.HasSuffix(line, "/")
				if dirOnly {
					line = line[:len(line)-1]
				}

				anchored := strings.HasPrefix(line, "/")
				matchLine := line
				if anchored {
					matchLine = line[1:]
				}

				regexStr := regexp.QuoteMeta(matchLine)
				regexStr = strings.ReplaceAll(regexStr, `\*`, "[^/]*")
				regexStr = strings.ReplaceAll(regexStr, `\?`, ".")

				if anchored {
					regexStr = "^" + regexStr
				} else {
					regexStr = "(^|.*/)" + regexStr
				}
				regexStr += "(/.*)?$"

				re, err := regexp.Compile(regexStr)
				if err != nil {
					continue
				}
				allRules = append(allRules, gitignoreRule{
					pattern: line,
					negate:  negate,
					dirOnly: dirOnly,
					regex:   re,
				})
			}
		}
	}

	gitignoreCache[dirAbs] = allRules
	return allRules
}
