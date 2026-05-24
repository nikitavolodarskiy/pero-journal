package main

import (
	"context"

	"pero/internal/config"
	"pero/internal/journal"
	"pero/internal/service"
)

// App is the Wails application struct.
// Every exported method on App is automatically bound to the Svelte frontend.
type App struct {
	ctx context.Context
	svc service.Service
}

func NewApp() *App {
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load pero config: " + err.Error())
	}
	return &App{svc: service.New(cfg)}
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
