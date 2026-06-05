package main

import (
	"context"
	"fmt"
	"sync"

	"pero/internal/config"
	"pero/internal/journal"
	"pero/internal/service"
)

// App is the Wails application struct.
// Every exported method on App is automatically bound to the Svelte frontend.
type App struct {
	ctx context.Context
	mu  sync.Mutex // guards cfg and svc; only SaveConfig writes them after init
	cfg *config.Config
	svc service.Service
}

func NewApp() *App {
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load pero config: " + err.Error())
	}
	return &App{cfg: cfg, svc: service.New(cfg)}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetToday returns (or creates) today's journal entry.
func (a *App) GetToday() (journal.Entry, error) {
	return a.svc.Today()
}

// GetEntries returns all journal entries, newest-first.
func (a *App) GetEntries() ([]journal.Entry, error) {
	return a.svc.List()
}

// GetEntry returns an existing entry for the given date (yyyy-mm-dd).
func (a *App) GetEntry(date string) (journal.Entry, error) {
	return a.svc.Get(date)
}

// ReadEntry returns the text content of the entry for the given date.
func (a *App) ReadEntry(date string) (string, error) {
	entry, err := a.svc.Get(date)
	if err != nil {
		return "", err
	}
	return a.svc.ReadEntry(entry)
}

// WriteEntry saves content for the given date.
func (a *App) WriteEntry(date string, content string) error {
	entry, err := a.svc.Get(date)
	if err != nil {
		return err
	}
	return a.svc.WriteEntry(entry, content)
}

// GetStats returns aggregated journal statistics.
func (a *App) GetStats() (service.Stats, error) {
	return a.svc.Stats()
}

// GetConfig returns a snapshot of the current configuration.
func (a *App) GetConfig() config.Config {
	a.mu.Lock()
	cfg := *a.cfg
	a.mu.Unlock()
	return cfg
}

// SaveConfig writes new settings to disk and rebuilds the service so changes
// take effect immediately without restarting the app.
// We do NOT call Load() here because Load() re-applies the PERO_JOURNAL_DIR
// env-var override, which would ignore the user's explicit choice.
func (a *App) SaveConfig(journalDir, editor string) error {
	if err := config.Save(journalDir, editor); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	resolved, err := config.ResolvePath(journalDir)
	if err != nil {
		return fmt.Errorf("resolve journal dir: %w", err)
	}
	cfg := &config.Config{JournalDir: resolved, Editor: editor}
	a.mu.Lock()
	a.cfg = cfg
	a.svc = service.New(cfg)
	a.mu.Unlock()
	return nil
}
