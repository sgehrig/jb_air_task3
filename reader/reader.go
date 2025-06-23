package reader

import (
    "encoding/json"
    "fmt"
    "io"
    "os"

    "github.com/xuri/excelize/v2"
)

func ReadSurveyData(filename string) (*SurveyData, error) {
    f, err := excelize.OpenFile(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer f.Close()

    // Read schema
    schema := make(Schema)
    rows, err := f.GetRows("schema")
    if err != nil {
        return nil, fmt.Errorf("failed to read schema sheet: %w", err)
    }
    for i, row := range rows {
        if i == 0 {
            continue // skip header
        }
        if len(row) < 3 {
            continue
        }
        key := row[0]
        text := row[1]
        qtype := QuestionType(row[2])
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
    header := rawRows[0]
    var responses []Response
    for _, row := range rawRows[1:] {
        resp := make(Response)
        for i, cell := range row {
            if i >= len(header) {
                break
            }
            key := header[i]
            entry, ok := schema[key]
            if !ok {
                continue
            }
            resp[key] = entry.ParseValue(cell)
        }
        responses = append(responses, resp)
    }

    return &SurveyData{
        Schema:    schema,
        Responses: responses,
    }, nil
}

func LoadSurveyData(r io.Reader) (*SurveyData, error) {
    var sd SurveyData
    dec := json.NewDecoder(r)
    if err := dec.Decode(&sd); err != nil {
        return nil, err
    }
    return &sd, nil
}

func LoadSurveyDataFromFile(filename string) (*SurveyData, error) {
    f, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    return LoadSurveyData(f)
}
