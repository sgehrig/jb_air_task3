package survey

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

func (s *SchemaEntry) addUsedOptions(vals []string) {
    for _, v := range vals {
        if v == "" {
            continue
        }
        if !slices.Contains(s.UsedOptions, v) {
            s.UsedOptions = append(s.UsedOptions, v)
        }
    }
    sort.Strings(s.UsedOptions)
}

func (s *SchemaEntry) ParseValue(val string) ResponseValue {
    if val == "" || val == "NA" {
        return ResponseValue{Val: nil}
    }
    switch s.QType {
    case SC:
        s.addUsedOptions([]string{val})
        return ResponseValue{Val: val}
    case MC:
        vals := strings.Split(val, ";")
        s.addUsedOptions(vals)
        return ResponseValue{Val: vals}
    case TE:
        return ResponseValue{Val: val}
    default:
        return ResponseValue{Val: nil}
    }
}

func (s *SchemaEntry) matches(str string) bool {
    if strings.Contains(strings.ToLower(s.Key), str) ||
        strings.Contains(strings.ToLower(s.Text), str) {
        return true
    }
    for _, opt := range s.UsedOptions {
        if strings.Contains(strings.ToLower(opt), str) {
            return true
        }
    }
    return false
}

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
    *s = append(*s, &SchemaEntry{Key: key, Text: text, QType: qtype, UsedOptions: make([]string, 0)})
}

type ResponseValue struct {
    Val any
}

func (rv ResponseValue) AsString() (string, bool) {
    if rv.Val == nil {
        return "", false
    }
    switch v := rv.Val.(type) {
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
    if rv.Val == nil {
        return nil, false
    }
    switch v := rv.Val.(type) {
    case []string:
        return v, true
    case []any:
        var out []string
        for _, elem := range v {
            s, ok := elem.(string)
            if !ok {
                return nil, false
            }
            out = append(out, s)
        }
        return out, true
    case string:
        return []string{v}, true
    default:
        return nil, false
    }
}

func (rv ResponseValue) Present() bool {
    return rv.Val != nil
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

func (sd *SurveyData) CreateSubset(questionKey string, optionSearch string) []Response {
    var result []Response
    entry, found := sd.Schema.Get(questionKey)
    if !found {
        return result
    }
    optionSearchLower := strings.ToLower(optionSearch)
    for _, resp := range sd.Responses {
        val, ok := resp[questionKey]
        if !ok || !val.Present() {
            continue
        }
        switch entry.QType {
        case SC:
            s, ok := val.AsString()
            if ok && strings.Contains(strings.ToLower(s), optionSearchLower) {
                result = append(result, resp)
            }
        case MC:
            ss, ok := val.AsStringSlice()
            if ok {
                for _, opt := range ss {
                    if strings.Contains(strings.ToLower(opt), optionSearchLower) {
                        result = append(result, resp)
                        break
                    }
                }
            }
        }
    }
    return result
}
