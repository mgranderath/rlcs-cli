package output

import (
	"fmt"
	"io"

	"github.com/mgranderath/rlcs-cli/internal/domain"
	"gopkg.in/yaml.v3"
)

// GamesYAMLFormatter outputs games as YAML
type GamesYAMLFormatter struct{}

func (f *GamesYAMLFormatter) Format(w io.Writer, games []domain.GameListing) error {
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)

	if err := encoder.Encode(games); err != nil {
		return fmt.Errorf("failed to encode games to YAML: %w", err)
	}

	if err := encoder.Close(); err != nil {
		return fmt.Errorf("failed to close YAML encoder: %w", err)
	}

	return nil
}
