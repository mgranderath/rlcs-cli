package output

import (
	"encoding/json"
	"io"

	"github.com/mgranderath/rlcs-cli/internal/domain"
)

// GamesJSONFormatter outputs games as formatted JSON
type GamesJSONFormatter struct{}

func (f *GamesJSONFormatter) Format(w io.Writer, games []domain.GameListing) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(games)
}
