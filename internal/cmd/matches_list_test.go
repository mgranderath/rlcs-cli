package cmd

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/mgranderath/rlcs-cli/internal/domain"
	"github.com/mgranderath/rlcs-cli/internal/output"
	"github.com/stretchr/testify/assert"
)

func TestMatchesListCmd_matchesFilters(t *testing.T) {
	tests := []struct {
		name     string
		cmd      MatchesListCmd
		match    domain.Match
		expected bool
	}{
		{
			name:     "no filters - match all",
			cmd:      MatchesListCmd{},
			match:    domain.Match{IsCompleted: true, IsLive: false},
			expected: true,
		},
		{
			name:     "completed only - match",
			cmd:      MatchesListCmd{CompletedOnly: true},
			match:    domain.Match{IsCompleted: true, IsLive: false},
			expected: true,
		},
		{
			name:     "completed only - no match",
			cmd:      MatchesListCmd{CompletedOnly: true},
			match:    domain.Match{IsCompleted: false, IsLive: true},
			expected: false,
		},
		{
			name:     "live only - match",
			cmd:      MatchesListCmd{LiveOnly: true},
			match:    domain.Match{IsLive: true, IsCompleted: false},
			expected: true,
		},
		{
			name:     "live only - no match",
			cmd:      MatchesListCmd{LiveOnly: true},
			match:    domain.Match{IsLive: false, IsCompleted: true},
			expected: false,
		},
		{
			name:     "upcoming only - match",
			cmd:      MatchesListCmd{UpcomingOnly: true},
			match:    domain.Match{IsLive: false, IsCompleted: false},
			expected: true,
		},
		{
			name:     "upcoming only - no match (completed)",
			cmd:      MatchesListCmd{UpcomingOnly: true},
			match:    domain.Match{IsLive: false, IsCompleted: true},
			expected: false,
		},
		{
			name:     "upcoming only - no match (live)",
			cmd:      MatchesListCmd{UpcomingOnly: true},
			match:    domain.Match{IsLive: true, IsCompleted: false},
			expected: false,
		},
		{
			name:     "team filter - match team A by name",
			cmd:      MatchesListCmd{Team: "Vitality"},
			match:    domain.Match{TeamA: domain.MatchTeam{Name: "Vitality"}, TeamB: domain.MatchTeam{Name: "KC"}},
			expected: true,
		},
		{
			name:     "team filter - match team B by name",
			cmd:      MatchesListCmd{Team: "KC"},
			match:    domain.Match{TeamA: domain.MatchTeam{Name: "Vitality"}, TeamB: domain.MatchTeam{Name: "KC"}},
			expected: true,
		},
		{
			name:     "team filter - case insensitive name",
			cmd:      MatchesListCmd{Team: "vitality"},
			match:    domain.Match{TeamA: domain.MatchTeam{Name: "Vitality"}, TeamB: domain.MatchTeam{Name: "KC"}},
			expected: true,
		},
		{
			name:     "team filter - match by shorthand",
			cmd:      MatchesListCmd{Team: "kc"},
			match:    domain.Match{TeamA: domain.MatchTeam{Name: "Vitality", Shorthand: "vitality"}, TeamB: domain.MatchTeam{Name: "Karmine Corp", Shorthand: "kc"}},
			expected: true,
		},
		{
			name:     "team filter - case insensitive shorthand",
			cmd:      MatchesListCmd{Team: "KC"},
			match:    domain.Match{TeamA: domain.MatchTeam{Name: "Vitality", Shorthand: "vitality"}, TeamB: domain.MatchTeam{Name: "Karmine Corp", Shorthand: "kc"}},
			expected: true,
		},
		{
			name:     "team filter - partial match on name",
			cmd:      MatchesListCmd{Team: "Vita"},
			match:    domain.Match{TeamA: domain.MatchTeam{Name: "Vitality"}, TeamB: domain.MatchTeam{Name: "KC"}},
			expected: true,
		},
		{
			name:     "team filter - partial match on shorthand",
			cmd:      MatchesListCmd{Team: "vit"},
			match:    domain.Match{TeamA: domain.MatchTeam{Name: "Vitality", Shorthand: "vitality"}, TeamB: domain.MatchTeam{Name: "KC", Shorthand: "kc"}},
			expected: true,
		},
		{
			name:     "team filter - no match",
			cmd:      MatchesListCmd{Team: "Furia"},
			match:    domain.Match{TeamA: domain.MatchTeam{Name: "Vitality"}, TeamB: domain.MatchTeam{Name: "KC"}},
			expected: false,
		},
		{
			name:     "match type filter - match",
			cmd:      MatchesListCmd{MatchType: "BO5"},
			match:    domain.Match{Type: "BO5"},
			expected: true,
		},
		{
			name:     "match type filter - case insensitive",
			cmd:      MatchesListCmd{MatchType: "bo5"},
			match:    domain.Match{Type: "BO5"},
			expected: true,
		},
		{
			name:     "match type filter - no match",
			cmd:      MatchesListCmd{MatchType: "BO7"},
			match:    domain.Match{Type: "BO5"},
			expected: false,
		},
		{
			name:     "multiple filters - all match",
			cmd:      MatchesListCmd{CompletedOnly: true, Team: "Vitality"},
			match:    domain.Match{IsCompleted: true, TeamA: domain.MatchTeam{Name: "Vitality"}},
			expected: true,
		},
		{
			name:     "multiple filters - one fails",
			cmd:      MatchesListCmd{CompletedOnly: true, Team: "Vitality"},
			match:    domain.Match{IsCompleted: false, TeamA: domain.MatchTeam{Name: "Vitality"}},
			expected: false,
		},
		{
			name:     "team name with special characters",
			cmd:      MatchesListCmd{Team: "Gen.G"},
			match:    domain.Match{TeamA: domain.MatchTeam{Name: "Gen.G Mobil1 Racing"}, TeamB: domain.MatchTeam{Name: "Other Team"}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.matchesFilters(tt.match)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMatchesListCmd_Run_HTTPMock(t *testing.T) {
	defer gock.Off()

	t.Run("successful fetch", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/games/rl/tournaments/test-tournament-id/matches").
			Reply(200).
			JSON([]map[string]interface{}{
				{
					"id":          "match-1",
					"name":        "Grand Final",
					"scheduledAt": "2026-01-15T18:00:00.000Z",
					"type":        "BO7",
					"index":       1,
					"externalId":  "external-1",
					"circuit": map[string]interface{}{
						"id":     "2026",
						"name":   "2026",
						"gameId": "rl",
					},
					"tournament": map[string]interface{}{
						"id":        "tournament-1",
						"name":      "Test Tournament",
						"startDate": "2026-01-15",
						"endDate":   "2026-01-17",
					},
					"stage": map[string]interface{}{
						"id":   "stage-1",
						"name": "Playoffs",
					},
					"teamA": map[string]interface{}{
						"id":          "team-a",
						"name":        "Team A",
						"shortName":   "teama",
						"nationality": "US",
						"externalId":  nil,
						"metadata":    nil,
					},
					"teamB": map[string]interface{}{
						"id":          "team-b",
						"name":        "Team B",
						"shortName":   "teamb",
						"nationality": "EU",
						"externalId":  nil,
						"metadata":    nil,
					},
					"teamAScore": 3,
					"teamBScore": 1,
					"maps": []map[string]interface{}{
						{
							"id":          "map-1",
							"name":        "Stadium_P",
							"scheduledAt": "2026-01-15T18:00:00.000Z",
							"startedAt":   "2026-01-15T18:05:00.000Z",
							"endedAt":     "2026-01-15T18:15:00.000Z",
							"externalId":  "MANUAL",
							"teamAScore":  3,
							"teamBScore":  1,
						},
					},
					"metadata": map[string]interface{}{
						"_t":                "rl_match",
						"teamBlueTeamId":    "team-a",
						"teamOrangeTeamId":  "team-b",
						"externalStreamUrl": "https://example.com/stream",
					},
				},
			})

		cmd := &MatchesListCmd{
			TournamentID: "test-tournament-id",
			Output:       output.MatchesFormatTable,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.NoError(t, err)
		assert.True(t, gock.IsDone())
	})

	t.Run("404 response - tournament not found", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/games/rl/tournaments/invalid-id/matches").
			Reply(404)

		cmd := &MatchesListCmd{
			TournamentID: "invalid-id",
			Output:       output.MatchesFormatTable,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tournament not found")
	})

	t.Run("500 response", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/games/rl/tournaments/test-id/matches").
			Reply(500)

		cmd := &MatchesListCmd{
			TournamentID: "test-id",
			Output:       output.MatchesFormatTable,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code: 500")
	})

	t.Run("empty matches response", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/games/rl/tournaments/empty-tournament/matches").
			Reply(200).
			JSON([]map[string]interface{}{})

		cmd := &MatchesListCmd{
			TournamentID: "empty-tournament",
			Output:       output.MatchesFormatTable,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("multiple matches", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/games/rl/tournaments/multi-match-tournament/matches").
			Reply(200).
			JSON([]map[string]interface{}{
				{
					"id":          "match-1",
					"name":        "Match 1",
					"scheduledAt": "2026-01-15T12:00:00.000Z",
					"type":        "BO5",
					"teamA":       map[string]interface{}{"id": "a", "name": "Team A"},
					"teamB":       map[string]interface{}{"id": "b", "name": "Team B"},
					"teamAScore":  3,
					"teamBScore":  1,
					"maps":        []map[string]interface{}{},
				},
				{
					"id":          "match-2",
					"name":        "Match 2",
					"scheduledAt": "2026-01-15T15:00:00.000Z",
					"type":        "BO5",
					"teamA":       map[string]interface{}{"id": "c", "name": "Team C"},
					"teamB":       map[string]interface{}{"id": "d", "name": "Team D"},
					"teamAScore":  2,
					"teamBScore":  3,
					"maps":        []map[string]interface{}{},
				},
			})

		cmd := &MatchesListCmd{
			TournamentID: "multi-match-tournament",
			Output:       output.MatchesFormatTable,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.NoError(t, err)
	})
}

func TestMatchesListCmd_Run_Validation(t *testing.T) {
	t.Run("conflicting status filters - completed and live", func(t *testing.T) {
		cmd := &MatchesListCmd{
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
		cmd := &MatchesListCmd{
			TournamentID:  "test-id",
			CompletedOnly: true,
			LiveOnly:      true,
			UpcomingOnly:  true,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.Error(t, err)
	})

	t.Run("conflicting status filters - live and upcoming", func(t *testing.T) {
		cmd := &MatchesListCmd{
			TournamentID: "test-id",
			LiveOnly:     true,
			UpcomingOnly: true,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot use multiple status filters")
	})
}

func TestMatchesListCmd_applyFilters(t *testing.T) {
	cmd := &MatchesListCmd{
		CompletedOnly: true,
		Team:          "Vitality",
	}

	matches := []domain.Match{
		{
			UUID:        "m1",
			Name:        "Match 1",
			IsCompleted: true,
			TeamA:       domain.MatchTeam{Name: "Vitality"},
			TeamB:       domain.MatchTeam{Name: "KC"},
		},
		{
			UUID:        "m2",
			Name:        "Match 2",
			IsCompleted: true,
			TeamA:       domain.MatchTeam{Name: "Furia"},
			TeamB:       domain.MatchTeam{Name: "G2"},
		},
		{
			UUID:        "m3",
			Name:        "Match 3",
			IsCompleted: false,
			TeamA:       domain.MatchTeam{Name: "Vitality"},
			TeamB:       domain.MatchTeam{Name: "KC"},
		},
	}

	filtered := cmd.applyFilters(matches)

	// Only m1 should match (completed AND has Vitality)
	assert.Len(t, filtered, 1)
	assert.Equal(t, "m1", filtered[0].UUID)
}

func TestMatchesListCmd_applyFilters_NoFilters(t *testing.T) {
	cmd := &MatchesListCmd{} // No filters set

	matches := []domain.Match{
		{UUID: "m1", Name: "Match 1"},
		{UUID: "m2", Name: "Match 2"},
		{UUID: "m3", Name: "Match 3"},
	}

	filtered := cmd.applyFilters(matches)

	// All matches should be returned when no filters are set
	assert.Len(t, filtered, 3)
}
