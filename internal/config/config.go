package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Config struct {
	WatchedDirectories []string `json:"watched_directories"`
	IgnoredPatterns    []string `json:"ignored_patterns"`
	Port               int      `json:"port"`
	Host               string   `json:"host"`
}

type FileInfo struct {
	Path        string    `json:"path"`
	Name        string    `json:"name"`
	Directory   string    `json:"directory"`
	ModifiedAt  time.Time `json:"modified_at"`
	Size        int64     `json:"size"`
}

const DefaultConfigPath = "config.json"

func Load(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = DefaultConfigPath
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				WatchedDirectories: []string{},
				IgnoredPatterns:    []string{},
				Port:               8080,
				Host:               "127.0.0.1",
			}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Port == 0 {
		cfg.Port = 8080
	}
	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}

	return &cfg, nil
}

func (c *Config) Save(configPath string) error {
	if configPath == "" {
		configPath = DefaultConfigPath
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func (c *Config) AddDirectory(dir string) bool {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return false
	}

	for _, d := range c.WatchedDirectories {
		if d == absDir {
			return false
		}
	}

	c.WatchedDirectories = append(c.WatchedDirectories, absDir)
	return true
}

func (c *Config) RemoveDirectory(dir string) bool {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return false
	}

	for i, d := range c.WatchedDirectories {
		if d == absDir {
			c.WatchedDirectories = append(c.WatchedDirectories[:i], c.WatchedDirectories[i+1:]...)
			return true
		}
	}

	return false
}

func (c *Config) ScanMarkdownFiles() ([]FileInfo, error) {
	var files []FileInfo
	seen := make(map[string]bool)

	globalPatterns := GetGlobalIgnorePatterns()

	for _, dir := range c.WatchedDirectories {
		localPatterns := GetLocalIgnorePatterns(dir)
		patterns := MergeIgnorePatterns(DefaultIgnoredPatterns, globalPatterns, localPatterns, c.IgnoredPatterns)

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if IsIgnored(path, info, patterns) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			if info.IsDir() {
				return nil
			}

			if strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
				if !seen[path] {
					seen[path] = true
					files = append(files, FileInfo{
						Path:       path,
						Name:       info.Name(),
						Directory:  filepath.Dir(path),
						ModifiedAt: info.ModTime(),
						Size:       info.Size(),
					})
				}
			}

			return nil
		})

		if err != nil {
			continue
		}
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].Directory == files[j].Directory {
			return files[i].Name < files[j].Name
		}
		return files[i].Directory < files[j].Directory
	})

	return files, nil
}

func (c *Config) IsPathAllowed(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	for _, dir := range c.WatchedDirectories {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			continue
		}

		if strings.HasPrefix(absPath, absDir) {
			return true
		}
	}

	return false
}
