package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMatchesFormatter(t *testing.T) {
	tests := []struct {
		name        string
		format      MatchesFormat
		expectError bool
		checkType   string
	}{
		{
			name:        "table format",
			format:      MatchesFormatTable,
			expectError: false,
			checkType:   "*output.MatchesTableFormatter",
		},
		{
			name:        "json format",
			format:      MatchesFormatJSON,
			expectError: false,
			checkType:   "*output.MatchesJSONFormatter",
		},
		{
			name:        "yaml format",
			format:      MatchesFormatYAML,
			expectError: false,
			checkType:   "*output.MatchesYAMLFormatter",
		},
		{
			name:        "invalid format",
			format:      MatchesFormat("xml"),
			expectError: true,
		},
		{
			name:        "empty format",
			format:      MatchesFormat(""),
			expectError: true,
		},
		{
			name:        "uppercase TABLE",
			format:      MatchesFormat("TABLE"),
			expectError: true,
		},
		{
			name:        "mixed case json",
			format:      MatchesFormat("Json"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter, err := GetMatchesFormatter(tt.format)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, formatter)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, formatter)
			assert.Equal(t, tt.checkType, formatterTypeName(formatter))
		})
	}
}

// Helper to get type name for verification
func formatterTypeName(f MatchesFormatter) string {
	if f == nil {
		return "nil"
	}
	// Return string representation of the type
	switch f.(type) {
	case *MatchesTableFormatter:
		return "*output.MatchesTableFormatter"
	case *MatchesJSONFormatter:
		return "*output.MatchesJSONFormatter"
	case *MatchesYAMLFormatter:
		return "*output.MatchesYAMLFormatter"
	default:
		return "unknown"
	}
}
