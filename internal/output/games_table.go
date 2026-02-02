package output

import (
	"fmt"
	"io"

	"github.com/mgranderath/rlcs-cli/internal/domain"
)

// GamesTableFormatter outputs games as a simplified ASCII table
type GamesTableFormatter struct{}

func (f *GamesTableFormatter) Format(w io.Writer, games []domain.GameListing) error {
	if len(games) == 0 {
		fmt.Fprintln(w, "No games found")
		return nil
	}

	// Write header
	fmt.Fprintln(w, "┌───────────────────────┬───────────────────────────────┬─────────────────────────────────────┬─────────┬─────────────┐")
	fmt.Fprintln(w, "│ Tournament            │ Match                         │ Teams                               │ Score   │ Status      │")
	fmt.Fprintln(w, "├───────────────────────┼───────────────────────────────┼─────────────────────────────────────┼─────────┼─────────────┤")

	// Write games
	for _, game := range games {
		tournament := truncate(game.TournamentName, 21)
		name := truncate(game.Match.Name, 29)
		teams := fmt.Sprintf("%s vs %s", truncate(game.Match.TeamA.Name, 15), truncate(game.Match.TeamB.Name, 15))
		score := fmt.Sprintf("%d - %d", game.Match.TeamAScore, game.Match.TeamBScore)
		status := formatMatchStatus(game.Match)

		fmt.Fprintf(w, "│ %-21s │ %-29s │ %-35s │ %-7s │ %-11s │\n",
			tournament, name, teams, score, status)
	}

	fmt.Fprintln(w, "└───────────────────────┴───────────────────────────────┴─────────────────────────────────────┴─────────┴─────────────┘")

	return nil
}

func formatMatchStatus(match domain.Match) string {
	if match.IsLive {
		return "LIVE"
	}
	if match.IsCompleted {
		return "Completed"
	}
	return "Upcoming"
}
