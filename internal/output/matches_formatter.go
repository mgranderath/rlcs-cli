package output

import (
	"fmt"
	"io"

	"github.com/mgranderath/rlcs-cli/internal/domain"
)

// MatchesFormatter defines the interface for matches output formatters
type MatchesFormatter interface {
	Format(w io.Writer, matches []domain.Match) error
}

// MatchesFormat represents the output format for matches
type MatchesFormat string

const (
	MatchesFormatTable MatchesFormat = "table"
	MatchesFormatJSON  MatchesFormat = "json"
	MatchesFormatYAML  MatchesFormat = "yaml"
)

// matchesRegistry holds all registered matches formatters
var matchesRegistry = map[MatchesFormat]MatchesFormatter{
	MatchesFormatTable: &MatchesTableFormatter{},
	MatchesFormatJSON:  &MatchesJSONFormatter{},
	MatchesFormatYAML:  &MatchesYAMLFormatter{},
}

// GetMatchesFormatter returns the formatter for the given format
func GetMatchesFormatter(format MatchesFormat) (MatchesFormatter, error) {
	formatter, ok := matchesRegistry[format]
	if !ok {
		return nil, fmt.Errorf("no formatter registered for format: %s", format)
	}
	return formatter, nil
}
