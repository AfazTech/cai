package core

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

func getCurrentProject() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	cwdAbs, _ := normalizePath(cwd)
	names, _ := ListProjectNames()
	for _, name := range names {
		cfg, err := LoadProject(name)
		if err != nil {
			continue
		}
		rootAbs, _ := normalizePath(cfg.RootPath)
		if rootAbs == cwdAbs {
			return name, nil
		}
	}
	return "", fmt.Errorf("no project found for current directory")
}

func runInit(name string) error {
	if _, err := LoadProject(name); err == nil {
		return fmt.Errorf("project '%s' already exists", name)
	}
	root, _ := os.Getwd()
	cfg := DefaultConfig(name, root)
	return SaveProject(name, cfg)
}

func runSetDescription(proj, desc string) error {
	cfg, err := LoadProject(proj)
	if err != nil {
		return err
	}
	cfg.Description = desc
	return SaveProject(proj, cfg)
}

func runSetSize(proj string, size int) error {
	if size <= 0 {
		return fmt.Errorf("size must be positive")
	}
	cfg, err := LoadProject(proj)
	if err != nil {
		return err
	}
	cfg.MaxSizeMB = size
	return SaveProject(proj, cfg)
}

func runListProjects() error {
	names, _ := ListProjectNames()
	if len(names) == 0 {
		fmt.Println("No projects. Run 'cai init <name>'.")
		return nil
	}
	fmt.Println("Projects:")
	for _, n := range names {
		cfg, _ := LoadProject(n)
		fmt.Printf("  %s (root: %s)\n", n, cfg.RootPath)
	}
	return nil
}

func runAddIgnore(proj, pattern string) error {
	cfg, err := LoadProject(proj)
	if err != nil {
		return err
	}
	if slices.Contains(cfg.Ignore, pattern) {
		fmt.Printf("pattern '%s' already in ignore list\n", pattern)
		return nil
	}
	cfg.Ignore = append(cfg.Ignore, pattern)
	return SaveProject(proj, cfg)
}

func runRemoveIgnore(proj, pattern string) error {
	cfg, err := LoadProject(proj)
	if err != nil {
		return err
	}
	orig := len(cfg.Ignore)
	cfg.Ignore = slices.DeleteFunc(cfg.Ignore, func(s string) bool { return s == pattern })
	if len(cfg.Ignore) == orig {
		fmt.Printf("pattern '%s' not found in ignore list\n", pattern)
	}
	return SaveProject(proj, cfg)
}

func runAddInclude(proj, pattern string) error {
	cfg, err := LoadProject(proj)
	if err != nil {
		return err
	}
	if slices.Contains(cfg.Include, pattern) {
		fmt.Printf("pattern '%s' already in include list\n", pattern)
		return nil
	}
	cfg.Include = append(cfg.Include, pattern)
	return SaveProject(proj, cfg)
}

func runRemoveInclude(proj, pattern string) error {
	cfg, err := LoadProject(proj)
	if err != nil {
		return err
	}
	orig := len(cfg.Include)
	cfg.Include = slices.DeleteFunc(cfg.Include, func(s string) bool { return s == pattern })
	if len(cfg.Include) == orig {
		fmt.Printf("pattern '%s' not found in include list\n", pattern)
	}
	return SaveProject(proj, cfg)
}

var groups = map[string][]string{
	"frontend": {".html", ".htm", ".xhtml", ".xml", ".css", ".scss", ".sass", ".less", ".styl",
		".js", ".mjs", ".cjs", ".jsx", ".ts", ".tsx", ".vue", ".svelte", ".astro", ".mdx"},
	"markup": {".html", ".htm", ".xhtml", ".xml"},
	"styles": {".css", ".scss", ".sass", ".less", ".styl"},
	"js":     {".js", ".mjs", ".cjs", ".jsx"},
	"ts":     {".ts", ".tsx"},
	"react":  {".jsx", ".tsx"},
	"vue":    {".vue"},
	"svelte": {".svelte"},
	"astro":  {".astro"},
	"php":    {".php", ".php3", ".php4", ".php5", ".php7", ".phtml", ".inc"},
	"python": {".py", ".pyi", ".pyx", ".pxd", ".pyd", ".pyw"},
	"c":      {".c", ".h"},
	"cpp":    {".cpp", ".cc", ".cxx", ".c++", ".hpp", ".hh", ".hxx", ".h++", ".inl"},
	"csharp": {".cs", ".csx"},
	"java":   {".java"},
	"go":     {".go"},
	"rust":   {".rs"},
	"ruby":   {".rb", ".rbw", ".rake", ".gemspec"},
	"swift":  {".swift"},
	"kotlin": {".kt", ".kts"},
	"scala":  {".scala", ".sc"},
	"perl":   {".pl", ".pm", ".t"},
	"lua":    {".lua"},
	"r":      {".r", ".R", ".Rmd", ".Rnw"},
	"shell":  {".sh", ".bash", ".zsh", ".fish", ".csh", ".tcsh", ".ksh"},
	"sql":    {".sql"},
	"docker": {"Dockerfile", ".dockerignore"},
	"yaml":   {".yaml", ".yml"},
	"json":   {".json", ".jsonc"},
	"toml":   {".toml"},
	"ini":    {".ini", ".cfg", ".conf"},
	"make":   {"Makefile", "makefile", ".mk"},
	"cmake":  {"CMakeLists.txt", ".cmake"},
}

func groupNames() string {
	keys := []string{}
	for k := range groups {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return strings.Join(keys, ", ")
}

func runAddGroup(proj, group string) error {
	patterns, ok := groups[group]
	if !ok {
		return fmt.Errorf("unknown group '%s'. available: %s", group, groupNames())
	}
	cfg, err := LoadProject(proj)
	if err != nil {
		return err
	}
	added := 0
	for _, p := range patterns {
		if !slices.Contains(cfg.Include, p) {
			cfg.Include = append(cfg.Include, p)
			added++
		}
	}
	if err := SaveProject(proj, cfg); err != nil {
		return err
	}
	fmt.Printf("added %d patterns from group '%s'\n", added, group)
	return nil
}

func runRemoveGroup(proj, group string) error {
	patterns, ok := groups[group]
	if !ok {
		return fmt.Errorf("unknown group '%s'. available: %s", group, groupNames())
	}
	cfg, err := LoadProject(proj)
	if err != nil {
		return err
	}
	removed := 0
	for _, p := range patterns {
		old := len(cfg.Include)
		cfg.Include = slices.DeleteFunc(cfg.Include, func(s string) bool { return s == p })
		if len(cfg.Include) < old {
			removed++
		}
	}
	if err := SaveProject(proj, cfg); err != nil {
		return err
	}
	fmt.Printf("removed %d patterns of group '%s'\n", removed, group)
	return nil
}
