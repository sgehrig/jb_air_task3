package reader

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// QuestionType represents the type of a question in the schema.
type QuestionType string

const (
	SC QuestionType = "SC" // Single Choice
	MC QuestionType = "MC" // Multiple Choice
	TE QuestionType = "TE" // Text Entry
)

// SchemaEntry represents a single question's schema definition.
type SchemaEntry struct {
	Key     string
	Text    string
	QType   QuestionType
	Options []string // Used options for SC and MC questions
}

func (s *SchemaEntry) addUsedOptions(vals []string) {
	for _, v := range vals {
		if v == "" {
			continue
		}
		if !slices.Contains(s.Options, v) {
			s.Options = append(s.Options, v)
		}
	}
	sort.Strings(s.Options)
}

func (s *SchemaEntry) ParseValue(value string) ResponseValue {
	if value == "" || value == "NA" {
		return ResponseValue{value: nil}
	}
	switch s.QType {
	case SC:
		s.addUsedOptions([]string{value})
		return ResponseValue{value: value}
	case MC:
		vals := strings.Split(value, ";")
		s.addUsedOptions(vals)
		return ResponseValue{value: vals}
	case TE:
		return ResponseValue{value: value}
	default:
		return ResponseValue{value: nil}
	}
}

func (s *SchemaEntry) matches(str string) bool {
	if strings.Contains(strings.ToLower(s.Key), str) ||
		strings.Contains(strings.ToLower(s.Text), str) {
		return true
	}
	for _, opt := range s.Options {
		if strings.Contains(strings.ToLower(opt), str) {
			return true
		}
	}
	return false
}

// Schema maps question keys to their schema entry.
type Schema []*SchemaEntry

func (s Schema) Get(key string) (*SchemaEntry, bool) {
	for _, entry := range s {
		if entry.Key == key {
			return entry, true
		}
	}
	return nil, false
}

func (s Schema) SearchForString(str string) []*SchemaEntry {
	var out []*SchemaEntry
	str = strings.ToLower(str)
	for _, entry := range s {
		if entry.matches(str) {
			out = append(out, entry)
		}
	}
	return out
}

func (s *Schema) add(key, text string, qtype QuestionType) {
	_, exists := s.Get(key)
	if exists {
		return
	}
	*s = append(*s, &SchemaEntry{Key: key, Text: text, QType: qtype, Options: make([]string, 0)})
}

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
