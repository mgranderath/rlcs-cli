package cmd

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/mgranderath/rlcs-cli/internal/domain"
	"github.com/mgranderath/rlcs-cli/internal/output"
	"github.com/stretchr/testify/assert"
)

func TestTournamentsBracketsCmd_matchesFilters(t *testing.T) {
	tests := []struct {
		name     string
		cmd      TournamentsBracketsCmd
		match    domain.Match
		expected bool
	}{
		{
			name:     "no filters - match all",
			cmd:      TournamentsBracketsCmd{},
			match:    domain.Match{IsCompleted: true, IsLive: false},
			expected: true,
		},
		{
			name:     "completed only - match",
			cmd:      TournamentsBracketsCmd{CompletedOnly: true},
			match:    domain.Match{IsCompleted: true, IsLive: false},
			expected: true,
		},
		{
			name:     "completed only - no match",
			cmd:      TournamentsBracketsCmd{CompletedOnly: true},
			match:    domain.Match{IsCompleted: false, IsLive: true},
			expected: false,
		},
		{
			name:     "live only - match",
			cmd:      TournamentsBracketsCmd{LiveOnly: true},
			match:    domain.Match{IsLive: true, IsCompleted: false},
			expected: true,
		},
		{
			name:     "live only - no match",
			cmd:      TournamentsBracketsCmd{LiveOnly: true},
			match:    domain.Match{IsLive: false, IsCompleted: true},
			expected: false,
		},
		{
			name:     "upcoming only - match",
			cmd:      TournamentsBracketsCmd{UpcomingOnly: true},
			match:    domain.Match{IsLive: false, IsCompleted: false},
			expected: true,
		},
		{
			name:     "upcoming only - no match (completed)",
			cmd:      TournamentsBracketsCmd{UpcomingOnly: true},
			match:    domain.Match{IsLive: false, IsCompleted: true},
			expected: false,
		},
		{
			name:     "upcoming only - no match (live)",
			cmd:      TournamentsBracketsCmd{UpcomingOnly: true},
			match:    domain.Match{IsLive: true, IsCompleted: false},
			expected: false,
		},
		{
			name:     "team filter - match team A",
			cmd:      TournamentsBracketsCmd{Team: "Vitality"},
			match:    domain.Match{TeamA: domain.MatchTeam{Name: "Vitality"}, TeamB: domain.MatchTeam{Name: "KC"}},
			expected: true,
		},
		{
			name:     "team filter - match team B",
			cmd:      TournamentsBracketsCmd{Team: "KC"},
			match:    domain.Match{TeamA: domain.MatchTeam{Name: "Vitality"}, TeamB: domain.MatchTeam{Name: "KC"}},
			expected: true,
		},
		{
			name:     "team filter - case insensitive",
			cmd:      TournamentsBracketsCmd{Team: "vitality"},
			match:    domain.Match{TeamA: domain.MatchTeam{Name: "Vitality"}, TeamB: domain.MatchTeam{Name: "KC"}},
			expected: true,
		},
		{
			name:     "team filter - partial match",
			cmd:      TournamentsBracketsCmd{Team: "Vita"},
			match:    domain.Match{TeamA: domain.MatchTeam{Name: "Vitality"}, TeamB: domain.MatchTeam{Name: "KC"}},
			expected: true,
		},
		{
			name:     "team filter - no match",
			cmd:      TournamentsBracketsCmd{Team: "Furia"},
			match:    domain.Match{TeamA: domain.MatchTeam{Name: "Vitality"}, TeamB: domain.MatchTeam{Name: "KC"}},
			expected: false,
		},
		{
			name:     "match type filter - match",
			cmd:      TournamentsBracketsCmd{MatchType: "BO5"},
			match:    domain.Match{Type: "BO5"},
			expected: true,
		},
		{
			name:     "match type filter - case insensitive",
			cmd:      TournamentsBracketsCmd{MatchType: "bo5"},
			match:    domain.Match{Type: "BO5"},
			expected: true,
		},
		{
			name:     "match type filter - no match",
			cmd:      TournamentsBracketsCmd{MatchType: "BO7"},
			match:    domain.Match{Type: "BO5"},
			expected: false,
		},
		{
			name:     "multiple filters - all match",
			cmd:      TournamentsBracketsCmd{CompletedOnly: true, Team: "Vitality"},
			match:    domain.Match{IsCompleted: true, TeamA: domain.MatchTeam{Name: "Vitality"}},
			expected: true,
		},
		{
			name:     "multiple filters - one fails",
			cmd:      TournamentsBracketsCmd{CompletedOnly: true, Team: "Vitality"},
			match:    domain.Match{IsCompleted: false, TeamA: domain.MatchTeam{Name: "Vitality"}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.matchesFilters(tt.match)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTournamentsBracketsCmd_Run_HTTPMock(t *testing.T) {
	defer gock.Off()

	t.Run("successful fetch", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/games/rl/tournaments/test-tournament-id/brackets").
			Reply(200).
			JSON([]map[string]interface{}{
				{
					"tournamentUuid":       "bracket-1",
					"tournamentName":       "Group A",
					"parentTournamentName": "Test Tournament",
					"startDate":            "2026-01-15T10:00:00.000Z",
					"endDate":              "2026-01-17T18:00:00.000Z",
					"label":                "Group A",
					"format":               "double-elim-8",
					"matches": []map[string]interface{}{
						{
							"uuid":         "match-1",
							"type":         "BO5",
							"name":         "Round 1",
							"timeOfSeries": "2026-01-15T12:00:00.000Z",
							"teamA": map[string]interface{}{
								"uuid":         "team-a",
								"name":         "Team A",
								"shorthand":    "ta",
								"location":     "US",
								"isEliminated": false,
							},
							"teamB": map[string]interface{}{
								"uuid":         "team-b",
								"name":         "Team B",
								"shorthand":    "tb",
								"location":     "EU",
								"isEliminated": false,
							},
							"teamAScore":  3,
							"teamBScore":  1,
							"maps":        []map[string]interface{}{},
							"isLive":      false,
							"isCompleted": true,
						},
					},
				},
			})

		cmd := &TournamentsBracketsCmd{
			TournamentID: "test-tournament-id",
			Output:       output.BracketsFormatTable,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.NoError(t, err)
		assert.True(t, gock.IsDone())
	})

	t.Run("404 response - tournament not found", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/games/rl/tournaments/invalid-id/brackets").
			Reply(404)

		cmd := &TournamentsBracketsCmd{
			TournamentID: "invalid-id",
			Output:       output.BracketsFormatTable,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tournament not found")
	})

	t.Run("500 response", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/games/rl/tournaments/test-id/brackets").
			Reply(500)

		cmd := &TournamentsBracketsCmd{
			TournamentID: "test-id",
			Output:       output.BracketsFormatTable,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code: 500")
	})

	t.Run("empty brackets response", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/games/rl/tournaments/empty-tournament/brackets").
			Reply(200).
			JSON([]map[string]interface{}{})

		cmd := &TournamentsBracketsCmd{
			TournamentID: "empty-tournament",
			Output:       output.BracketsFormatTable,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.NoError(t, err)
	})
}

func TestTournamentsBracketsCmd_Run_Validation(t *testing.T) {
	t.Run("conflicting status filters - completed and live", func(t *testing.T) {
		cmd := &TournamentsBracketsCmd{
			TournamentID:  "test-id",
			CompletedOnly: true,
			LiveOnly:      true,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot use multiple status filters")
	})

	t.Run("conflicting status filters - all three", func(t *testing.T) {
		cmd := &TournamentsBracketsCmd{
			TournamentID:  "test-id",
			CompletedOnly: true,
			LiveOnly:      true,
			UpcomingOnly:  true,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.Error(t, err)
	})
}
