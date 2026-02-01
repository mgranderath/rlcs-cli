package mapper

import (
	"testing"

	"github.com/mgranderath/rlcs-cli/internal/api/blast"
	"github.com/mgranderath/rlcs-cli/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToDomainTournament(t *testing.T) {
	tests := []struct {
		name        string
		api         blast.Tournament
		expected    domain.Tournament
		expectError bool
	}{
		{
			name: "valid tournament - EU Open",
			api: blast.Tournament{
				ID:            "rlcs-open-1-eu-2026",
				Name:          "RLCS Open 1 EU 2026",
				StartDate:     "2026-01-15",
				EndDate:       "2026-01-17",
				CircuitID:     "2026",
				PrizePool:     "$50,000",
				Location:      "Online",
				NumberOfTeams: 16,
				Region:        "EU",
				Grouping:      "RLCS Open 1 2026",
				Description:   "EU Open Tournament",
			},
			expected: domain.Tournament{
				ID:          "rlcs-open-1-eu-2026",
				Name:        "RLCS Open 1 EU 2026",
				CircuitID:   "2026",
				PrizePool:   "$50,000",
				Location:    "Online",
				TeamCount:   16,
				Region:      domain.RegionEU,
				Type:        domain.TypeOpen,
				Description: "EU Open Tournament",
				IsOnline:    true,
				IsMajor:     false,
			},
			expectError: false,
		},
		{
			name: "valid tournament - Major (no region/grouping)",
			api: blast.Tournament{
				ID:            "major-1-2026",
				Name:          "RLCS Major 1 2026",
				StartDate:     "2026-03-15",
				EndDate:       "2026-03-20",
				CircuitID:     "2026",
				PrizePool:     "$300,000",
				Location:      "Berlin",
				NumberOfTeams: 16,
				Region:        "",
				Grouping:      "",
				Description:   "First Major of 2026",
			},
			expected: domain.Tournament{
				ID:          "major-1-2026",
				Name:        "RLCS Major 1 2026",
				CircuitID:   "2026",
				PrizePool:   "$300,000",
				Location:    "Berlin",
				TeamCount:   16,
				Region:      domain.RegionNone,
				Type:        domain.TypeMajor,
				Description: "First Major of 2026",
				IsOnline:    false,
				IsMajor:     true,
			},
			expectError: false,
		},
		{
			name: "valid tournament - World Championship",
			api: blast.Tournament{
				ID:            "world-championship-2026",
				Name:          "RLCS World Championship 2026",
				StartDate:     "2026-08-01",
				EndDate:       "2026-08-10",
				CircuitID:     "2026",
				PrizePool:     "$1,000,000",
				Location:      "Los Angeles",
				NumberOfTeams: 24,
				Region:        "",
				Grouping:      "",
				Description:   "World Championship",
			},
			expected: domain.Tournament{
				ID:          "world-championship-2026",
				Name:        "RLCS World Championship 2026",
				CircuitID:   "2026",
				PrizePool:   "$1,000,000",
				Location:    "Los Angeles",
				TeamCount:   24,
				Region:      domain.RegionNone,
				Type:        domain.TypeWorldChampionship,
				Description: "World Championship",
				IsOnline:    false,
				IsMajor:     true,
			},
			expectError: false,
		},
		{
			name: "invalid date format",
			api: blast.Tournament{
				ID:        "invalid",
				StartDate: "invalid-date",
				EndDate:   "2026-01-17",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToDomainTournament(tt.api)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Region, result.Region)
			assert.Equal(t, tt.expected.Type, result.Type)
			assert.Equal(t, tt.expected.IsOnline, result.IsOnline)
			assert.Equal(t, tt.expected.IsMajor, result.IsMajor)
		})
	}
}

func TestToDomainTournaments(t *testing.T) {
	tests := []struct {
		name        string
		api         []blast.Tournament
		expectError bool
		expectedLen int
	}{
		{
			name: "multiple tournaments",
			api: []blast.Tournament{
				{ID: "1", Name: "T1", StartDate: "2026-01-01", EndDate: "2026-01-02", Region: "NA"},
				{ID: "2", Name: "T2", StartDate: "2026-01-03", EndDate: "2026-01-04", Region: "EU"},
			},
			expectError: false,
			expectedLen: 2,
		},
		{
			name:        "empty slice",
			api:         []blast.Tournament{},
			expectError: false,
			expectedLen: 0,
		},
		{
			name: "one invalid tournament",
			api: []blast.Tournament{
				{ID: "1", Name: "T1", StartDate: "2026-01-01", EndDate: "2026-01-02"},
				{ID: "2", Name: "T2", StartDate: "invalid", EndDate: "2026-01-04"},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToDomainTournaments(tt.api)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, result, tt.expectedLen)
		})
	}
}

func TestParseRegion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected domain.Region
	}{
		{"NA uppercase", "NA", domain.RegionNA},
		{"na lowercase", "na", domain.RegionNA},
		{"EU uppercase", "EU", domain.RegionEU},
		{"eu lowercase", "eu", domain.RegionEU},
		{"APAC", "APAC", domain.RegionAPAC},
		{"apac", "apac", domain.RegionAPAC},
		{"SAM", "SAM", domain.RegionSAM},
		{"OCE", "OCE", domain.RegionOCE},
		{"MENA", "MENA", domain.RegionMENA},
		{"SSA", "SSA", domain.RegionSSA},
		{"empty string", "", domain.RegionNone},
		{"unknown region", "UNKNOWN", domain.RegionNone},
		{"mixed case", "Na", domain.RegionNA},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseRegion(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetermineTournamentType(t *testing.T) {
	tests := []struct {
		name     string
		tourName string
		grouping string
		region   string
		expected domain.TournamentType
	}{
		{"World Championship", "RLCS World Championship 2026", "", "", domain.TypeWorldChampionship},
		{"Kickoff tournament", "RLCS Kick-Off Tournament 2026", "", "EU", domain.TypeKickoff},
		{"Kickoff variant", "RLCS Kickoff 2026", "", "NA", domain.TypeKickoff},
		{"Major (no region/grouping)", "RLCS Major 1 2026", "", "", domain.TypeMajor},
		{"Open tournament", "RLCS Open 1 EU 2026", "RLCS Open 1 2026", "EU", domain.TypeOpen},
		{"Default to Open", "Random Tournament", "group", "NA", domain.TypeOpen},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineTournamentType(tt.tourName, tt.grouping, tt.region)
			assert.Equal(t, tt.expected, result)
		})
	}
}
