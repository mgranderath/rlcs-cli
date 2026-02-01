package output

import (
	"fmt"
	"io"

	"github.com/mgranderath/rlcs-cli/internal/domain"
	"gopkg.in/yaml.v3"
)

// MatchesYAMLFormatter outputs matches as YAML
type MatchesYAMLFormatter struct{}

func (f *MatchesYAMLFormatter) Format(w io.Writer, matches []domain.Match) error {
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)

	// Encode the matches
	if err := encoder.Encode(matches); err != nil {
		encoder.Close()
		return fmt.Errorf("failed to encode matches to YAML: %w", err)
	}

	// Close the encoder and handle any flush errors
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("failed to close YAML encoder: %w", err)
	}

	return nil
}
