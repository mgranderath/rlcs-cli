package output

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/mgranderath/rlcs-cli/internal/domain"
)

// CSVFormatter outputs tournaments as CSV
type CSVFormatter struct{}

func (f *CSVFormatter) Format(w io.Writer, tournaments []domain.Tournament) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Name", "StartDate", "EndDate", "CircuitID", "PrizePool", "Location", "TeamCount", "Region", "Type", "Description", "IsOnline", "IsMajor"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write rows
	for _, t := range tournaments {
		record := []string{
			t.ID,
			t.Name,
			t.StartDate.Format("2006-01-02"),
			t.EndDate.Format("2006-01-02"),
			t.CircuitID,
			t.PrizePool,
			t.Location,
			fmt.Sprintf("%d", t.TeamCount),
			string(t.Region),
			string(t.Type),
			t.Description,
			fmt.Sprintf("%t", t.IsOnline),
			fmt.Sprintf("%t", t.IsMajor),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}
