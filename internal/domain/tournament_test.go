package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTournament_IsUpcoming(t *testing.T) {
	tests := []struct {
		name       string
		tournament Tournament
		now        time.Time
		expected   bool
	}{
		{
			name: "upcoming - start date in future",
			tournament: Tournament{
				StartDate: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 1, 17, 0, 0, 0, 0, time.UTC),
			},
			now:      time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name: "not upcoming - already started",
			tournament: Tournament{
				StartDate: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 1, 17, 0, 0, 0, 0, time.UTC),
			},
			now:      time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name: "not upcoming - same day as start",
			tournament: Tournament{
				StartDate: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 1, 17, 0, 0, 0, 0, time.UTC),
			},
			now:      time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tournament.IsUpcoming(tt.now)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTournament_IsOngoing(t *testing.T) {
	tests := []struct {
		name       string
		tournament Tournament
		now        time.Time
		expected   bool
	}{
		{
			name: "ongoing - in middle of tournament",
			tournament: Tournament{
				StartDate: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC),
			},
			now:      time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name: "ongoing - on start date",
			tournament: Tournament{
				StartDate: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC),
			},
			now:      time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name: "ongoing - on end date",
			tournament: Tournament{
				StartDate: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			now:      time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name: "not ongoing - before start",
			tournament: Tournament{
				StartDate: time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC),
			},
			now:      time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name: "not ongoing - after end",
			tournament: Tournament{
				StartDate: time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
			},
			now:      time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tournament.IsOngoing(tt.now)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTournament_IsPast(t *testing.T) {
	tests := []struct {
		name       string
		tournament Tournament
		now        time.Time
		expected   bool
	}{
		{
			name: "past - tournament ended",
			tournament: Tournament{
				StartDate: time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
			},
			now:      time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name: "not past - still ongoing",
			tournament: Tournament{
				StartDate: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC),
			},
			now:      time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name: "not past - same day as end",
			tournament: Tournament{
				StartDate: time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			now:      time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tournament.IsPast(tt.now)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTournament_Duration(t *testing.T) {
	tests := []struct {
		name       string
		tournament Tournament
		expected   time.Duration
	}{
		{
			name: "multi-day tournament",
			tournament: Tournament{
				StartDate: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			expected: 5 * 24 * time.Hour,
		},
		{
			name: "single day tournament",
			tournament: Tournament{
				StartDate: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tournament.Duration()
			assert.Equal(t, tt.expected, result)
		})
	}
}
