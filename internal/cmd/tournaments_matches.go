package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mgranderath/rlcs-cli/internal/api/blast"
	"github.com/mgranderath/rlcs-cli/internal/domain"
	"github.com/mgranderath/rlcs-cli/internal/mapper"
	"github.com/mgranderath/rlcs-cli/internal/output"
)

// TournamentsMatchesCmd retrieves ongoing and upcoming games across tournaments in a circuit
type TournamentsMatchesCmd struct {
	Circuit       string             `help:"Circuit/year to fetch tournaments from (e.g., 2025, 2026)" default:""`
	Region        string             `help:"Filter by region (NA, EU, APAC, SAM, OCE, MENA, SSA)"`
	Online        bool               `help:"Show only online tournaments"`
	Major         bool               `help:"Show only major tournaments (empty region/grouping)"`
	Grouping      string             `help:"Filter by tournament grouping (e.g., 'RLCS Open 1 2026')"`
	MinTeams      int                `help:"Minimum number of teams"`
	LiveOnly      bool               `help:"Show only live matches"`
	UpcomingOnly  bool               `help:"Show only upcoming matches"`
	CompletedOnly bool               `help:"Show only completed matches"`
	Limit         int                `help:"Maximum number of matches to return (after filtering)"`
	Output        output.GamesFormat `help:"Output format (table, json, yaml)" default:"table" short:"o"`

	// now is a function that returns the current time, can be overridden for testing
	now func() time.Time `kong:"-"`
}

func (l *TournamentsMatchesCmd) Run(ctx *Context) error {
	// Validate conflicting status filters
	filterCount := 0
	if l.LiveOnly {
		filterCount++
	}
	if l.UpcomingOnly {
		filterCount++
	}
	if l.CompletedOnly {
		filterCount++
	}
	if filterCount > 1 {
		return fmt.Errorf("cannot use multiple status filters together (completed-only, live-only, upcoming-only are mutually exclusive)")
	}
	if l.Limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}

	if l.now == nil {
		l.now = time.Now
	}

	circuit := l.Circuit
	if circuit == "" {
		circuit = fmt.Sprintf("%d", l.now().Year())
	}

	tournaments, err := l.fetchTournaments(circuit)
	if err != nil {
		return err
	}

	filteredTournaments := make([]domain.Tournament, 0, len(tournaments))
	for _, t := range tournaments {
		if l.matchesTournamentFilters(t) {
			filteredTournaments = append(filteredTournaments, t)
		}
	}

	games := make([]domain.GameListing, 0)
	for _, t := range filteredTournaments {
		matches, err := l.fetchMatches(t.ID)
		if err != nil {
			return err
		}

		for _, match := range matches {
			if !l.matchesStatusFilter(match) {
				continue
			}
			games = append(games, domain.GameListing{
				TournamentID:   t.ID,
				TournamentName: t.Name,
				Match:          match,
			})
		}
	}

	sortGames(games)

	if l.Limit > 0 && len(games) > l.Limit {
		games = games[:l.Limit]
	}

	formatter, err := output.GetGamesFormatter(l.Output)
	if err != nil {
		return fmt.Errorf("failed to get formatter: %w", err)
	}

	if err := formatter.Format(os.Stdout, games); err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	return nil
}

func (l *TournamentsMatchesCmd) matchesTournamentFilters(t domain.Tournament) bool {
	if l.Region != "" && !strings.EqualFold(string(t.Region), l.Region) {
		return false
	}
	if l.Online && !t.IsOnline {
		return false
	}
	if l.Major && !t.IsMajor {
		return false
	}
	if l.Grouping != "" && !strings.Contains(t.Name, l.Grouping) {
		return false
	}
	if l.MinTeams > 0 && t.TeamCount < l.MinTeams {
		return false
	}
	return true
}

func (l *TournamentsMatchesCmd) matchesStatusFilter(match domain.Match) bool {
	if l.LiveOnly {
		return match.IsLive
	}
	if l.UpcomingOnly {
		return !match.IsLive && !match.IsCompleted
	}
	if l.CompletedOnly {
		return match.IsCompleted
	}
	return match.IsLive || (!match.IsLive && !match.IsCompleted)
}

func (l *TournamentsMatchesCmd) fetchTournaments(circuit string) ([]domain.Tournament, error) {
	url := fmt.Sprintf("%s/circuits/%s/tournaments?game=rl", blast.BaseURL, circuit)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiTournaments []blast.Tournament
	if err := json.Unmarshal(body, &apiTournaments); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	tournaments, err := mapper.ToDomainTournaments(apiTournaments)
	if err != nil {
		return nil, fmt.Errorf("failed to map tournaments: %w", err)
	}

	return tournaments, nil
}

func (l *TournamentsMatchesCmd) fetchMatches(tournamentID string) ([]domain.Match, error) {
	url := fmt.Sprintf("%s/games/rl/tournaments/%s/matches", blast.BaseURL, tournamentID)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("tournament not found: %s", tournamentID)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiMatches []blast.MatchResponse
	if err := json.Unmarshal(body, &apiMatches); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	matches, err := mapper.ToDomainMatchesFromResponse(apiMatches)
	if err != nil {
		return nil, fmt.Errorf("failed to map matches: %w", err)
	}

	return matches, nil
}

func sortGames(games []domain.GameListing) {
	sort.Slice(games, func(i, j int) bool {
		a := games[i].Match
		b := games[j].Match
		rankA := matchStatusRank(a)
		rankB := matchStatusRank(b)
		if rankA != rankB {
			return rankA < rankB
		}
		if !a.TimeOfSeries.Equal(b.TimeOfSeries) {
			return a.TimeOfSeries.Before(b.TimeOfSeries)
		}
		if games[i].TournamentName != games[j].TournamentName {
			return games[i].TournamentName < games[j].TournamentName
		}
		return games[i].Match.Name < games[j].Match.Name
	})
}

func matchStatusRank(match domain.Match) int {
	if match.IsLive {
		return 0
	}
	if match.IsCompleted {
		return 2
	}
	return 1
}
