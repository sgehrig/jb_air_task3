package reader

import (
    "compress/gzip"
    "encoding/json"
    "io"
    "os"
    "strconv"
    "strings"
)

type QuestionType string

const (
    SC QuestionType = "SC" // Single Choice
    MC QuestionType = "MC" // Multiple Choice
    TE QuestionType = "TE" // Text Entry
)

type SchemaEntry struct {
    Key         string
    Text        string
    QType       QuestionType
    UsedOptions []string // Tracks used options for SC and MC questions
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
        return ResponseValue{val: val}
    default:
        return ResponseValue{val: nil}
    }
}

type Schema map[string]SchemaEntry

type ResponseValue struct {
    val any
}

func (rv ResponseValue) AsString() (string, bool) {
    if rv.val == nil {
        return "", false
    }
    switch v := rv.val.(type) {
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

func (rv ResponseValue) AsStringSlice() ([]string, bool) {
    if rv.val == nil {
        return nil, false
    }
    switch v := rv.val.(type) {
    case []string:
        return v, true
    case string:
        return []string{v}, true
    default:
        return nil, false
    }
}

func (rv ResponseValue) Present() bool {
    return rv.val != nil
}

type Response map[string]ResponseValue

type SurveyData struct {
    Schema    Schema
    Responses []Response
}

func (sd *SurveyData) WriteJSON(w io.Writer) error {
    enc := json.NewEncoder(w)
    enc.SetIndent("", "  ")
    return enc.Encode(sd)
}

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
