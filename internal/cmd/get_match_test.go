package cmd

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/mgranderath/rlcs-cli/internal/output"
	"github.com/stretchr/testify/assert"
)

func TestGetMatchCmd_Run_HTTPMock(t *testing.T) {
	defer gock.Off()

	t.Run("successful fetch", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/matches/test-match-id/detailed").
			Reply(200).
			JSON(map[string]interface{}{
				"id":          "test-match-id",
				"name":        "Quarter Final 1",
				"scheduledAt": "2026-02-01T10:00:00.000Z",
				"type":        "BO7",
				"index":       4,
				"externalId":  nil,
				"circuit": map[string]interface{}{
					"gameId": "rl",
					"id":     "2026",
					"name":   "2026",
				},
				"tournament": map[string]interface{}{
					"id":         "rlcs-open-3-apac-2026",
					"name":       "RLCS Open 3 APAC 2026",
					"startDate":  "2026-01-30",
					"endDate":    "2026-02-01",
					"prizePool":  "$29,700",
					"externalId": "regional-3-21op85bf52",
				},
				"stage": map[string]interface{}{
					"id":            "bdc700c8-cddc-44fd-8f59-faa3d4fb696a",
					"name":          "Playoffs",
					"format":        "afl-final-eight",
					"numberOfTeams": nil,
					"metadata":      nil,
					"startDate":     "2026-01-31T16:00:00",
					"endDate":       "2026-02-01T20:00:00",
					"index":         2,
				},
				"teamA": map[string]interface{}{
					"id":          "49737583-e29b-4056-a706-41ac13db0d39",
					"name":        "God Speed",
					"shortName":   "godspeed",
					"nationality": "MY",
					"externalId":  nil,
					"metadata":    nil,
				},
				"teamB": map[string]interface{}{
					"id":          "20607998-f78d-4a56-b4ba-54078cc29752",
					"name":        "Ground Zero Gaming",
					"shortName":   "gzg",
					"nationality": "AU",
					"externalId":  nil,
					"metadata":    nil,
				},
				"teamAScore": 1,
				"teamBScore": 3,
				"maps":       []map[string]interface{}{},
				"metadata":   nil,
			})

		cmd := &GetMatchCmd{
			MatchID: "test-match-id",
			Output:  output.MatchesFormatTable,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.NoError(t, err)
		assert.True(t, gock.IsDone())
	})

	t.Run("404 response - match not found", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/matches/invalid-id/detailed").
			Reply(404)

		cmd := &GetMatchCmd{
			MatchID: "invalid-id",
			Output:  output.MatchesFormatTable,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "match not found")
	})

	t.Run("500 response", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/matches/test-id/detailed").
			Reply(500)

		cmd := &GetMatchCmd{
			MatchID: "test-id",
			Output:  output.MatchesFormatTable,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code: 500")
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/matches/test-id/detailed").
			Reply(200).
			BodyString("invalid json")

		cmd := &GetMatchCmd{
			MatchID: "test-id",
			Output:  output.MatchesFormatTable,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse JSON")
	})

	t.Run("match with maps", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/matches/match-with-maps/detailed").
			Reply(200).
			JSON(map[string]interface{}{
				"id":          "match-with-maps",
				"name":        "Grand Final",
				"scheduledAt": "2026-01-15T18:00:00.000Z",
				"type":        "BO7",
				"teamA": map[string]interface{}{
					"id":          "team-a",
					"name":        "Team A",
					"shortName":   "teama",
					"nationality": "US",
				},
				"teamB": map[string]interface{}{
					"id":          "team-b",
					"name":        "Team B",
					"shortName":   "teamb",
					"nationality": "EU",
				},
				"teamAScore": 4,
				"teamBScore": 2,
				"maps": []map[string]interface{}{
					{
						"id":          "map-1",
						"name":        "Stadium_P",
						"scheduledAt": "2026-01-15T18:00:00.000Z",
						"startedAt":   "2026-01-15T18:05:00.000Z",
						"endedAt":     "2026-01-15T18:15:00.000Z",
						"teamAScore":  3,
						"teamBScore":  2,
						"externalId":  "MANUAL",
					},
					{
						"id":          "map-2",
						"name":        "Urban_P",
						"scheduledAt": "2026-01-15T18:20:00.000Z",
						"startedAt":   "2026-01-15T18:25:00.000Z",
						"endedAt":     "2026-01-15T18:35:00.000Z",
						"teamAScore":  4,
						"teamBScore":  3,
						"externalId":  "MANUAL",
					},
				},
			})

		cmd := &GetMatchCmd{
			MatchID: "match-with-maps",
			Output:  output.MatchesFormatTable,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.NoError(t, err)
		assert.True(t, gock.IsDone())
	})

	t.Run("JSON output format", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/matches/json-test/detailed").
			Reply(200).
			JSON(map[string]interface{}{
				"id":          "json-test",
				"name":        "Test Match",
				"scheduledAt": "2026-01-15T18:00:00.000Z",
				"type":        "BO5",
				"teamA":       map[string]interface{}{"id": "a", "name": "Team A"},
				"teamB":       map[string]interface{}{"id": "b", "name": "Team B"},
				"teamAScore":  3,
				"teamBScore":  1,
				"maps":        []map[string]interface{}{},
			})

		cmd := &GetMatchCmd{
			MatchID: "json-test",
			Output:  output.MatchesFormatJSON,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("YAML output format", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/matches/yaml-test/detailed").
			Reply(200).
			JSON(map[string]interface{}{
				"id":          "yaml-test",
				"name":        "Test Match",
				"scheduledAt": "2026-01-15T18:00:00.000Z",
				"type":        "BO5",
				"teamA":       map[string]interface{}{"id": "a", "name": "Team A"},
				"teamB":       map[string]interface{}{"id": "b", "name": "Team B"},
				"teamAScore":  3,
				"teamBScore":  1,
				"maps":        []map[string]interface{}{},
			})

		cmd := &GetMatchCmd{
			MatchID: "yaml-test",
			Output:  output.MatchesFormatYAML,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("invalid time format in response", func(t *testing.T) {
		gock.New("https://api.blast.tv").
			Get("/v2/matches/invalid-time/detailed").
			Reply(200).
			JSON(map[string]interface{}{
				"id":          "invalid-time",
				"name":        "Test Match",
				"scheduledAt": "invalid-time-format",
				"type":        "BO5",
				"teamA":       map[string]interface{}{"id": "a", "name": "Team A"},
				"teamB":       map[string]interface{}{"id": "b", "name": "Team B"},
				"maps":        []map[string]interface{}{},
			})

		cmd := &GetMatchCmd{
			MatchID: "invalid-time",
			Output:  output.MatchesFormatTable,
		}
		ctx := &Context{Debug: false}

		err := cmd.Run(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to map match")
	})
}
