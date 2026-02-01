package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/mgranderath/rlcs-cli/internal/domain"
)

// TableFormatter outputs tournaments as an ASCII table
type TableFormatter struct{}

func (f *TableFormatter) Format(w io.Writer, tournaments []domain.Tournament) error {
	// Write header
	fmt.Fprintln(w, "┌────────────────────────────┬───────────────────────────────┬──────────────────┬─────────────┬────────┬───────┬────────────────────┐")
	fmt.Fprintln(w, "│ ID                         │ Name                          │ Dates            │ Prize Pool  │ Region │ Teams │ Type               │")
	fmt.Fprintln(w, "├────────────────────────────┼───────────────────────────────┼──────────────────┼─────────────┼────────┼───────┼────────────────────┤")

	for _, t := range tournaments {
		id := truncate(t.ID, 26)
		name := truncate(t.Name, 29)
		dates := formatDateRange(t.StartDate, t.EndDate)
		prizePool := truncate(t.PrizePool, 13)
		region := string(t.Region)
		if region == "" {
			region = "-"
		}
		teams := fmt.Sprintf("%d", t.TeamCount)
		typ := string(t.Type)

		fmt.Fprintf(w, "│ %-26s │ %-29s │ %-16s │ %-11s │ %-6s │ %-5s │ %-18s │\n",
			id, name, dates, prizePool, region, teams, typ)
	}

	fmt.Fprintln(w, "└────────────────────────────┴───────────────────────────────┴──────────────────┴─────────────┴────────┴───────┴────────────────────┘")
	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatDateRange(start, end interface{}) string {
	// Handle both string and time.Time
	startStr := fmt.Sprintf("%v", start)
	endStr := fmt.Sprintf("%v", end)

	// Parse if they contain T (time format)
	if strings.Contains(startStr, "T") {
		startStr = startStr[:10]
	}
	if strings.Contains(endStr, "T") {
		endStr = endStr[:10]
	}

	// If same month, show abbreviated format
	if len(startStr) >= 10 && len(endStr) >= 10 {
		startMonth := startStr[5:7]
		endMonth := endStr[5:7]
		startDay := startStr[8:10]
		endDay := endStr[8:10]
		year := startStr[2:4]

		monthNames := map[string]string{
			"01": "Jan", "02": "Feb", "03": "Mar", "04": "Apr",
			"05": "May", "06": "Jun", "07": "Jul", "08": "Aug",
			"09": "Sep", "10": "Oct", "11": "Nov", "12": "Dec",
		}

		if startMonth == endMonth {
			return fmt.Sprintf("%s %s-%s '%s", monthNames[startMonth], startDay, endDay, year)
		}
		return fmt.Sprintf("%s %s-%s %s '%s", monthNames[startMonth], startDay, monthNames[endMonth], endDay, year)
	}

	return fmt.Sprintf("%s - %s", startStr, endStr)
}
