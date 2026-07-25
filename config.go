package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the configuration loaded from cai.json
type Config struct {
	Project     string   `json:"project"`
	Tree        bool     `json:"tree"`
	Files       bool     `json:"files"`
	Ignore      []string `json:"ignore"`
	Include     []string `json:"include"`
	Gitignore   bool     `json:"gitignore"`
	Description string   `json:"description"` // optional project description
}

// LoadConfig reads cai.json from root, or returns defaults if not present
func LoadConfig(root string) (Config, error) {
	config := Config{
		Project:   filepath.Base(root),
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
		},
		Include:     []string{"*"},
		Gitignore:   false,
		Description: "",
	}

	configPath := filepath.Join(root, "cai.json")
	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return config, nil
	}
	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("invalid cai.json: %w", err)
	}

	return config, nil
}
