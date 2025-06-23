package reader

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

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

// LoadJSON loads SurveyData from the given io.Reader.
func LoadJSON(r io.Reader) (*SurveyData, error) {
	var sd SurveyData
	dec := json.NewDecoder(r)
	if err := dec.Decode(&sd); err != nil {
		return nil, err
	}
	return &sd, nil
}

// LoadJSONFromFile loads SurveyData from the specified file path.
func LoadJSONFromFile(filename string) (*SurveyData, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadJSON(f)
}
