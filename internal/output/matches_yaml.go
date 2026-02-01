package output

import (
	"io"

	"github.com/mgranderath/rlcs-cli/internal/domain"
	"gopkg.in/yaml.v3"
)

// MatchesYAMLFormatter outputs matches as YAML
type MatchesYAMLFormatter struct{}

func (f *MatchesYAMLFormatter) Format(w io.Writer, matches []domain.Match) error {
	encoder := yaml.NewEncoder(w)
	defer encoder.Close()
	encoder.SetIndent(2)
	return encoder.Encode(matches)
}
