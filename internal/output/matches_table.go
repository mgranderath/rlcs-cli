package output

import (
	"fmt"
	"io"

	"github.com/mgranderath/rlcs-cli/internal/domain"
)

// MatchesTableFormatter outputs matches as a simplified ASCII table
type MatchesTableFormatter struct{}

func (f *MatchesTableFormatter) Format(w io.Writer, matches []domain.Match) error {
	if len(matches) == 0 {
		fmt.Fprintln(w, "No matches found")
		return nil
	}

	// Write header
	fmt.Fprintln(w, "┌───────────────────────────────┬─────────────────────────────────────┬─────────┬─────────────┐")
	fmt.Fprintln(w, "│ Match                         │ Teams                               │ Score   │ Status      │")
	fmt.Fprintln(w, "├───────────────────────────────┼─────────────────────────────────────┼─────────┼─────────────┤")

	// Write matches
	for _, match := range matches {
		name := truncate(match.Name, 29)
		teams := fmt.Sprintf("%s vs %s", truncate(match.TeamA.Name, 15), truncate(match.TeamB.Name, 15))
		score := fmt.Sprintf("%d - %d", match.TeamAScore, match.TeamBScore)
		status := f.formatStatus(match)

		fmt.Fprintf(w, "│ %-29s │ %-35s │ %-7s │ %-11s │\n",
			name, teams, score, status)
	}

	fmt.Fprintln(w, "└───────────────────────────────┴─────────────────────────────────────┴─────────┴─────────────┘")

	return nil
}

func (f *MatchesTableFormatter) formatStatus(match domain.Match) string {
	if match.IsLive {
		return "LIVE"
	}
	if match.IsCompleted {
		return "Completed"
	}
	return "Upcoming"
}
