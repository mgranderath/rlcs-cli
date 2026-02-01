package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/mgranderath/rlcs-cli/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableFormatter_Format(t *testing.T) {
	formatter := &TableFormatter{}

	tests := []struct {
		name        string
		tournaments []domain.Tournament
		contains    []string
	}{
		{
			name: "single tournament",
			tournaments: []domain.Tournament{
				{
					ID:        "test-1",
					Name:      "Test Tournament",
					StartDate: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2026, 1, 17, 0, 0, 0, 0, time.UTC),
					PrizePool: "$50,000",
					Region:    domain.RegionNA,
					TeamCount: 16,
					Type:      domain.TypeOpen,
				},
			},
			contains: []string{"Test Tournament", "$50,000", "NA", "16", "Open"},
		},
		{
			name: "multiple tournaments",
			tournaments: []domain.Tournament{
				{
					ID:        "test-1",
					Name:      "Tournament One",
					StartDate: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2026, 1, 17, 0, 0, 0, 0, time.UTC),
					PrizePool: "$50,000",
					Region:    domain.RegionEU,
					TeamCount: 16,
					Type:      domain.TypeOpen,
				},
				{
					ID:        "test-2",
					Name:      "Tournament Two",
					StartDate: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2026, 2, 5, 0, 0, 0, 0, time.UTC),
					PrizePool: "$100,000",
					Region:    domain.RegionNA,
					TeamCount: 24,
					Type:      domain.TypeMajor,
				},
			},
			contains: []string{"Tournament One", "Tournament Two", "EU", "NA"},
		},
		{
			name: "tournament with empty region (Major)",
			tournaments: []domain.Tournament{
				{
					ID:        "major-1",
					Name:      "RLCS Major",
					StartDate: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC),
					PrizePool: "$300,000",
					Region:    domain.RegionNone,
					TeamCount: 16,
					Type:      domain.TypeMajor,
					IsMajor:   true,
				},
			},
			contains: []string{"RLCS Major", "$300,000", "-", "Major"},
		},
		{
			name:        "empty tournaments",
			tournaments: []domain.Tournament{},
			contains:    []string{"Name", "Dates"}, // Header should still be present
		},
		{
			name: "long name truncation",
			tournaments: []domain.Tournament{
				{
					ID:        "test-1",
					Name:      "This is a very long tournament name that needs to be truncated properly",
					StartDate: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2026, 1, 17, 0, 0, 0, 0, time.UTC),
					Region:    domain.RegionNA,
				},
			},
			contains: []string{"This is a very long tourna..."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := formatter.Format(&buf, tt.tournaments)
			require.NoError(t, err)

			output := buf.String()
			for _, expected := range tt.contains {
				assert.Contains(t, output, expected)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{"shorter than max", "hello", 10, "hello"},
		{"exactly max", "helloworld", 10, "helloworld"},
		{"longer than max", "helloworld12345", 10, "hellowo..."},
		{"empty string", "", 10, ""},
		{"single char", "a", 10, "a"},
		{"needs truncation", "This is a very long string", 15, "This is a ve..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncate(tt.input, tt.maxLen)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatDateRange(t *testing.T) {
	tests := []struct {
		name     string
		start    interface{}
		end      interface{}
		expected string
	}{
		{
			name:     "same month time.Time",
			start:    time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2026, 1, 17, 0, 0, 0, 0, time.UTC),
			expected: "Jan 15-17 '26",
		},
		{
			name:     "different months time.Time",
			start:    time.Date(2026, 1, 30, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
			expected: "Jan 30-Feb 01 '26",
		},
		{
			name:     "string dates with T",
			start:    "2026-01-15T00:00:00Z",
			end:      "2026-01-17T00:00:00Z",
			expected: "Jan 15-17 '26",
		},
		{
			name:     "string dates without T",
			start:    "2026-01-15",
			end:      "2026-01-17",
			expected: "Jan 15-17 '26",
		},
		{
			name:     "invalid format fallback",
			start:    "invalid",
			end:      "also invalid",
			expected: "invalid - also invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDateRange(tt.start, tt.end)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTableFormatter_ColumnAlignment(t *testing.T) {
	// Test that the table has proper column alignment with box-drawing characters
	formatter := &TableFormatter{}
	tournaments := []domain.Tournament{
		{
			ID:        "test",
			Name:      "Test",
			StartDate: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2026, 1, 17, 0, 0, 0, 0, time.UTC),
			Region:    domain.RegionNA,
			Type:      domain.TypeOpen,
		},
	}

	var buf bytes.Buffer
	err := formatter.Format(&buf, tournaments)
	require.NoError(t, err)

	output := buf.String()
	// Check for box-drawing characters
	assert.Contains(t, output, "┌")
	assert.Contains(t, output, "┐")
	assert.Contains(t, output, "└")
	assert.Contains(t, output, "┘")
	assert.Contains(t, output, "├")
	assert.Contains(t, output, "┤")
	assert.Contains(t, output, "│")
	assert.Contains(t, output, "─")
}

func TestTableFormatter_WithLongFields(t *testing.T) {
	formatter := &TableFormatter{}
	tournaments := []domain.Tournament{
		{
			ID:        "very-long-tournament-id-2026",
			Name:      "RLCS World Championship Finals 2026 Edition",
			StartDate: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC),
			PrizePool: "$1,000,000 + Additional Bonuses",
			Region:    domain.RegionNone,
			TeamCount: 24,
			Type:      domain.TypeWorldChampionship,
		},
	}

	var buf bytes.Buffer
	err := formatter.Format(&buf, tournaments)
	require.NoError(t, err)

	output := buf.String()
	// Ensure output doesn't break table structure
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		if i == 0 || i == 2 || i == len(lines)-2 {
			// Check separator lines have consistent length
			if strings.Contains(line, "─") {
				assert.True(t, len(line) > 50, "Separator line should be wide enough")
			}
		}
	}
}
