package cmd

import (
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/mgranderath/rlcs-cli/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestListTournamentsCmd_matchesFilters(t *testing.T) {
	today := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		cmd      ListTournamentsCmd
		tour     domain.Tournament
		expected bool
	}{
		{
			name:     "no filters - match all",
			cmd:      ListTournamentsCmd{},
			tour:     domain.Tournament{Region: domain.RegionNA, IsOnline: true},
			expected: true,
		},
		{
			name:     "region filter - match",
			cmd:      ListTournamentsCmd{Region: "NA"},
			tour:     domain.Tournament{Region: domain.RegionNA},
			expected: true,
		},
		{
			name:     "region filter - no match",
			cmd:      ListTournamentsCmd{Region: "NA"},
			tour:     domain.Tournament{Region: domain.RegionEU},
			expected: false,
		},
		{
			name:     "region filter - case insensitive",
			cmd:      ListTournamentsCmd{Region: "eu"},
			tour:     domain.Tournament{Region: domain.RegionEU},
			expected: true,
		},
		{
			name:     "online filter - match",
			cmd:      ListTournamentsCmd{Online: true},
			tour:     domain.Tournament{IsOnline: true},
			expected: true,
		},
		{
			name:     "online filter - no match",
			cmd:      ListTournamentsCmd{Online: true},
			tour:     domain.Tournament{IsOnline: false},
			expected: false,
		},
		{
			name:     "major filter - match",
			cmd:      ListTournamentsCmd{Major: true},
			tour:     domain.Tournament{IsMajor: true},
			expected: true,
		},
		{
			name:     "major filter - no match",
			cmd:      ListTournamentsCmd{Major: true},
			tour:     domain.Tournament{IsMajor: false},
			expected: false,
		},
		{
			name:     "grouping filter - partial match",
			cmd:      ListTournamentsCmd{Grouping: "Open 1"},
			tour:     domain.Tournament{Name: "RLCS Open 1 EU 2026"},
			expected: true,
		},
		{
			name:     "grouping filter - no match",
			cmd:      ListTournamentsCmd{Grouping: "Open 2"},
			tour:     domain.Tournament{Name: "RLCS Open 1 EU 2026"},
			expected: false,
		},
		{
			name:     "min teams filter - match",
			cmd:      ListTournamentsCmd{MinTeams: 16},
			tour:     domain.Tournament{TeamCount: 24},
			expected: true,
		},
		{
			name:     "min teams filter - no match",
			cmd:      ListTournamentsCmd{MinTeams: 16},
			tour:     domain.Tournament{TeamCount: 8},
			expected: false,
		},
		{
			name:     "upcoming filter - match",
			cmd:      ListTournamentsCmd{Upcoming: true},
			tour:     domain.Tournament{StartDate: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)},
			expected: true,
		},
		{
			name:     "upcoming filter - no match",
			cmd:      ListTournamentsCmd{Upcoming: true},
			tour:     domain.Tournament{StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
			expected: false,
		},
		{
			name:     "ongoing filter - match",
			cmd:      ListTournamentsCmd{Ongoing: true},
			tour:     domain.Tournament{StartDate: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)},
			expected: true,
		},
		{
			name:     "past filter - match",
			cmd:      ListTournamentsCmd{Past: true},
			tour:     domain.Tournament{EndDate: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)},
			expected: true,
		},
		{
			name:     "past filter - no match",
			cmd:      ListTournamentsCmd{Past: true},
			tour:     domain.Tournament{EndDate: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)},
			expected: false,
		},
		{
			name:     "multiple filters - all match",
			cmd:      ListTournamentsCmd{Region: "NA", Online: true},
			tour:     domain.Tournament{Region: domain.RegionNA, IsOnline: true},
			expected: true,
		},
		{
			name:     "multiple filters - one fails",
			cmd:      ListTournamentsCmd{Region: "NA", Online: true},
			tour:     domain.Tournament{Region: domain.RegionNA, IsOnline: false},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.matchesFilters(tt.tour, today)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestListTournamentsCmd_matchesFilters_ConflictingTemporal(t *testing.T) {
	today := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	// Test that upcoming and past are mutually exclusive in filter logic
	t.Run("upcoming and past both true should not match anything", func(t *testing.T) {
		cmd := &ListTournamentsCmd{Upcoming: true, Past: true}
		tour := domain.Tournament{
			StartDate: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		result := cmd.matchesFilters(tour, today)
		assert.False(t, result)
	})
}

func TestListTournamentsCmd_Run_HTTPMock(t *testing.T) {
	defer gock.Off()

	t.Run("successful fetch with default circuit", func(t *testing.T) {
		// Mock the API response - use fixed time for deterministic testing
		gock.New("https://api.blast.tv").
			Get("/v2/circuits/2026/tournaments").
			MatchParam("game", "rl").
			Reply(200).
			JSON([]map[string]interface{}{
				{
					"id":            "test-tournament",
					"name":          "Test Tournament",
					"startDate":     "2026-01-15",
					"endDate":       "2026-01-17",
					"circuitId":     "2026",
					"region":        "NA",
					"numberOfTeams": 16,
				},
			})

		// Use fixed time injection to make test deterministic
		cmd := &ListTournamentsCmd{
			Output: "table",
			now: func() time.Time {
				return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
			},
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.NoError(t, err)
		assert.True(t, gock.IsDone())
	})

	t.Run("successful fetch with explicit circuit", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/circuits/2025/tournaments").
			MatchParam("game", "rl").
			Reply(200).
			JSON([]map[string]interface{}{
				{
					"id":            "old-tournament",
					"name":          "Old Tournament",
					"startDate":     "2025-01-15",
					"endDate":       "2025-01-17",
					"circuitId":     "2025",
					"region":        "EU",
					"numberOfTeams": 16,
				},
			})

		cmd := &ListTournamentsCmd{Circuit: "2025", Output: "table"}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.NoError(t, err)
		assert.True(t, gock.IsDone())
	})

	t.Run("404 response", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/circuits/2026/tournaments").
			Reply(404)

		// Use fixed time injection to make test deterministic
		cmd := &ListTournamentsCmd{
			Output: "table",
			now: func() time.Time {
				return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
			},
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code: 404")
	})

	t.Run("500 response", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/circuits/2026/tournaments").
			Reply(500)

		// Use fixed time injection to make test deterministic
		cmd := &ListTournamentsCmd{
			Output: "table",
			now: func() time.Time {
				return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
			},
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code: 500")
	})
}

func TestListTournamentsCmd_Run_Validation(t *testing.T) {
	t.Run("conflicting temporal filters should error", func(t *testing.T) {
		cmd := &ListTournamentsCmd{Upcoming: true, Past: true}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot use --upcoming and --past together")
	})
}
