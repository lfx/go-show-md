package config

import (
	"os"
	"testing"
	"time"
)

type mockFileInfo struct {
	name  string
	isDir bool
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return 0 }
func (m mockFileInfo) Mode() os.FileMode  { return 0 }
func (m mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() interface{}   { return nil }

func TestMergeIgnorePatterns(t *testing.T) {
	defaults := []string{".git", "node_modules"}
	global := []string{"*.tmp", ".DS_Store"}
	local := []string{"vendor", ".env"}
	configPats := []string{"ignore_me.md", ".git"} // .git is a duplicate

	merged := MergeIgnorePatterns(defaults, global, local, configPats)

	expected := []string{".git", "node_modules", "*.tmp", ".DS_Store", "vendor", ".env", "ignore_me.md"}

	if len(merged) != len(expected) {
		t.Fatalf("Expected %d patterns, got %d", len(expected), len(merged))
	}

	patternMap := make(map[string]bool)
	for _, p := range merged {
		patternMap[p] = true
	}

	for _, p := range expected {
		if !patternMap[p] {
			t.Errorf("Expected pattern %s to be in merged list", p)
		}
	}
}

func TestIsIgnored(t *testing.T) {
	patterns := []string{".git", "node_modules", "*.tmp", "secret.md"}

	tests := []struct {
		name     string
		fileName string
		isDir    bool
		ignored  bool
	}{
		{"ignored dir .git", ".git", true, true},
		{"ignored dir node_modules", "node_modules", true, true},
		{"ignored file by ext", "test.tmp", false, true},
		{"ignored file by name", "secret.md", false, true},
		{"allowed dir", "src", true, false},
		{"allowed file", "test.md", false, false},
		{"almost match", "node_module", true, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			info := mockFileInfo{name: tc.fileName, isDir: tc.isDir}
			got := IsIgnored(tc.fileName, info, patterns)
			if got != tc.ignored {
				t.Errorf("Expected IsIgnored for %s to be %v, got %v", tc.fileName, tc.ignored, got)
			}
		})
	}
}
