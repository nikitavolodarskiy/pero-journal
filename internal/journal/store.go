package journal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// EnsureEntry creates the full directory hierarchy and the entry file if they
// do not already exist. The file is created blank (no template).
func EnsureEntry(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create journal directories: %w", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("create entry %s: %w", path, err)
		}
		if err := f.Close(); err != nil {
			return fmt.Errorf("close entry %s: %w", path, err)
		}
	}

	return nil
}

// EntryExists returns nil if the entry file at path exists, or a descriptive
// error if it does not. It is a low-level check; callers are responsible for
// constructing the path via EntryPath.
func EntryExists(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		name := strings.TrimSuffix(filepath.Base(path), ".md")
		return fmt.Errorf("no entry found for %s", name)
	}
	return fmt.Errorf("stat %s: %w", path, err)
}

// ListEntries walks journalDir and returns every file whose name matches the
// yyyy-mm-dd.md pattern. Results are sorted newest-first.
func ListEntries(journalDir string) ([]Entry, error) {
	if _, err := os.Stat(journalDir); os.IsNotExist(err) {
		return nil, nil
	}

	var entries []Entry

	err := filepath.Walk(journalDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		name := strings.TrimSuffix(filepath.Base(path), ".md")
		t, err := time.Parse("2006-01-02", name)
		if err != nil {
			return nil // skip files that don't match the naming convention
		}

		entries = append(entries, Entry{Date: t, Path: path})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk journal directory: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Date.After(entries[j].Date)
	})

	return entries, nil
}

// ReadContent returns the full text content of the entry at path.
func ReadContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read entry %s: %w", path, err)
	}
	return string(data), nil
}

// WriteContent saves content to the entry at path, replacing whatever was there.
func WriteContent(path, content string) error {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write entry %s: %w", path, err)
	}
	return nil
}
