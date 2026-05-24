package service

import (
	"testing"
	"time"

	"pero/internal/journal"
)

func TestComputeStreaks(t *testing.T) {
	// Helper to build an Entry with just a date (path not needed for streak logic).
	entry := func(dateStr string) journal.Entry {
		d, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			t.Fatalf("bad date in test: %s", dateStr)
		}
		return journal.Entry{Date: d}
	}

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	twoDaysAgo := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	threeDaysAgo := time.Now().AddDate(0, 0, -3).Format("2006-01-02")

	tests := []struct {
		name            string
		entries         []journal.Entry // newest-first, as ListEntries returns
		wantCurrent     int
		wantLongest     int
	}{
		{
			name:        "empty",
			entries:     nil,
			wantCurrent: 0,
			wantLongest: 0,
		},
		{
			name:        "only today",
			entries:     []journal.Entry{entry(today)},
			wantCurrent: 1,
			wantLongest: 1,
		},
		{
			name:        "today and yesterday",
			entries:     []journal.Entry{entry(today), entry(yesterday)},
			wantCurrent: 2,
			wantLongest: 2,
		},
		{
			name: "three consecutive days ending today",
			entries: []journal.Entry{
				entry(today), entry(yesterday), entry(twoDaysAgo),
			},
			wantCurrent: 3,
			wantLongest: 3,
		},
		{
			name: "gap breaks current streak",
			entries: []journal.Entry{
				// today is missing — gap before yesterday
				entry(yesterday), entry(twoDaysAgo), entry(threeDaysAgo),
			},
			wantCurrent: 0,
			wantLongest: 3,
		},
		{
			name: "current streak shorter than longest",
			entries: []journal.Entry{
				entry(today), entry(yesterday),
				// gap here
				entry(threeDaysAgo), entry(time.Now().AddDate(0, 0, -4).Format("2006-01-02")),
				entry(time.Now().AddDate(0, 0, -5).Format("2006-01-02")),
			},
			wantCurrent: 2,
			wantLongest: 3,
		},
		{
			name: "single old entry — no current streak",
			entries: []journal.Entry{
				entry("2025-01-01"),
			},
			wantCurrent: 0,
			wantLongest: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotCurrent, gotLongest := computeStreaks(tc.entries)
			if gotCurrent != tc.wantCurrent {
				t.Errorf("current streak: got %d, want %d", gotCurrent, tc.wantCurrent)
			}
			if gotLongest != tc.wantLongest {
				t.Errorf("longest streak: got %d, want %d", gotLongest, tc.wantLongest)
			}
		})
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"", 0},
		{"hello", 1},
		{"hello world", 2},
		{"  lots   of   spaces  ", 3},
		{"line one\nline two\nline three", 6},
		{"# Heading\n\nSome text here.", 5}, // '#' counts as a token
	}

	for _, tc := range tests {
		got := countWords(tc.input)
		if got != tc.want {
			t.Errorf("countWords(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}
}
