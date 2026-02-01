package output

import (
	"fmt"
	"io"

	"github.com/mgranderath/rlcs-cli/internal/domain"
)

// BracketsFormatter defines the interface for bracket output formatters
type BracketsFormatter interface {
	Format(w io.Writer, brackets []domain.Bracket) error
}

// BracketsFormat represents the output format for brackets
type BracketsFormat string

const (
	BracketsFormatTable BracketsFormat = "table"
	BracketsFormatJSON  BracketsFormat = "json"
	BracketsFormatYAML  BracketsFormat = "yaml"
)

// bracketsRegistry holds all registered bracket formatters
var bracketsRegistry = map[BracketsFormat]BracketsFormatter{
	BracketsFormatTable: &BracketsTableFormatter{},
	BracketsFormatJSON:  &BracketsJSONFormatter{},
	BracketsFormatYAML:  &BracketsYAMLFormatter{},
}

// GetBracketsFormatter returns the formatter for the given format
func GetBracketsFormatter(format BracketsFormat) (BracketsFormatter, error) {
	formatter, ok := bracketsRegistry[format]
	if !ok {
		return nil, fmt.Errorf("no formatter registered for format: %s", format)
	}
	return formatter, nil
}
