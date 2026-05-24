// Package service is the single brain of pero.
//
// All business logic lives here. Every interface layer — CLI, TUI, GUI, or
// anything else — talks only to Service and never imports internal/journal or
// internal/config directly.
package service

import (
	"fmt"
	"sync"
	"time"

	"pero/internal/config"
	"pero/internal/journal"
)

// Service is the contract every UI layer programs against.
// Adding a new interface to pero means wiring up this, nothing more.
type Service interface {
	// Today returns the entry for the current day, creating the file and any
	// missing directories if they do not exist yet.
	Today() (journal.Entry, error)

	// Get returns an existing entry for the given date (yyyy-mm-dd).
	// Returns an error if no entry exists for that date.
	Get(date string) (journal.Entry, error)

	// List returns all entries sorted newest-first.
	List() ([]journal.Entry, error)

	// Stats computes and returns aggregated journal statistics.
	Stats() (Stats, error)

	// ReadEntry returns the text content of the given entry.
	ReadEntry(entry journal.Entry) (string, error)

	// WriteEntry saves content to the given entry, overwriting what was there.
	WriteEntry(entry journal.Entry, content string) error

	// Editor returns the configured editor binary (e.g. "nvim").
	// UI layers that need to launch an external editor use this.
	Editor() string
}

// svc is the concrete implementation backed by the local filesystem.
type svc struct {
	cfg    *config.Config
	mu     sync.Mutex // guards cached
	cached *Stats     // nil means dirty; recomputed on next Stats() call
}

// New constructs a Service from an already-loaded Config.
func New(cfg *config.Config) Service {
	return &svc{cfg: cfg}
}

// Today implements Service.
func (s *svc) Today() (journal.Entry, error) {
	now := time.Now()
	path := journal.EntryPath(s.cfg.JournalDir, now)

	if err := journal.EnsureEntry(path); err != nil {
		return journal.Entry{}, fmt.Errorf("ensure today's entry: %w", err)
	}

	// A new file may have been created — invalidate the stats cache.
	s.mu.Lock()
	s.cached = nil
	s.mu.Unlock()

	return journal.Entry{Date: now, Path: path}, nil
}

// Get implements Service.
func (s *svc) Get(date string) (journal.Entry, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return journal.Entry{}, fmt.Errorf("invalid date %q — expected yyyy-mm-dd", date)
	}

	path := journal.EntryPath(s.cfg.JournalDir, t)
	if err := journal.EntryExists(path); err != nil {
		return journal.Entry{}, err
	}

	return journal.Entry{Date: t, Path: path}, nil
}

// List implements Service.
func (s *svc) List() ([]journal.Entry, error) {
	return journal.ListEntries(s.cfg.JournalDir)
}

// ReadEntry implements Service.
func (s *svc) ReadEntry(entry journal.Entry) (string, error) {
	return journal.ReadContent(entry.Path)
}

// WriteEntry implements Service.
func (s *svc) WriteEntry(entry journal.Entry, content string) error {
	if err := journal.WriteContent(entry.Path, content); err != nil {
		return err
	}
	// Content changed — word counts are stale.
	s.mu.Lock()
	s.cached = nil
	s.mu.Unlock()
	return nil
}

// Editor implements Service.
func (s *svc) Editor() string {
	return s.cfg.Editor
}
