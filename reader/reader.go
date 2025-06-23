package reader

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
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

// AsString returns the value as a string if possible.
func (rv ResponseValue) AsString() (string, error) {
	if rv.value == nil {
		return "", errors.New("no value present")
	}
	switch v := rv.value.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case []string:
		return strings.Join(v, ";"), nil
	default:
		return "", errors.New("unsupported type")
	}
}

// AsInt returns the value as an int if possible.
func (rv ResponseValue) AsInt() (int, error) {
	if rv.value == nil {
		return 0, errors.New("no value present")
	}
	switch v := rv.value.(type) {
	case int:
		return v, nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, errors.New("unsupported type")
	}
}

// AsStringSlice returns the value as a []string if possible.
func (rv ResponseValue) AsStringSlice() ([]string, error) {
	if rv.value == nil {
		return nil, errors.New("no value present")
	}
	switch v := rv.value.(type) {
	case []string:
		return v, nil
	case string:
		return []string{v}, nil
	default:
		return nil, errors.New("unsupported type")
	}
}

// IsPresent returns true if the value is not nil.
func (rv ResponseValue) IsPresent() bool {
	return rv.value != nil
}

// SurveyResponse maps question keys to their response value.
type SurveyResponse map[string]ResponseValue

// SurveyData holds the schema and all responses.
type SurveyData struct {
	Schema    Schema
	Responses []SurveyResponse
}

// ReadSurvey reads the Excel file and returns the parsed survey data.
func ReadSurvey(filename string) (*SurveyData, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	// Read schema
	schema := make(Schema)
	schemaRows, err := f.GetRows("schema")
	if err != nil {
		return nil, fmt.Errorf("failed to read schema sheet: %w", err)
	}
	for i, row := range schemaRows {
		if i == 0 {
			continue // skip header
		}
		if len(row) < 3 {
			continue
		}
		key := row[0]
		text := row[1]
		typeStr := row[2]
		qtype := QuestionType(typeStr)
		schema[key] = SchemaEntry{Key: key, Text: text, QType: qtype}
	}

	// Read raw data
	rawRows, err := f.GetRows("raw data")
	if err != nil {
		return nil, fmt.Errorf("failed to read \"raw data\" sheet: %w", err)
	}
	if len(rawRows) < 1 {
		return nil, fmt.Errorf("\"raw data\" sheet is empty")
	}
	headers := rawRows[0]
	var responses []SurveyResponse
	for _, row := range rawRows[1:] {
		resp := make(SurveyResponse)
		for colIdx, key := range headers {
			if colIdx >= len(row) {
				continue
			}
			val := row[colIdx]
			if val == "NA" {
				resp[key] = ResponseValue{value: nil}
				continue
			}
			entry, ok := schema[key]
			if !ok {
				continue
			}
			switch entry.QType {
			case SC:
				resp[key] = ResponseValue{value: val}
			case MC:
				if val == "" {
					resp[key] = ResponseValue{value: []string{}}
				} else {
					resp[key] = ResponseValue{value: strings.Split(val, ";")}
				}
			case TE:
				if val == "" {
					resp[key] = ResponseValue{value: nil}
				} else {
					parsed, err := strconv.Atoi(val)
					if err == nil {
						resp[key] = ResponseValue{value: parsed}
					} else {
						resp[key] = ResponseValue{value: nil}
					}
				}
			}
		}
		responses = append(responses, resp)
	}

	return &SurveyData{
		Schema:    schema,
		Responses: responses,
	}, nil
}
