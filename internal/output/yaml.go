package output

import (
	"io"

	"github.com/mgranderath/rlcs-cli/internal/domain"
	"gopkg.in/yaml.v3"
)

// YAMLFormatter outputs tournaments as YAML
type YAMLFormatter struct{}

func (f *YAMLFormatter) Format(w io.Writer, tournaments []domain.Tournament) error {
	encoder := yaml.NewEncoder(w)
	defer encoder.Close()
	encoder.SetIndent(2)
	return encoder.Encode(tournaments)
}
