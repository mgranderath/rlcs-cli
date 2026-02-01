package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat_Valid(t *testing.T) {
	tests := []struct {
		name     string
		format   Format
		expected bool
	}{
		{"table is valid", FormatTable, true},
		{"json is valid", FormatJSON, true},
		{"csv is valid", FormatCSV, true},
		{"yaml is valid", FormatYAML, true},
		{"uppercase TABLE", Format("TABLE"), false},
		{"invalid format", Format("xml"), false},
		{"empty format", Format(""), false},
		{"random string", Format("random"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.format.Valid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormat_String(t *testing.T) {
	tests := []struct {
		name     string
		format   Format
		expected string
	}{
		{"table", FormatTable, "table"},
		{"json", FormatJSON, "json"},
		{"csv", FormatCSV, "csv"},
		{"yaml", FormatYAML, "yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.format.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormat_UnmarshalFlag(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    Format
		expectError bool
	}{
		{"valid table", "table", FormatTable, false},
		{"valid json", "json", FormatJSON, false},
		{"valid csv", "csv", FormatCSV, false},
		{"valid yaml", "yaml", FormatYAML, false},
		{"invalid format", "xml", "", true},
		{"empty format", "", "", true},
		{"mixed case table", "TABLE", "", true},
		{"json with spaces", " json ", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var format Format
			err := format.UnmarshalFlag(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, format)
		})
	}
}
