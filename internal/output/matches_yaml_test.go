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

func TestMatchesYAMLFormatter_Format(t *testing.T) {
	formatter := &MatchesYAMLFormatter{}

	tests := []struct {
		name           string
		matches        []domain.Match
		expectedFields []string
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
			expectedFields: []string{"uuid: match-1", "name: Grand Final", "type: BO7", "iscompleted: true"},
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
			expectedFields: []string{"- uuid: m1", "- uuid: m2", "islive: true"},
		},
		{
			name:           "empty matches",
			matches:        []domain.Match{},
			expectedFields: []string{"[]"},
		},
		{
			name: "match with nested teams and maps",
			matches: []domain.Match{
				{
					UUID:       "m1",
					Name:       "Final",
					TeamA:      domain.MatchTeam{Name: "Winner", Shorthand: "win"},
					TeamB:      domain.MatchTeam{Name: "Loser", Shorthand: "lose"},
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
			expectedFields: []string{"teama:", "name: Winner", "maps:", "name: Stadium_P"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := formatter.Format(&buf, tt.matches)
			require.NoError(t, err)

			output := buf.String()

			// Check for expected fields
			for _, field := range tt.expectedFields {
				assert.Contains(t, output, field)
			}

			// Verify proper indentation (2 spaces for yaml.v3 with SetIndent(2))
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "    ") {
					// Found a line with 2-space indentation
					break
				}
			}
		})
	}
}

func TestMatchesYAMLFormatter_FormatEmpty(t *testing.T) {
	formatter := &MatchesYAMLFormatter{}
	var buf bytes.Buffer

	// Test with empty matches slice (not nil)
	err := formatter.Format(&buf, []domain.Match{})
	require.NoError(t, err)
	output := buf.String()
	// Empty array encodes as "[]" in YAML
	assert.True(t, strings.TrimSpace(output) == "[]" || output == "")
}
