package reader

import (
    "bytes"
    "encoding/json"
    "testing"
)

func TestResponseValue_AsMethods(t *testing.T) {
    rvInt := ResponseValue{val: 42}
    rvStr := ResponseValue{val: "foo"}
    rvSlice := ResponseValue{val: []string{"a", "b"}}
    rvNil := ResponseValue{val: nil}

    if v, ok := rvInt.AsInt(); !ok || v != 42 {
        t.Errorf("AsInt failed: got %v, %v", v, ok)
    }
    if v, ok := rvStr.AsString(); !ok || v != "foo" {
        t.Errorf("AsString failed: got %v, %v", v, ok)
    }
    if v, ok := rvSlice.AsStringSlice(); !ok || len(v) != 2 || v[0] != "a" || v[1] != "b" {
        t.Errorf("AsStringSlice failed: got %v, %v", v, ok)
    }
    if rvNil.Present() {
        t.Error("Present should be false for nil value")
    }
    if !rvInt.Present() || !rvStr.Present() || !rvSlice.Present() {
        t.Error("Present should be true for non-nil values")
    }
}

func TestSurveyData_WriteJSON_and_LoadSurveyData(t *testing.T) {
    // Minimal SurveyData for roundtrip
    sd := &SurveyData{
        Schema: Schema{
            "Q1": SchemaEntry{Key: "Q1", Text: "Question 1", QType: SC},
            "Q2": SchemaEntry{Key: "Q2", Text: "Question 2", QType: MC},
            "Q3": SchemaEntry{Key: "Q3", Text: "Question 3", QType: TE},
        },
        Responses: []Response{
            {
                "Q1": ResponseValue{val: "foo"},
                "Q2": ResponseValue{val: []string{"a", "b"}},
                "Q3": ResponseValue{val: 123},
            },
            {
                "Q1": ResponseValue{val: nil},
                "Q2": ResponseValue{val: []string{}},
                "Q3": ResponseValue{val: nil},
            },
        },
    }

    var buf bytes.Buffer
    if err := sd.WriteJSON(&buf); err != nil {
        t.Fatalf("WriteJSON failed: %v", err)
    }

    // Validate JSON is valid and roundtrips
    var out SurveyData
    if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
        t.Fatalf("json.Unmarshal failed: %v", err)
    }

    loaded, err := LoadSurveyData(bytes.NewReader(buf.Bytes()))
    if err != nil {
        t.Fatalf("LoadSurveyData failed: %v", err)
    }

    if len(loaded.Schema) != len(sd.Schema) {
        t.Errorf("Schema length mismatch: got %d, want %d", len(loaded.Schema), len(sd.Schema))
    }
    if len(loaded.Responses) != len(sd.Responses) {
        t.Errorf("Responses length mismatch: got %d, want %d", len(loaded.Responses), len(sd.Responses))
    }
}
