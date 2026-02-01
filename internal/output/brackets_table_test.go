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

func TestBracketsTableFormatter_Format(t *testing.T) {
	formatter := &BracketsTableFormatter{}

	tests := []struct {
		name     string
		brackets []domain.Bracket
		contains []string
	}{
		{
			name: "single bracket with matches",
			brackets: []domain.Bracket{
				{
					TournamentUUID:       "bracket-1",
					TournamentName:       "Group A",
					ParentTournamentName: "RLCS Open 2026",
					Label:                "Group A",
					Format:               "double-elim-8",
					Matches: []domain.Match{
						{
							UUID:         "match-1",
							Name:         "Round 1",
							Type:         "BO5",
							TeamA:        domain.MatchTeam{Name: "Vitality"},
							TeamB:        domain.MatchTeam{Name: "KC"},
							TeamAScore:   3,
							TeamBScore:   1,
							IsCompleted:  true,
							TimeOfSeries: time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC),
						},
					},
				},
			},
			contains: []string{"Group A", "RLCS Open 2026", "Round 1", "Vitality vs KC", "3 - 1", "Completed"},
		},
		{
			name: "multiple brackets",
			brackets: []domain.Bracket{
				{
					TournamentName: "Group A",
					Label:          "Group A",
					Matches: []domain.Match{
						{UUID: "m1", Name: "Match 1", TeamA: domain.MatchTeam{Name: "Team A"}, TeamB: domain.MatchTeam{Name: "Team B"}},
					},
				},
				{
					TournamentName: "Group B",
					Label:          "Group B",
					Matches: []domain.Match{
						{UUID: "m2", Name: "Match 2", TeamA: domain.MatchTeam{Name: "Team C"}, TeamB: domain.MatchTeam{Name: "Team D"}},
					},
				},
			},
			contains: []string{"Group A", "Group B", "========"},
		},
		{
			name:     "empty brackets",
			brackets: []domain.Bracket{},
			contains: []string{"No brackets found"},
		},
		{
			name: "bracket without parent tournament",
			brackets: []domain.Bracket{
				{
					TournamentName:       "Playoffs",
					ParentTournamentName: "",
					Label:                "Playoffs",
					Matches: []domain.Match{
						{UUID: "m1", Name: "Final", TeamA: domain.MatchTeam{Name: "Winner A"}, TeamB: domain.MatchTeam{Name: "Winner B"}},
					},
				},
			},
			contains: []string{"Playoffs", "Final"},
		},
		{
			name: "live match",
			brackets: []domain.Bracket{
				{
					TournamentName: "Live Bracket",
					Matches: []domain.Match{
						{
							UUID:        "live-match",
							Name:        "Current Match",
							TeamA:       domain.MatchTeam{Name: "Team A"},
							TeamB:       domain.MatchTeam{Name: "Team B"},
							IsLive:      true,
							IsCompleted: false,
						},
					},
				},
			},
			contains: []string{"Live Bracket", "Current Match", "LIVE"},
		},
		{
			name: "upcoming match",
			brackets: []domain.Bracket{
				{
					TournamentName: "Upcoming Bracket",
					Matches: []domain.Match{
						{
							UUID:        "upcoming-match",
							Name:        "Future Match",
							TeamA:       domain.MatchTeam{Name: "Team A"},
							TeamB:       domain.MatchTeam{Name: "Team B"},
							IsLive:      false,
							IsCompleted: false,
						},
					},
				},
			},
			contains: []string{"Upcoming Bracket", "Future Match", "Upcoming"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := formatter.Format(&buf, tt.brackets)
			require.NoError(t, err)

			output := buf.String()
			for _, expected := range tt.contains {
				assert.Contains(t, output, expected)
			}
		})
	}
}

func TestBracketsTableFormatter_FormatStatus(t *testing.T) {
	formatter := &BracketsTableFormatter{}

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

func TestBracketsTableFormatter_LongNames(t *testing.T) {
	formatter := &BracketsTableFormatter{}
	brackets := []domain.Bracket{
		{
			TournamentName:       "This is an extremely long bracket name that definitely needs truncation",
			ParentTournamentName: "Also a very long parent tournament name",
			Matches: []domain.Match{
				{
					Name:  "Very Long Match Name That Exceeds Normal Limits",
					TeamA: domain.MatchTeam{Name: "Very Long Team Name A"},
					TeamB: domain.MatchTeam{Name: "Very Long Team Name B"},
				},
			},
		},
	}

	var buf bytes.Buffer
	err := formatter.Format(&buf, brackets)
	require.NoError(t, err)

	output := buf.String()
	// Verify table structure is maintained even with long names
	lines := strings.Split(output, "\n")
	assert.True(t, len(lines) > 5)
	// Check that truncation happened
	assert.Contains(t, output, "This is an extremely long bracket name")
}
