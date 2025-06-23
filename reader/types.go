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
