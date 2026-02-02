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

type ListTournamentsCmd struct {
	Circuit  string        `help:"Circuit/year to fetch tournaments from (e.g., 2025, 2026)" default:""`
	Region   string        `help:"Filter by region (NA, EU, APAC, SAM, OCE, MENA, SSA)"`
	Online   bool          `help:"Show only online tournaments"`
	Major    bool          `help:"Show only major tournaments (empty region/grouping)"`
	Grouping string        `help:"Filter by tournament grouping (e.g., 'RLCS Open 1 2026')"`
	Upcoming bool          `help:"Show only upcoming tournaments (start date > today)"`
	Ongoing  bool          `help:"Show only ongoing tournaments (start date <= today <= end date)"`
	Past     bool          `help:"Show only past tournaments (end date < today)"`
	MinTeams int           `help:"Minimum number of teams"`
	Output   output.Format `help:"Output format (table, json, csv, yaml)" default:"table" short:"o"`

	// now is a function that returns the current time, can be overridden for testing
	now func() time.Time `kong:"-"`
}

func (l *ListTournamentsCmd) matchesFilters(t domain.Tournament, today time.Time) bool {
	// Region filter (case-insensitive)
	if l.Region != "" && !strings.EqualFold(string(t.Region), l.Region) {
		return false
	}

	// Online filter
	if l.Online && !t.IsOnline {
		return false
	}

	// Major filter
	if l.Major && !t.IsMajor {
		return false
	}

	// Grouping filter (partial match on name for simplicity)
	if l.Grouping != "" && !strings.Contains(t.Name, l.Grouping) {
		return false
	}

	// MinTeams filter
	if l.MinTeams > 0 && t.TeamCount < l.MinTeams {
		return false
	}

	// Temporal filters
	// Check for conflicting filters
	if l.Upcoming && l.Past {
		return false // Cannot be both upcoming and past
	}

	// Upcoming: start date > today
	if l.Upcoming && !t.IsUpcoming(today) {
		return false
	}

	// Past: end date < today
	if l.Past && !t.IsPast(today) {
		return false
	}

	// Ongoing: start date <= today <= end date
	if l.Ongoing && !t.IsOngoing(today) {
		return false
	}

	return true
}

func (l *ListTournamentsCmd) Run(ctx *Context) error {
	// Validate conflicting temporal filters
	if l.Upcoming && l.Past {
		return fmt.Errorf("cannot use --upcoming and --past together (they are mutually exclusive)")
	}

	// Initialize now function if not set (allows for dependency injection in tests)
	if l.now == nil {
		l.now = time.Now
	}

	// Determine circuit - use provided value or default to current year
	circuit := l.Circuit
	if circuit == "" {
		circuit = fmt.Sprintf("%d", l.now().Year())
	}

	url := fmt.Sprintf("%s/circuits/%s/tournaments?game=rl", blast.BaseURL, circuit)

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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var apiTournaments []blast.Tournament
	if err := json.Unmarshal(body, &apiTournaments); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Map API response to domain model
	tournaments, err := mapper.ToDomainTournaments(apiTournaments)
	if err != nil {
		return fmt.Errorf("failed to map tournaments: %w", err)
	}

	// Apply filters
	today := l.now().Truncate(24 * time.Hour)
	var filtered []domain.Tournament
	for _, t := range tournaments {
		if l.matchesFilters(t, today) {
			filtered = append(filtered, t)
		}
	}

	// Get the appropriate formatter
	formatter, err := output.GetFormatter(l.Output)
	if err != nil {
		return fmt.Errorf("failed to get formatter: %w", err)
	}

	// Output using the selected formatter
	if err := formatter.Format(os.Stdout, filtered); err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	return nil
}
