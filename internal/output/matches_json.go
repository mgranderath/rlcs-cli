package output

import (
	"encoding/json"
	"io"

	"github.com/mgranderath/rlcs-cli/internal/domain"
)

// MatchesJSONFormatter outputs matches as formatted JSON
type MatchesJSONFormatter struct{}

func (f *MatchesJSONFormatter) Format(w io.Writer, matches []domain.Match) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(matches)
}
