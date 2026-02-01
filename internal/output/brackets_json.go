package output

import (
	"encoding/json"
	"io"

	"github.com/mgranderath/rlcs-cli/internal/domain"
)

// BracketsJSONFormatter outputs brackets as formatted JSON
type BracketsJSONFormatter struct{}

func (f *BracketsJSONFormatter) Format(w io.Writer, brackets []domain.Bracket) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(brackets)
}
