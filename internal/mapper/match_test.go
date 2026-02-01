package mapper

import (
	"testing"

	"github.com/mgranderath/rlcs-cli/internal/api/blast"
	"github.com/mgranderath/rlcs-cli/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToDomainMatchesFromResponse(t *testing.T) {
	tests := []struct {
		name        string
		api         []blast.MatchResponse
		expectError bool
		expectedLen int
	}{
		{
			name: "valid matches",
			api: []blast.MatchResponse{
				{
					ID:          "match-1",
					Name:        "Grand Final",
					ScheduledAt: "2026-01-15T18:00:00.000Z",
					Type:        "BO7",
					Index:       1,
					TeamA: blast.MatchResponseTeam{
						ID:          "team-a",
						Name:        "Vitality",
						ShortName:   "vitality",
						Nationality: "FR",
					},
					TeamB: blast.MatchResponseTeam{
						ID:          "team-b",
						Name:        "Karmine Corp",
						ShortName:   "kc",
						Nationality: "FR",
					},
					TeamAScore: 4,
					TeamBScore: 2,
					Maps:       []blast.MatchResponseMap{},
				},
				{
					ID:          "match-2",
					Name:        "Semi Final",
					ScheduledAt: "2026-01-15T15:00:00.000Z",
					Type:        "BO5",
					Index:       0,
					TeamA: blast.MatchResponseTeam{
						ID:   "team-c",
						Name: "Team C",
					},
					TeamB: blast.MatchResponseTeam{
						ID:   "team-d",
						Name: "Team D",
					},
					TeamAScore: 3,
					TeamBScore: 1,
					Maps:       []blast.MatchResponseMap{},
				},
			},
			expectError: false,
			expectedLen: 2,
		},
		{
			name:        "empty matches",
			api:         []blast.MatchResponse{},
			expectError: false,
			expectedLen: 0,
		},
		{
			name: "one invalid match",
			api: []blast.MatchResponse{
				{
					ID:          "match-1",
					Name:        "Valid Match",
					ScheduledAt: "2026-01-15T18:00:00.000Z",
					Type:        "BO5",
					TeamA:       blast.MatchResponseTeam{ID: "a", Name: "A"},
					TeamB:       blast.MatchResponseTeam{ID: "b", Name: "B"},
					Maps:        []blast.MatchResponseMap{},
				},
				{
					ID:          "match-2",
					Name:        "Invalid Match",
					ScheduledAt: "invalid-time",
					Type:        "BO5",
					TeamA:       blast.MatchResponseTeam{ID: "a", Name: "A"},
					TeamB:       blast.MatchResponseTeam{ID: "b", Name: "B"},
					Maps:        []blast.MatchResponseMap{},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToDomainMatchesFromResponse(tt.api)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, result, tt.expectedLen)

			if tt.expectedLen > 0 {
				assert.Equal(t, tt.api[0].ID, result[0].UUID)
				assert.Equal(t, tt.api[0].Name, result[0].Name)
			}
		})
	}
}

func TestToDomainMatchFromResponse(t *testing.T) {
	tests := []struct {
		name        string
		api         blast.MatchResponse
		expectError bool
		checkFields func(t *testing.T, result domain.Match)
	}{
		{
			name: "completed match with maps",
			api: blast.MatchResponse{
				ID:          "match-1",
				Name:        "Grand Final",
				ScheduledAt: "2026-01-15T18:00:00.000Z",
				Type:        "BO7",
				Index:       1,
				ExternalID:  "external-1",
				Circuit: blast.Circuit{
					ID:     "2026",
					Name:   "2026",
					GameID: "rl",
				},
				Tournament: blast.Tournament{
					ID:        "tournament-1",
					Name:      "Test Tournament",
					StartDate: "2026-01-15",
					EndDate:   "2026-01-17",
				},
				Stage: blast.Stage{
					ID:   "stage-1",
					Name: "Playoffs",
				},
				TeamA: blast.MatchResponseTeam{
					ID:          "team-a",
					Name:        "Vitality",
					ShortName:   "vitality",
					Nationality: "FR",
				},
				TeamB: blast.MatchResponseTeam{
					ID:          "team-b",
					Name:        "Karmine Corp",
					ShortName:   "kc",
					Nationality: "FR",
				},
				TeamAScore: 4,
				TeamBScore: 2,
				Maps: []blast.MatchResponseMap{
					{
						ID:          "map-1",
						Name:        "Stadium_P",
						ScheduledAt: "2026-01-15T18:00:00.000Z",
						StartedAt:   "2026-01-15T18:05:00.000Z",
						EndedAt:     "2026-01-15T18:15:00.000Z",
						TeamAScore:  3,
						TeamBScore:  2,
						ExternalID:  "MANUAL",
					},
					{
						ID:          "map-2",
						Name:        "Urban_P",
						ScheduledAt: "2026-01-15T18:20:00.000Z",
						StartedAt:   "2026-01-15T18:25:00.000Z",
						EndedAt:     "2026-01-15T18:35:00.000Z",
						TeamAScore:  4,
						TeamBScore:  3,
						ExternalID:  "MANUAL",
					},
				},
				Metadata: blast.MatchResponseMetadata{
					T:                 "rl_match",
					TeamBlueTeamID:    "team-a",
					TeamOrangeTeamID:  "team-b",
					ExternalStreamURL: "https://example.com/stream",
				},
			},
			expectError: false,
			checkFields: func(t *testing.T, result domain.Match) {
				assert.Equal(t, "match-1", result.UUID)
				assert.Equal(t, "Grand Final", result.Name)
				assert.Equal(t, "BO7", result.Type)
				assert.Equal(t, 1, result.Index)
				assert.Equal(t, "external-1", result.ExternalID)
				assert.Equal(t, "Vitality", result.TeamA.Name)
				assert.Equal(t, "vitality", result.TeamA.Shorthand)
				assert.Equal(t, "FR", result.TeamA.Location)
				assert.Equal(t, "Karmine Corp", result.TeamB.Name)
				assert.Equal(t, "kc", result.TeamB.Shorthand)
				assert.Equal(t, "FR", result.TeamB.Location)
				assert.Equal(t, 4, result.TeamAScore)
				assert.Equal(t, 2, result.TeamBScore)
				assert.True(t, result.IsCompleted)
				assert.False(t, result.IsLive)
				assert.Len(t, result.Maps, 2)
				assert.Nil(t, result.WinnerGoesTo)
				assert.Nil(t, result.LoserGoesTo)
			},
		},
		{
			name: "live match - some maps started but not all ended",
			api: blast.MatchResponse{
				ID:          "match-2",
				Name:        "Current Match",
				ScheduledAt: "2026-01-15T15:00:00.000Z",
				Type:        "BO5",
				TeamA: blast.MatchResponseTeam{
					ID:   "team-a",
					Name: "Team A",
				},
				TeamB: blast.MatchResponseTeam{
					ID:   "team-b",
					Name: "Team B",
				},
				TeamAScore: 1,
				TeamBScore: 0,
				Maps: []blast.MatchResponseMap{
					{
						ID:          "map-1",
						Name:        "Stadium_P",
						ScheduledAt: "2026-01-15T15:00:00.000Z",
						StartedAt:   "2026-01-15T15:05:00.000Z",
						EndedAt:     "2026-01-15T15:15:00.000Z",
						TeamAScore:  3,
						TeamBScore:  2,
					},
					{
						ID:          "map-2",
						Name:        "Urban_P",
						ScheduledAt: "2026-01-15T15:20:00.000Z",
						StartedAt:   "2026-01-15T15:25:00.000Z",
						EndedAt:     "",
						TeamAScore:  1,
						TeamBScore:  0,
					},
				},
			},
			expectError: false,
			checkFields: func(t *testing.T, result domain.Match) {
				assert.True(t, result.IsLive)
				assert.False(t, result.IsCompleted)
			},
		},
		{
			name: "upcoming match - no maps started",
			api: blast.MatchResponse{
				ID:          "match-3",
				Name:        "Future Match",
				ScheduledAt: "2026-01-16T12:00:00.000Z",
				Type:        "BO5",
				TeamA: blast.MatchResponseTeam{
					ID:   "team-a",
					Name: "Team A",
				},
				TeamB: blast.MatchResponseTeam{
					ID:   "team-b",
					Name: "Team B",
				},
				TeamAScore: 0,
				TeamBScore: 0,
				Maps: []blast.MatchResponseMap{
					{
						ID:          "map-1",
						Name:        "Stadium_P",
						ScheduledAt: "2026-01-16T12:00:00.000Z",
						StartedAt:   "",
						EndedAt:     "",
						TeamAScore:  0,
						TeamBScore:  0,
					},
				},
			},
			expectError: false,
			checkFields: func(t *testing.T, result domain.Match) {
				assert.False(t, result.IsLive)
				assert.False(t, result.IsCompleted)
			},
		},
		{
			name: "match with no maps",
			api: blast.MatchResponse{
				ID:          "match-4",
				Name:        "Empty Match",
				ScheduledAt: "2026-01-16T12:00:00.000Z",
				Type:        "BO5",
				TeamA:       blast.MatchResponseTeam{ID: "a", Name: "A"},
				TeamB:       blast.MatchResponseTeam{ID: "b", Name: "B"},
				Maps:        []blast.MatchResponseMap{},
			},
			expectError: false,
			checkFields: func(t *testing.T, result domain.Match) {
				assert.False(t, result.IsLive)
				assert.False(t, result.IsCompleted)
				assert.Empty(t, result.Maps)
			},
		},
		{
			name: "invalid scheduled time",
			api: blast.MatchResponse{
				ID:          "match-5",
				Name:        "Invalid Match",
				ScheduledAt: "invalid-time",
				Type:        "BO5",
				TeamA:       blast.MatchResponseTeam{ID: "a", Name: "A"},
				TeamB:       blast.MatchResponseTeam{ID: "b", Name: "B"},
				Maps:        []blast.MatchResponseMap{},
			},
			expectError: true,
		},
		{
			name: "external ID is null",
			api: blast.MatchResponse{
				ID:          "match-6",
				Name:        "Match with null external ID",
				ScheduledAt: "2026-01-15T18:00:00.000Z",
				Type:        "BO5",
				ExternalID:  "",
				TeamA: blast.MatchResponseTeam{
					ID:         "team-a",
					Name:       "Team A",
					ExternalID: nil,
				},
				TeamB: blast.MatchResponseTeam{
					ID:   "team-b",
					Name: "Team B",
				},
				Maps: []blast.MatchResponseMap{},
			},
			expectError: false,
			checkFields: func(t *testing.T, result domain.Match) {
				assert.Equal(t, "", result.ExternalID)
				assert.Equal(t, "Team A", result.TeamA.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toDomainMatchFromResponse(tt.api)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.checkFields != nil {
				tt.checkFields(t, result)
			}
		})
	}
}

func TestToDomainMatchMapFromResponse(t *testing.T) {
	tests := []struct {
		name        string
		api         blast.MatchResponseMap
		expectError bool
	}{
		{
			name: "valid map with all times",
			api: blast.MatchResponseMap{
				ID:          "map-1",
				Name:        "Stadium_P",
				ScheduledAt: "2026-01-15T12:00:00.000Z",
				StartedAt:   "2026-01-15T12:05:00.000Z",
				EndedAt:     "2026-01-15T12:15:00.000Z",
				TeamAScore:  3,
				TeamBScore:  1,
				ExternalID:  "MANUAL",
			},
			expectError: false,
		},
		{
			name: "map with empty started and ended times",
			api: blast.MatchResponseMap{
				ID:          "map-2",
				Name:        "Urban_P",
				ScheduledAt: "2026-01-15T12:00:00.000Z",
				StartedAt:   "",
				EndedAt:     "",
				TeamAScore:  0,
				TeamBScore:  0,
				ExternalID:  "MANUAL",
			},
			expectError: false,
		},
		{
			name: "invalid scheduled time",
			api: blast.MatchResponseMap{
				ID:          "map-3",
				Name:        "Invalid",
				ScheduledAt: "invalid-time",
				StartedAt:   "",
				EndedAt:     "",
			},
			expectError: true,
		},
		{
			name: "invalid started time",
			api: blast.MatchResponseMap{
				ID:          "map-4",
				Name:        "Invalid",
				ScheduledAt: "2026-01-15T12:00:00.000Z",
				StartedAt:   "invalid",
				EndedAt:     "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toDomainMatchMapFromResponse(tt.api)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.api.ID, result.UUID)
			assert.Equal(t, tt.api.Name, result.Name)
			assert.Equal(t, tt.api.TeamAScore, result.TeamAScore)
			assert.Equal(t, tt.api.TeamBScore, result.TeamBScore)
		})
	}
}

func TestInferMatchStatus(t *testing.T) {
	tests := []struct {
		name        string
		maps        []blast.MatchResponseMap
		isCompleted bool
		isLive      bool
	}{
		{
			name:        "empty maps",
			maps:        []blast.MatchResponseMap{},
			isCompleted: false,
			isLive:      false,
		},
		{
			name: "completed - all maps ended",
			maps: []blast.MatchResponseMap{
				{StartedAt: "2026-01-15T12:00:00.000Z", EndedAt: "2026-01-15T12:10:00.000Z"},
				{StartedAt: "2026-01-15T12:15:00.000Z", EndedAt: "2026-01-15T12:25:00.000Z"},
			},
			isCompleted: true,
			isLive:      false,
		},
		{
			name: "live - some maps started, not all ended",
			maps: []blast.MatchResponseMap{
				{StartedAt: "2026-01-15T12:00:00.000Z", EndedAt: "2026-01-15T12:10:00.000Z"},
				{StartedAt: "2026-01-15T12:15:00.000Z", EndedAt: ""},
			},
			isCompleted: false,
			isLive:      true,
		},
		{
			name: "upcoming - no maps started",
			maps: []blast.MatchResponseMap{
				{StartedAt: "", EndedAt: ""},
				{StartedAt: "", EndedAt: ""},
			},
			isCompleted: false,
			isLive:      false,
		},
		{
			name: "live - first map started but not ended",
			maps: []blast.MatchResponseMap{
				{StartedAt: "2026-01-15T12:00:00.000Z", EndedAt: ""},
			},
			isCompleted: false,
			isLive:      true,
		},
		{
			name: "completed - single map started and ended",
			maps: []blast.MatchResponseMap{
				{StartedAt: "2026-01-15T12:00:00.000Z", EndedAt: "2026-01-15T12:10:00.000Z"},
			},
			isCompleted: true,
			isLive:      false,
		},
		{
			name: "live - multiple maps with last not ended",
			maps: []blast.MatchResponseMap{
				{StartedAt: "2026-01-15T12:00:00.000Z", EndedAt: "2026-01-15T12:10:00.000Z"},
				{StartedAt: "2026-01-15T12:15:00.000Z", EndedAt: "2026-01-15T12:25:00.000Z"},
				{StartedAt: "2026-01-15T12:30:00.000Z", EndedAt: ""},
			},
			isCompleted: false,
			isLive:      true,
		},
		{
			name: "upcoming - mix of empty maps and one not started map",
			maps: []blast.MatchResponseMap{
				{StartedAt: "", EndedAt: ""},
				{StartedAt: "", EndedAt: ""},
				{StartedAt: "", EndedAt: ""},
			},
			isCompleted: false,
			isLive:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isCompleted, isLive := inferMatchStatus(tt.maps)
			assert.Equal(t, tt.isCompleted, isCompleted, "isCompleted mismatch")
			assert.Equal(t, tt.isLive, isLive, "isLive mismatch")
		})
	}
}
