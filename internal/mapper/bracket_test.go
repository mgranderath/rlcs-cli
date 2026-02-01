package mapper

import (
	"testing"

	"github.com/mgranderath/rlcs-cli/internal/api/blast"
	"github.com/mgranderath/rlcs-cli/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToDomainBracket(t *testing.T) {
	tests := []struct {
		name        string
		api         blast.Bracket
		expectError bool
	}{
		{
			name: "valid bracket with matches",
			api: blast.Bracket{
				TournamentUUID:       "test-uuid",
				TournamentName:       "Group A",
				ParentTournamentName: "RLCS Open 2026",
				StartDate:            "2026-01-15T10:00:00.000Z",
				EndDate:              "2026-01-17T18:00:00.000Z",
				Label:                "Group A",
				Format:               "double-elim-8",
				Matches: []blast.Match{
					{
						UUID:         "match-1",
						Type:         "BO5",
						Name:         "Round 1",
						TimeOfSeries: "2026-01-15T12:00:00.000Z",
						TeamA:        blast.Team{UUID: "team-a", Name: "Team A"},
						TeamB:        blast.Team{UUID: "team-b", Name: "Team B"},
						TeamAScore:   3,
						TeamBScore:   1,
						IsCompleted:  true,
					},
				},
			},
			expectError: false,
		},
		{
			name: "bracket with nil number of teams",
			api: blast.Bracket{
				TournamentUUID: "test-uuid",
				TournamentName: "Playoffs",
				StartDate:      "2026-01-15T10:00:00.000Z",
				EndDate:        "2026-01-17T18:00:00.000Z",
				NumberOfTeams:  nil,
				Matches:        []blast.Match{},
			},
			expectError: false,
		},
		{
			name: "invalid start date",
			api: blast.Bracket{
				TournamentUUID: "test-uuid",
				StartDate:      "invalid-date",
				EndDate:        "2026-01-17T18:00:00.000Z",
				Matches:        []blast.Match{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToDomainBracket(tt.api)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.api.TournamentUUID, result.TournamentUUID)
			assert.Equal(t, tt.api.TournamentName, result.TournamentName)
		})
	}
}

func TestToDomainMatch(t *testing.T) {
	tests := []struct {
		name        string
		api         blast.Match
		expectError bool
		checkFields func(t *testing.T, result domain.Match)
	}{
		{
			name: "completed match with winner/loser destinations",
			api: blast.Match{
				UUID:         "match-1",
				Type:         "BO7",
				Name:         "Grand Final",
				TimeOfSeries: "2026-01-15T18:00:00.000Z",
				TeamA: blast.Team{
					UUID:         "team-a",
					Name:         "Vitality",
					Shorthand:    "vitality",
					Location:     "FR",
					IsEliminated: false,
				},
				TeamB: blast.Team{
					UUID:         "team-b",
					Name:         "Karmine Corp",
					Shorthand:    "kc",
					Location:     "FR",
					IsEliminated: true,
				},
				TeamAScore: 4,
				TeamBScore: 2,
				Maps: []blast.Map{
					{
						UUID:               "map-1",
						ScheduledStartTime: "2026-01-15T18:00:00.000Z",
						ActualStartTime:    "2026-01-15T18:05:00.000Z",
						Name:               "Stadium_P",
						MatchEndedTime:     "2026-01-15T18:10:00.000Z",
						TeamAScore:         3,
						TeamBScore:         2,
						ExternalID:         "MANUAL",
					},
				},
				WinnerGoesTo: &blast.BracketDestination{
					TournamentUUID:  "next-tournament",
					SeriesUUID:      "next-match",
					BracketPosition: "POSITION_A",
				},
				LoserGoesTo: nil,
				IsLive:      false,
				IsCompleted: true,
			},
			expectError: false,
			checkFields: func(t *testing.T, result domain.Match) {
				assert.Equal(t, "match-1", result.UUID)
				assert.Equal(t, "BO7", result.Type)
				assert.Equal(t, "Vitality", result.TeamA.Name)
				assert.Equal(t, "Karmine Corp", result.TeamB.Name)
				assert.Equal(t, 4, result.TeamAScore)
				assert.Equal(t, 2, result.TeamBScore)
				assert.True(t, result.IsCompleted)
				assert.False(t, result.IsLive)
				require.NotNil(t, result.WinnerGoesTo)
				assert.Equal(t, "next-tournament", result.WinnerGoesTo.TournamentUUID)
				assert.Nil(t, result.LoserGoesTo)
			},
		},
		{
			name: "live match",
			api: blast.Match{
				UUID:         "match-2",
				Type:         "BO5",
				Name:         "Semi Final",
				TimeOfSeries: "2026-01-15T15:00:00.000Z",
				TeamA:        blast.Team{Name: "Team A"},
				TeamB:        blast.Team{Name: "Team B"},
				TeamAScore:   2,
				TeamBScore:   1,
				IsLive:       true,
				IsCompleted:  false,
				Maps:         []blast.Map{},
			},
			expectError: false,
			checkFields: func(t *testing.T, result domain.Match) {
				assert.True(t, result.IsLive)
				assert.False(t, result.IsCompleted)
			},
		},
		{
			name: "invalid time format",
			api: blast.Match{
				UUID:         "match-3",
				TimeOfSeries: "invalid-time",
				TeamA:        blast.Team{Name: "Team A"},
				TeamB:        blast.Team{Name: "Team B"},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToDomainMatch(tt.api)

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

func TestToDomainMap(t *testing.T) {
	tests := []struct {
		name        string
		api         blast.Map
		expectError bool
	}{
		{
			name: "valid map",
			api: blast.Map{
				UUID:               "map-1",
				ScheduledStartTime: "2026-01-15T12:00:00.000Z",
				ActualStartTime:    "2026-01-15T12:05:00.000Z",
				Name:               "Stadium_P",
				MatchEndedTime:     "2026-01-15T12:15:00.000Z",
				TeamAScore:         3,
				TeamBScore:         1,
				ExternalID:         "MANUAL",
			},
			expectError: false,
		},
		{
			name: "invalid scheduled start time",
			api: blast.Map{
				UUID:               "map-1",
				ScheduledStartTime: "invalid",
				ActualStartTime:    "2026-01-15T12:05:00.000Z",
				MatchEndedTime:     "2026-01-15T12:15:00.000Z",
			},
			expectError: true,
		},
		{
			name: "map with empty actual and ended times",
			api: blast.Map{
				UUID:               "map-2",
				ScheduledStartTime: "2026-01-15T12:00:00.000Z",
				ActualStartTime:    "",
				MatchEndedTime:     "",
				Name:               "Stadium_P",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToDomainMap(tt.api)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.api.UUID, result.UUID)
			assert.Equal(t, tt.api.Name, result.Name)
		})
	}
}

func TestToDomainBrackets(t *testing.T) {
	tests := []struct {
		name        string
		api         []blast.Bracket
		expectError bool
		expectedLen int
	}{
		{
			name: "multiple brackets",
			api: []blast.Bracket{
				{
					TournamentUUID: "bracket-1",
					TournamentName: "Group A",
					StartDate:      "2026-01-15T10:00:00.000Z",
					EndDate:        "2026-01-17T18:00:00.000Z",
					Matches:        []blast.Match{},
				},
				{
					TournamentUUID: "bracket-2",
					TournamentName: "Group B",
					StartDate:      "2026-01-15T10:00:00.000Z",
					EndDate:        "2026-01-17T18:00:00.000Z",
					Matches:        []blast.Match{},
				},
			},
			expectError: false,
			expectedLen: 2,
		},
		{
			name:        "empty brackets",
			api:         []blast.Bracket{},
			expectError: false,
			expectedLen: 0,
		},
		{
			name: "one invalid bracket",
			api: []blast.Bracket{
				{
					TournamentUUID: "bracket-1",
					StartDate:      "2026-01-15T10:00:00.000Z",
					EndDate:        "2026-01-17T18:00:00.000Z",
					Matches:        []blast.Match{},
				},
				{
					TournamentUUID: "bracket-2",
					StartDate:      "invalid-date",
					EndDate:        "2026-01-17T18:00:00.000Z",
					Matches:        []blast.Match{},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToDomainBrackets(tt.api)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, result, tt.expectedLen)
		})
	}
}
