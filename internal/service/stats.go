package service

import (
	"strings"
	"time"

	"pero/internal/journal"
)

// Stats holds aggregated statistics about the journal.
// All fields are safe to display directly in any interface.
type Stats struct {
	TotalEntries  int
	TotalWords    int
	AverageWords  int
	CurrentStreak int // consecutive days up to and including today
	LongestStreak int // all-time best consecutive-day run
	FirstEntry    *time.Time
	LastEntry     *time.Time
}

// Stats implements Service.
//
// Results are cached between writes: the first call after startup (or after
// any WriteEntry/Today call) reads all files; subsequent calls return the
// cached value immediately.
func (s *svc) Stats() (Stats, error) {
	// Fast path: return cached value if still valid.
	s.mu.Lock()
	if s.cached != nil {
		c := *s.cached
		s.mu.Unlock()
		return c, nil
	}
	s.mu.Unlock()

	entries, err := journal.ListEntries(s.cfg.JournalDir)
	if err != nil {
		return Stats{}, err
	}

	if len(entries) == 0 {
		return Stats{}, nil
	}

	totalWords := 0
	for _, e := range entries {
		content, err := journal.ReadContent(e.Path)
		if err != nil {
			return Stats{}, err
		}
		totalWords += countWords(content)
	}

	avg := totalWords / len(entries)
	current, longest := computeStreaks(entries)

	// entries are newest-first; last chronologically is the final element.
	first := entries[len(entries)-1].Date
	last := entries[0].Date

	result := Stats{
		TotalEntries:  len(entries),
		TotalWords:    totalWords,
		AverageWords:  avg,
		CurrentStreak: current,
		LongestStreak: longest,
		FirstEntry:    &first,
		LastEntry:     &last,
	}

	s.mu.Lock()
	s.cached = &result
	s.mu.Unlock()

	return result, nil
}

// countWords returns the number of whitespace-separated tokens in s.
// Works correctly for markdown content (counts words across lines).
func countWords(s string) int {
	return len(strings.Fields(s))
}

// computeStreaks returns the current and longest consecutive-day streaks
// derived from the given entries (which may be in any order).
//
// Current streak: consecutive days running up to and including today.
// If today has no entry, the current streak is 0.
// Longest streak: the longest consecutive-day run across all entries.
func computeStreaks(entries []journal.Entry) (current, longest int) {
	if len(entries) == 0 {
		return 0, 0
	}

	// Build a date-keyed set for O(1) lookup (yyyy-mm-dd → true).
	dateSet := make(map[string]bool, len(entries))
	for _, e := range entries {
		dateSet[e.DateString()] = true
	}

	today := truncateToDay(time.Now())

	// Current streak: walk backwards from today while each day has an entry.
	current = 0
	for d := today; ; d = d.AddDate(0, 0, -1) {
		if dateSet[d.Format("2006-01-02")] {
			current++
		} else {
			break
		}
	}

	// Longest streak: walk all entries in chronological order, count runs.
	// entries is newest-first, so iterate in reverse.
	longest = 0
	run := 1
	for i := len(entries) - 1; i > 0; i-- {
		prev := truncateToDay(entries[i].Date)
		next := truncateToDay(entries[i-1].Date)
		if prev.AddDate(0, 0, 1).Equal(next) {
			run++
		} else {
			if run > longest {
				longest = run
			}
			run = 1
		}
	}
	// Flush the last run.
	if run > longest {
		longest = run
	}

	return current, longest
}

// truncateToDay strips the time component from t, returning midnight local time.
func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}
