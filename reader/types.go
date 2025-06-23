package reader

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"
)

// QuestionType represents the type of a question in the schema.
type QuestionType string

const (
	SC QuestionType = "SC" // Single Choice
	MC QuestionType = "MC" // Multiple Choice
	TE QuestionType = "TE" // Text Entry (integer)
)

// SchemaEntry represents a single question's schema definition.
type SchemaEntry struct {
	Key   string
	Text  string
	QType QuestionType
}

// Schema maps question keys to their schema entry.
type Schema map[string]SchemaEntry

// ResponseValue holds the value for a question, respecting its type.
type ResponseValue struct {
	value any
}

// SurveyResponse maps question keys to their response value.
type SurveyResponse map[string]ResponseValue

// SurveyData holds the schema and all responses.
type SurveyData struct {
	Schema    Schema
	Responses []SurveyResponse
}

// AsString returns the value as a string if possible.
func (rv ResponseValue) AsString() (string, bool) {
	if rv.value == nil {
		return "", false
	}
	switch v := rv.value.(type) {
	case string:
		return v, true
	case int:
		return strconv.Itoa(v), true
	case []string:
		return strings.Join(v, ";"), true
	default:
		return "", false
	}
}

// AsInt returns the value as an int if possible.
func (rv ResponseValue) AsInt() (int, bool) {
	if rv.value == nil {
		return 0, false
	}
	switch v := rv.value.(type) {
	case int:
		return v, true
	case string:
		ival, err := strconv.Atoi(v)
		if err == nil {
			return ival, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// AsStringSlice returns the value as a []string if possible.
func (rv ResponseValue) AsStringSlice() ([]string, bool) {
	if rv.value == nil {
		return nil, false
	}
	switch v := rv.value.(type) {
	case []string:
		return v, true
	case string:
		return []string{v}, true
	default:
		return nil, false
	}
}

// IsPresent returns true if the value is not nil.
func (rv ResponseValue) IsPresent() bool {
	return rv.value != nil
}

// WriteJSON writes the SurveyData as JSON to the given io.Writer.
func (sd *SurveyData) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(sd)
}

// WriteJSONToFile writes the SurveyData as JSON to the specified file path, gzip-compressed if the filename ends with .gz.
func (sd *SurveyData) WriteJSONToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if len(filename) > 3 && filename[len(filename)-3:] == ".gz" {
		gw := gzip.NewWriter(f)
		defer gw.Close()
		return sd.WriteJSON(gw)
	}
	return sd.WriteJSON(f)
}
