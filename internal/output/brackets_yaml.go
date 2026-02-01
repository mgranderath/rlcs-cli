package output

import (
	"io"

	"github.com/mgranderath/rlcs-cli/internal/domain"
	"gopkg.in/yaml.v3"
)

// BracketsYAMLFormatter outputs brackets as YAML
type BracketsYAMLFormatter struct{}

func (f *BracketsYAMLFormatter) Format(w io.Writer, brackets []domain.Bracket) error {
	encoder := yaml.NewEncoder(w)
	defer encoder.Close()
	encoder.SetIndent(2)
	return encoder.Encode(brackets)
}
