package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/mgranderath/rlcs-cli/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchesTableFormatter_Format(t *testing.T) {
	formatter := &MatchesTableFormatter{}

	tests := []struct {
		name     string
		matches  []domain.Match
		contains []string
	}{
		{
			name: "single match",
			matches: []domain.Match{
				{
					UUID:         "match-1",
					Name:         "Grand Final",
					Type:         "BO7",
					TeamA:        domain.MatchTeam{Name: "Vitality"},
					TeamB:        domain.MatchTeam{Name: "KC"},
					TeamAScore:   4,
					TeamBScore:   2,
					IsCompleted:  true,
					TimeOfSeries: time.Date(2026, 1, 15, 18, 0, 0, 0, time.UTC),
				},
			},
			contains: []string{"Grand Final", "Vitality vs KC", "4 - 2", "Completed"},
		},
		{
			name: "multiple matches",
			matches: []domain.Match{
				{
					UUID:        "m1",
					Name:        "Match 1",
					TeamA:       domain.MatchTeam{Name: "Team A"},
					TeamB:       domain.MatchTeam{Name: "Team B"},
					TeamAScore:  3,
					TeamBScore:  1,
					IsCompleted: true,
				},
				{
					UUID:        "m2",
					Name:        "Match 2",
					TeamA:       domain.MatchTeam{Name: "Team C"},
					TeamB:       domain.MatchTeam{Name: "Team D"},
					TeamAScore:  2,
					TeamBScore:  3,
					IsCompleted: true,
				},
			},
			contains: []string{"Match 1", "Match 2", "Team A vs Team B", "Team C vs Team D"},
		},
		{
			name:     "empty matches",
			matches:  []domain.Match{},
			contains: []string{"No matches found"},
		},
		{
			name: "live match",
			matches: []domain.Match{
				{
					UUID:        "live-match",
					Name:        "Current Match",
					TeamA:       domain.MatchTeam{Name: "Team A"},
					TeamB:       domain.MatchTeam{Name: "Team B"},
					IsLive:      true,
					IsCompleted: false,
				},
			},
			contains: []string{"Current Match", "LIVE"},
		},
		{
			name: "upcoming match",
			matches: []domain.Match{
				{
					UUID:        "upcoming-match",
					Name:        "Future Match",
					TeamA:       domain.MatchTeam{Name: "Team A"},
					TeamB:       domain.MatchTeam{Name: "Team B"},
					IsLive:      false,
					IsCompleted: false,
				},
			},
			contains: []string{"Future Match", "Upcoming"},
		},
		{
			name: "match with type",
			matches: []domain.Match{
				{
					UUID:        "m1",
					Name:        "Semi Final",
					Type:        "BO5",
					TeamA:       domain.MatchTeam{Name: "Vitality"},
					TeamB:       domain.MatchTeam{Name: "KC"},
					TeamAScore:  3,
					TeamBScore:  2,
					IsCompleted: true,
				},
			},
			contains: []string{"Semi Final", "Vitality vs KC", "Completed"},
		},
		{
			name: "match with zero scores",
			matches: []domain.Match{
				{
					UUID:        "m1",
					Name:        "Upcoming Final",
					TeamA:       domain.MatchTeam{Name: "Team A"},
					TeamB:       domain.MatchTeam{Name: "Team B"},
					TeamAScore:  0,
					TeamBScore:  0,
					IsLive:      false,
					IsCompleted: false,
				},
			},
			contains: []string{"Upcoming Final", "0 - 0", "Upcoming"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := formatter.Format(&buf, tt.matches)
			require.NoError(t, err)

			output := buf.String()
			for _, expected := range tt.contains {
				assert.Contains(t, output, expected)
			}
		})
	}
}

func TestMatchesTableFormatter_FormatStatus(t *testing.T) {
	formatter := &MatchesTableFormatter{}

	tests := []struct {
		name     string
		match    domain.Match
		expected string
	}{
		{
			name:     "live match",
			match:    domain.Match{IsLive: true, IsCompleted: false},
			expected: "LIVE",
		},
		{
			name:     "completed match",
			match:    domain.Match{IsLive: false, IsCompleted: true},
			expected: "Completed",
		},
		{
			name:     "upcoming match",
			match:    domain.Match{IsLive: false, IsCompleted: false},
			expected: "Upcoming",
		},
		{
			name:     "edge case - both live and completed",
			match:    domain.Match{IsLive: true, IsCompleted: true},
			expected: "LIVE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.formatStatus(tt.match)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMatchesTableFormatter_LongNames(t *testing.T) {
	formatter := &MatchesTableFormatter{}
	matches := []domain.Match{
		{
			UUID:       "m1",
			Name:       "Very Long Match Name That Exceeds Normal Limits",
			TeamA:      domain.MatchTeam{Name: "Very Long Team Name A"},
			TeamB:      domain.MatchTeam{Name: "Very Long Team Name B"},
			TeamAScore: 3,
			TeamBScore: 1,
		},
	}

	var buf bytes.Buffer
	err := formatter.Format(&buf, matches)
	require.NoError(t, err)

	output := buf.String()
	// Verify table structure is maintained even with long names
	lines := strings.Split(output, "\n")
	assert.True(t, len(lines) > 3)
	// Check that truncation happened in match fields (teams and match name)
	assert.Contains(t, output, "...")
	// Verify the match is still displayed
	assert.Contains(t, output, "Very Long")
}
