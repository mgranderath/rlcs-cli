package output

import (
	"fmt"
	"io"

	"github.com/mgranderath/rlcs-cli/internal/domain"
)

// GamesFormatter defines the interface for games output formatters
type GamesFormatter interface {
	Format(w io.Writer, games []domain.GameListing) error
}

// GamesFormat represents the output format for games
type GamesFormat string

const (
	GamesFormatTable GamesFormat = "table"
	GamesFormatJSON  GamesFormat = "json"
	GamesFormatYAML  GamesFormat = "yaml"
)

// gamesRegistry holds all registered games formatters
var gamesRegistry = map[GamesFormat]GamesFormatter{
	GamesFormatTable: &GamesTableFormatter{},
	GamesFormatJSON:  &GamesJSONFormatter{},
	GamesFormatYAML:  &GamesYAMLFormatter{},
}

// GetGamesFormatter returns the formatter for the given format
func GetGamesFormatter(format GamesFormat) (GamesFormatter, error) {
	formatter, ok := gamesRegistry[format]
	if !ok {
		return nil, fmt.Errorf("no formatter registered for format: %s", format)
	}
	return formatter, nil
}
