package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/mgranderath/rlcs-cli/internal/api/blast"
	"github.com/mgranderath/rlcs-cli/internal/domain"
	"github.com/mgranderath/rlcs-cli/internal/mapper"
	"github.com/mgranderath/rlcs-cli/internal/output"
)

// GetMatchCmd retrieves detailed information for a specific match
type GetMatchCmd struct {
	MatchID string               `arg:"" help:"Match ID"`
	Output  output.MatchesFormat `help:"Output format (table, json, yaml)" default:"table" short:"o"`
}

func (g *GetMatchCmd) Run(ctx *Context) error {
	url := fmt.Sprintf("https://api.blast.tv/v2/matches/%s/detailed", g.MatchID)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("match not found: %s", g.MatchID)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var apiMatch blast.MatchResponse
	if err := json.Unmarshal(body, &apiMatch); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Map API response to domain model
	match, err := mapper.ToDomainMatchFromDetailResponse(apiMatch)
	if err != nil {
		return fmt.Errorf("failed to map match: %w", err)
	}

	// Wrap single match in a slice for formatter compatibility
	matches := []domain.Match{match}

	// Get the appropriate formatter
	formatter, err := output.GetMatchesFormatter(g.Output)
	if err != nil {
		return fmt.Errorf("failed to get formatter: %w", err)
	}

	// Output using the selected formatter
	if err := formatter.Format(os.Stdout, matches); err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	return nil
}
