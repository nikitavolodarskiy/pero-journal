package journal

import (
	"path/filepath"
	"time"
)

// Entry represents a single journal file on disk.
type Entry struct {
	Date time.Time
	Path string
}

// DateString returns the entry's date formatted as yyyy-mm-dd.
func (e Entry) DateString() string {
	return e.Date.Format("2006-01-02")
}

// EntryPath returns the absolute path for the journal entry on date t.
//
// Structure: <journalDir>/<year>/<month>/<yyyy-mm-dd>.md
// e.g. ~/pero/2026/05/2026-05-24.md
func EntryPath(journalDir string, t time.Time) string {
	year := t.Format("2006")
	month := t.Format("01")
	filename := t.Format("2006-01-02") + ".md"
	return filepath.Join(journalDir, year, month, filename)
}

// TodayPath returns the absolute path for today's journal entry.
func TodayPath(journalDir string) string {
	return EntryPath(journalDir, time.Now())
}
