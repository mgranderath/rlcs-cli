package output

import "fmt"

// Format is a strongly-typed output format type
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatCSV   Format = "csv"
	FormatYAML  Format = "yaml"
)

// Valid checks if the format is supported
func (f Format) Valid() bool {
	switch f {
	case FormatTable, FormatJSON, FormatCSV, FormatYAML:
		return true
	}
	return false
}

// String implements fmt.Stringer
func (f Format) String() string {
	return string(f)
}

// UnmarshalFlag implements kong.FlagUnmarshaller for CLI parsing
func (f *Format) UnmarshalFlag(value string) error {
	format := Format(value)
	if !format.Valid() {
		return fmt.Errorf("invalid output format %q, must be one of: table, json, csv, yaml", value)
	}
	*f = format
	return nil
}
