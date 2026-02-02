package cmd

import (
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/mgranderath/rlcs-cli/internal/domain"
	"github.com/mgranderath/rlcs-cli/internal/output"
	"github.com/stretchr/testify/assert"
)

func TestTournamentsMatchesCmd_matchesStatusFilter(t *testing.T) {
	cmd := &TournamentsMatchesCmd{}

	live := domain.Match{IsLive: true, IsCompleted: false}
	upcoming := domain.Match{IsLive: false, IsCompleted: false}
	completed := domain.Match{IsLive: false, IsCompleted: true}

	assert.True(t, cmd.matchesStatusFilter(live))
	assert.True(t, cmd.matchesStatusFilter(upcoming))
	assert.False(t, cmd.matchesStatusFilter(completed))

	cmd = &TournamentsMatchesCmd{LiveOnly: true}
	assert.True(t, cmd.matchesStatusFilter(live))
	assert.False(t, cmd.matchesStatusFilter(upcoming))
	assert.False(t, cmd.matchesStatusFilter(completed))

	cmd = &TournamentsMatchesCmd{UpcomingOnly: true}
	assert.False(t, cmd.matchesStatusFilter(live))
	assert.True(t, cmd.matchesStatusFilter(upcoming))
	assert.False(t, cmd.matchesStatusFilter(completed))

	cmd = &TournamentsMatchesCmd{CompletedOnly: true}
	assert.False(t, cmd.matchesStatusFilter(live))
	assert.False(t, cmd.matchesStatusFilter(upcoming))
	assert.True(t, cmd.matchesStatusFilter(completed))
}

func TestTournamentsMatchesCmd_matchesTournamentFilters(t *testing.T) {
	cmd := &TournamentsMatchesCmd{
		Region:   "EU",
		Online:   true,
		Major:    true,
		Grouping: "Open 1",
		MinTeams: 16,
	}

	tournament := domain.Tournament{
		Name:      "RLCS Open 1 2026",
		Region:    domain.RegionEU,
		IsOnline:  true,
		IsMajor:   true,
		TeamCount: 16,
	}

	assert.True(t, cmd.matchesTournamentFilters(tournament))

	tournament.Region = domain.RegionNA
	assert.False(t, cmd.matchesTournamentFilters(tournament))
}

func TestSortGames(t *testing.T) {
	timeA := time.Date(2026, 1, 10, 10, 0, 0, 0, time.UTC)
	timeB := time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC)

	games := []domain.GameListing{
		{TournamentName: "T2", Match: domain.Match{Name: "Completed", IsCompleted: true, TimeOfSeries: timeA}},
		{TournamentName: "T1", Match: domain.Match{Name: "Upcoming B", TimeOfSeries: timeB}},
		{TournamentName: "T1", Match: domain.Match{Name: "Live", IsLive: true, TimeOfSeries: timeB}},
		{TournamentName: "T1", Match: domain.Match{Name: "Upcoming A", TimeOfSeries: timeA}},
	}

	sortGames(games)

	assert.Equal(t, "Live", games[0].Match.Name)
	assert.Equal(t, "Upcoming A", games[1].Match.Name)
	assert.Equal(t, "Upcoming B", games[2].Match.Name)
	assert.Equal(t, "Completed", games[3].Match.Name)
}

func TestTournamentsMatchesCmd_Run_HTTPMock(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.blast.tv").
		Get("/v2/circuits/2026/tournaments").
		MatchParam("game", "rl").
		Reply(200).
		JSON([]map[string]interface{}{
			{
				"id":            "tournament-1",
				"name":          "Tournament One",
				"startDate":     "2026-01-10",
				"endDate":       "2026-01-12",
				"circuitId":     "2026",
				"region":        "EU",
				"numberOfTeams": 16,
				"location":      "Online",
				"grouping":      "",
			},
			{
				"id":            "tournament-2",
				"name":          "Tournament Two",
				"startDate":     "2026-01-15",
				"endDate":       "2026-01-17",
				"circuitId":     "2026",
				"region":        "NA",
				"numberOfTeams": 16,
				"location":      "Online",
				"grouping":      "",
			},
		})

	gock.New("https://api.blast.tv").
		Get("/v2/games/rl/tournaments/tournament-1/matches").
		Reply(200).
		JSON([]map[string]interface{}{
			{
				"id":          "match-live",
				"name":        "Live Match",
				"scheduledAt": "2026-01-10T18:00:00.000Z",
				"type":        "BO5",
				"teamA":       map[string]interface{}{"id": "a", "name": "Team A"},
				"teamB":       map[string]interface{}{"id": "b", "name": "Team B"},
				"teamAScore":  1,
				"teamBScore":  0,
				"maps": []map[string]interface{}{
					{"id": "map-1", "name": "Map 1", "scheduledAt": "2026-01-10T18:00:00.000Z", "startedAt": "2026-01-10T18:05:00.000Z", "endedAt": ""},
				},
			},
		})

	gock.New("https://api.blast.tv").
		Get("/v2/games/rl/tournaments/tournament-2/matches").
		Reply(200).
		JSON([]map[string]interface{}{
			{
				"id":          "match-upcoming",
				"name":        "Upcoming Match",
				"scheduledAt": "2026-01-15T18:00:00.000Z",
				"type":        "BO5",
				"teamA":       map[string]interface{}{"id": "c", "name": "Team C"},
				"teamB":       map[string]interface{}{"id": "d", "name": "Team D"},
				"teamAScore":  0,
				"teamBScore":  0,
				"maps":        []map[string]interface{}{},
			},
		})

	cmd := &TournamentsMatchesCmd{
		Output: output.GamesFormatTable,
		now: func() time.Time {
			return time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
		},
	}

	ctx := &Context{Debug: false}
	err := cmd.Run(ctx)
	assert.NoError(t, err)
	assert.True(t, gock.IsDone())
}
