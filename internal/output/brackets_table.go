package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/mgranderath/rlcs-cli/internal/domain"
)

// BracketsTableFormatter outputs brackets as a simplified ASCII table
type BracketsTableFormatter struct{}

func (f *BracketsTableFormatter) Format(w io.Writer, brackets []domain.Bracket) error {
	if len(brackets) == 0 {
		fmt.Fprintln(w, "No brackets found")
		return nil
	}

	// Display all brackets
	for i, bracket := range brackets {
		// Add separator between brackets (except before the first one)
		if i > 0 {
			fmt.Fprintln(w, "\n"+strings.Repeat("=", 80))
		}

		// Write bracket header
		fmt.Fprintf(w, "\n%s (%s)\n", bracket.TournamentName, bracket.Label)
		if bracket.ParentTournamentName != "" {
			fmt.Fprintf(w, "Part of: %s\n", bracket.ParentTournamentName)
		}
		fmt.Fprintln(w)

		// Write matches table header
		fmt.Fprintln(w, "┌───────────────────────────────┬─────────────────────────────────────┬─────────┬─────────────┐")
		fmt.Fprintln(w, "│ Match                         │ Teams                               │ Score   │ Status      │")
		fmt.Fprintln(w, "├───────────────────────────────┼─────────────────────────────────────┼─────────┼─────────────┤")

		// Write matches
		for _, match := range bracket.Matches {
			name := truncate(match.Name, 29)
			teams := fmt.Sprintf("%s vs %s", truncate(match.TeamA.Name, 15), truncate(match.TeamB.Name, 15))
			score := fmt.Sprintf("%d - %d", match.TeamAScore, match.TeamBScore)
			status := f.formatStatus(match)

			fmt.Fprintf(w, "│ %-29s │ %-35s │ %-7s │ %-11s │\n",
				name, teams, score, status)
		}

		fmt.Fprintln(w, "└───────────────────────────────┴─────────────────────────────────────┴─────────┴─────────────┘")
	}

	return nil
}

func (f *BracketsTableFormatter) formatStatus(match domain.Match) string {
	if match.IsLive {
		return "LIVE"
	}
	if match.IsCompleted {
		return "Completed"
	}
	return "Upcoming"
}
