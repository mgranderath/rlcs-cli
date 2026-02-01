package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mgranderath/rlcs-cli/internal/api/blast"
	"github.com/mgranderath/rlcs-cli/internal/domain"
	"github.com/mgranderath/rlcs-cli/internal/mapper"
	"github.com/mgranderath/rlcs-cli/internal/output"
)

// GetBracketsCmd retrieves tournament brackets
type GetBracketsCmd struct {
	TournamentID  string                `arg:"" help:"Tournament ID (UUID)"`
	CompletedOnly bool                  `help:"Show only completed matches"`
	LiveOnly      bool                  `help:"Show only live matches"`
	UpcomingOnly  bool                  `help:"Show only upcoming matches"`
	Team          string                `help:"Filter by team name (case-insensitive partial match)"`
	MatchType     string                `help:"Filter by match type (e.g., BO5, BO7)"`
	Output        output.BracketsFormat `help:"Output format (table, json, yaml)" default:"table" short:"o"`
}

func (g *GetBracketsCmd) matchesFilters(match domain.Match) bool {
	// Status filters
	if g.CompletedOnly && !match.IsCompleted {
		return false
	}
	if g.LiveOnly && !match.IsLive {
		return false
	}
	if g.UpcomingOnly && (match.IsCompleted || match.IsLive) {
		return false
	}

	// Team filter (case-insensitive partial match)
	if g.Team != "" {
		teamLower := strings.ToLower(g.Team)
		teamAMatch := strings.Contains(strings.ToLower(match.TeamA.Name), teamLower)
		teamBMatch := strings.Contains(strings.ToLower(match.TeamB.Name), teamLower)
		if !teamAMatch && !teamBMatch {
			return false
		}
	}

	// Match type filter (case-insensitive)
	if g.MatchType != "" && !strings.EqualFold(match.Type, g.MatchType) {
		return false
	}

	return true
}

func (g *GetBracketsCmd) Run(ctx *Context) error {
	// Validate conflicting filters
	filterCount := 0
	if g.CompletedOnly {
		filterCount++
	}
	if g.LiveOnly {
		filterCount++
	}
	if g.UpcomingOnly {
		filterCount++
	}
	if filterCount > 1 {
		return fmt.Errorf("cannot use multiple status filters together (completed-only, live-only, upcoming-only are mutually exclusive)")
	}

	url := fmt.Sprintf("https://api.blast.tv/v2/games/rl/tournaments/%s/brackets", g.TournamentID)

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
		return fmt.Errorf("tournament not found: %s", g.TournamentID)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var apiBrackets []blast.Bracket
	if err := json.Unmarshal(body, &apiBrackets); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Map API response to domain model
	brackets, err := mapper.ToDomainBrackets(apiBrackets)
	if err != nil {
		return fmt.Errorf("failed to map brackets: %w", err)
	}

	// Apply filters to matches within each bracket
	brackets = g.applyFilters(brackets)

	// Get the appropriate formatter
	formatter, err := output.GetBracketsFormatter(g.Output)
	if err != nil {
		return fmt.Errorf("failed to get formatter: %w", err)
	}

	// Output using the selected formatter
	if err := formatter.Format(os.Stdout, brackets); err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	return nil
}

func (g *GetBracketsCmd) applyFilters(brackets []domain.Bracket) []domain.Bracket {
	// Check if any filters are applied
	hasFilters := g.CompletedOnly || g.LiveOnly || g.UpcomingOnly || g.Team != "" || g.MatchType != ""
	if !hasFilters {
		return brackets
	}

	result := make([]domain.Bracket, 0, len(brackets))
	for _, bracket := range brackets {
		filteredMatches := make([]domain.Match, 0)
		for _, match := range bracket.Matches {
			if g.matchesFilters(match) {
				filteredMatches = append(filteredMatches, match)
			}
		}
		if len(filteredMatches) > 0 {
			bracket.Matches = filteredMatches
			result = append(result, bracket)
		}
	}
	return result
}
