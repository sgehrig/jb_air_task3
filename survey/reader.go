package survey

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

// ReadSurvey reads the Excel file and returns the parsed survey data.
func ReadSurvey(filename string) (*SurveyData, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	// Read schema
	schema := make(Schema, 0)
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
		schema.add(key, text, qtype)
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
			entry, ok := schema.Get(key)
			if !ok {
				continue
			}
			resp[key] = entry.ParseValue(row[colIdx])
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

// LoadJSONFromFile loads SurveyData from the specified file path, supporting gzip if the filename ends with .gz.
func LoadJSONFromFile(filename string) (*SurveyData, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if len(filename) > 3 && filename[len(filename)-3:] == ".gz" {
		gr, err := gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
		defer gr.Close()
		return LoadJSON(gr)
	}
	return LoadJSON(f)
}

func createCacheFilename(xlsxFile string) string {
	// Generate cache file name based on Excel file name, strip extension
	base := filepath.Base(xlsxFile)
	name := base
	if ext := filepath.Ext(base); ext != "" {
		name = base[:len(base)-len(ext)]
	}
	return "_" + name + ".cache.json.gz"
}

// ReadSurveyCached tries to load survey data from a JSON cache file derived from the XLSX filename.
// If not, it reads from the XLSX file, writes the JSON cache, and returns the data.
func ReadSurveyCached(xlsxFile string) (*SurveyData, error) {
	jsonFile := createCacheFilename(xlsxFile)

	// Try to load from JSON first
	if _, err := os.Stat(jsonFile); err == nil {
		data, err := LoadJSONFromFile(jsonFile)
		if err == nil {
			return data, nil
		}
		// If JSON exists but is invalid, fall back to XLSX
	}
	// Fallback: load from XLSX and write JSON
	if _, err := os.Stat(xlsxFile); err != nil {
		return nil, fmt.Errorf("could not find %s or %s", jsonFile, xlsxFile)
	}
	data, err := ReadSurvey(xlsxFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", xlsxFile, err)
	}
	if err := data.WriteJSONToFile(jsonFile); err != nil {
		return nil, fmt.Errorf("failed to write %s: %w", jsonFile, err)
	}
	return data, nil
}
