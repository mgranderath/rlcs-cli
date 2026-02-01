package output

import (
	"fmt"
	"io"

	"github.com/mgranderath/rlcs-cli/internal/domain"
)

// Formatter defines the interface for output formatters
type Formatter interface {
	Format(w io.Writer, tournaments []domain.Tournament) error
}

// registry holds all registered formatters
var registry = map[Format]Formatter{
	FormatJSON:  &JSONFormatter{},
	FormatTable: &TableFormatter{},
	FormatCSV:   &CSVFormatter{},
	FormatYAML:  &YAMLFormatter{},
}

// GetFormatter returns the formatter for the given format
func GetFormatter(format Format) (Formatter, error) {
	formatter, ok := registry[format]
	if !ok {
		return nil, fmt.Errorf("no formatter registered for format: %s", format)
	}
	return formatter, nil
}
