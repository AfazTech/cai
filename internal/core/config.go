package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Project     string   `json:"project"`
	RootPath    string   `json:"rootPath"`
	Tree        bool     `json:"tree"`
	Files       bool     `json:"files"`
	Ignore      []string `json:"ignore"`   // literal patterns, or "regex:..." for regex
	Include     []string `json:"include"`  // literal patterns, or "regex:..." for regex
	Gitignore   bool     `json:"gitignore"`
	Description string   `json:"description"`
	MaxSizeMB   int      `json:"maxSizeMB"`
}

func DefaultConfig(projectName, rootPath string) Config {
	return Config{
		Project:   projectName,
		RootPath:  rootPath,
		Tree:      true,
		Files:     true,
		Ignore: []string{
			".git",
			"vendor",
			"node_modules",
			"dist",
			"build",
			"storage",
			".env",
			"cai.json",
			"go.mod",
			"go.sum",
			".cai",
		},
		Include:     []string{"*"},
		Gitignore:   false,
		Description: "",
		MaxSizeMB:   50,
	}
}

func getCAIDir() string      { return ".cai" }
func getProjectsDir() string { return filepath.Join(getCAIDir(), "projects") }
func getProjectConfigPath(name string) string {
	return filepath.Join(getProjectsDir(), name, "config.json")
}

func LoadProject(name string) (Config, error) {
	path := getProjectConfigPath(name)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Config{}, fmt.Errorf("project '%s' not found. run 'cai init %s'", name, name)
	}
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("invalid config: %w", err)
	}
	return cfg, nil
}

func SaveProject(name string, cfg Config) error {
	dir := filepath.Dir(getProjectConfigPath(name))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(getProjectConfigPath(name), data, 0644)
}

func ListProjectNames() ([]string, error) {
	dir := getProjectsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	names := []string{}
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func normalizePath(p string) (string, error) {
	if p == "" {
		return "", fmt.Errorf("empty path")
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return filepath.Clean(abs), nil
}
