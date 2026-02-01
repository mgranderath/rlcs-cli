package output

import (
	"encoding/json"
	"io"

	"github.com/mgranderath/rlcs-cli/internal/domain"
)

// JSONFormatter outputs tournaments as formatted JSON
type JSONFormatter struct{}

func (f *JSONFormatter) Format(w io.Writer, tournaments []domain.Tournament) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tournaments)
}
