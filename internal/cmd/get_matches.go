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

// GetMatchesCmd retrieves all matches for a tournament
type GetMatchesCmd struct {
	TournamentID  string               `arg:"" help:"Tournament ID"`
	CompletedOnly bool                 `help:"Show only completed matches"`
	LiveOnly      bool                 `help:"Show only live matches"`
	UpcomingOnly  bool                 `help:"Show only upcoming matches"`
	Team          string               `help:"Filter by team name (case-insensitive partial match)"`
	MatchType     string               `help:"Filter by match type (e.g., BO5, BO7)"`
	Output        output.MatchesFormat `help:"Output format (table, json, yaml)" default:"table" short:"o"`
}

func (g *GetMatchesCmd) matchesFilters(match domain.Match) bool {
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

	// Team filter (case-insensitive partial match on name or shorthand)
	if g.Team != "" {
		teamLower := strings.ToLower(g.Team)
		teamAMatch := strings.Contains(strings.ToLower(match.TeamA.Name), teamLower) ||
			strings.Contains(strings.ToLower(match.TeamA.Shorthand), teamLower)
		teamBMatch := strings.Contains(strings.ToLower(match.TeamB.Name), teamLower) ||
			strings.Contains(strings.ToLower(match.TeamB.Shorthand), teamLower)
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

func (g *GetMatchesCmd) Run(ctx *Context) error {
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

	url := fmt.Sprintf("https://api.blast.tv/v2/games/rl/tournaments/%s/matches", g.TournamentID)

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

	var apiMatches []blast.MatchResponse
	if err := json.Unmarshal(body, &apiMatches); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Map API response to domain model
	matches, err := mapper.ToDomainMatchesFromResponse(apiMatches)
	if err != nil {
		return fmt.Errorf("failed to map matches: %w", err)
	}

	// Apply filters
	matches = g.applyFilters(matches)

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

func (g *GetMatchesCmd) applyFilters(matches []domain.Match) []domain.Match {
	// Check if any filters are applied
	hasFilters := g.CompletedOnly || g.LiveOnly || g.UpcomingOnly || g.Team != "" || g.MatchType != ""
	if !hasFilters {
		return matches
	}

	result := make([]domain.Match, 0, len(matches))
	for _, match := range matches {
		if g.matchesFilters(match) {
			result = append(result, match)
		}
	}
	return result
}
