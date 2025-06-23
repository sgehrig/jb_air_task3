package reader

import (
    "fmt"
    "strconv"
    "strings"

    "github.com/xuri/excelize/v2"
)

type QuestionType string

const (
    SC QuestionType = "SC"
    MC QuestionType = "MC"
    TE QuestionType = "TE"
)

type SchemaEntry struct {
    Key   string
    Text  string
    QType QuestionType
}

func (s SchemaEntry) ParseValue(val string) ResponseValue {
    if val == "" || val == "NA" {
        return ResponseValue{val: nil}
    }
    switch s.QType {
    case SC:
        return ResponseValue{val: val}
    case MC:
        return ResponseValue{val: strings.Split(val, ";")}
    case TE:
        num, err := strconv.Atoi(val)
        if err != nil {
            return ResponseValue{val: nil}
        }
        return ResponseValue{val: num}
    default:
        return ResponseValue{val: val}
    }
}

type Schema map[string]SchemaEntry

type ResponseValue struct {
    val any
}

func (rv ResponseValue) AsInt() (int, bool) {
    v, ok := rv.val.(int)
    return v, ok
}

func (rv ResponseValue) AsString() (string, bool) {
    v, ok := rv.val.(string)
    return v, ok
}

func (rv ResponseValue) AsStringSlice() ([]string, bool) {
    v, ok := rv.val.([]string)
    return v, ok
}

func (rv ResponseValue) Present() bool {
    return rv.val != nil
}

type Response map[string]ResponseValue

type SurveyData struct {
    Schema    Schema
    Responses []Response
}

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
