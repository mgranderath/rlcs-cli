package output

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/mgranderath/rlcs-cli/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchesJSONFormatter_Format(t *testing.T) {
	formatter := &MatchesJSONFormatter{}

	tests := []struct {
		name          string
		matches       []domain.Match
		validateJSON  bool
		expectedField string
	}{
		{
			name: "single match",
			matches: []domain.Match{
				{
					UUID:         "match-1",
					Name:         "Grand Final",
					Type:         "BO7",
					TeamA:        domain.MatchTeam{Name: "Vitality", Shorthand: "vitality"},
					TeamB:        domain.MatchTeam{Name: "KC", Shorthand: "kc"},
					TeamAScore:   4,
					TeamBScore:   2,
					IsCompleted:  true,
					TimeOfSeries: time.Date(2026, 1, 15, 18, 0, 0, 0, time.UTC),
				},
			},
			validateJSON:  true,
			expectedField: "Grand Final",
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
					IsCompleted: false,
					IsLive:      true,
				},
			},
			validateJSON:  true,
			expectedField: "Match 1",
		},
		{
			name:          "empty matches",
			matches:       []domain.Match{},
			validateJSON:  true,
			expectedField: "[]",
		},
		{
			name: "match with maps",
			matches: []domain.Match{
				{
					UUID:       "m1",
					Name:       "Final",
					TeamA:      domain.MatchTeam{Name: "Winner"},
					TeamB:      domain.MatchTeam{Name: "Loser"},
					TeamAScore: 4,
					TeamBScore: 2,
					Maps: []domain.MatchMap{
						{
							UUID:               "map-1",
							Name:               "Stadium_P",
							ScheduledStartTime: time.Date(2026, 1, 15, 18, 0, 0, 0, time.UTC),
							TeamAScore:         3,
							TeamBScore:         2,
						},
					},
					IsCompleted: true,
				},
			},
			validateJSON:  true,
			expectedField: "Stadium_P",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := formatter.Format(&buf, tt.matches)
			require.NoError(t, err)

			output := buf.String()

			if tt.validateJSON {
				// Verify it's valid JSON
				var result []map[string]interface{}
				err = json.Unmarshal([]byte(output), &result)
				require.NoError(t, err, "output should be valid JSON")

				// Verify expected field is present if specified
				if tt.expectedField != "" {
					assert.Contains(t, output, tt.expectedField)
				}
			}

			// Verify it's formatted with indentation (except for empty array)
			assert.Contains(t, output, "\n")
			if tt.name != "empty matches" {
				assert.Contains(t, output, "  ")
			}
		})
	}
}

func TestMatchesJSONFormatter_FormatEmpty(t *testing.T) {
	formatter := &MatchesJSONFormatter{}
	var buf bytes.Buffer

	// Test with nil matches
	err := formatter.Format(&buf, nil)
	require.NoError(t, err)
	assert.Equal(t, "null\n", buf.String())
}
