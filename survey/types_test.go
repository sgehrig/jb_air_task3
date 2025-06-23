package survey

import (
    "bytes"
    "encoding/json"
    "testing"
)

func TestResponseValue_AsMethods(t *testing.T) {
    rvInt := ResponseValue{Val: "42"}
    rvStr := ResponseValue{Val: "foo"}
    rvSlice := ResponseValue{Val: []string{"a", "b"}}
    rvNil := ResponseValue{Val: nil}

    if v, ok := rvInt.AsString(); !ok || v != "42" {
        t.Errorf("AsString failed: got %v, %v", v, ok)
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
            {Key: "Q1", Text: "Question 1", QType: SC, UsedOptions: []string{"foo"}},
            {Key: "Q2", Text: "Question 2", QType: MC, UsedOptions: []string{"a", "b"}},
            {Key: "Q3", Text: "Question 3", QType: TE, UsedOptions: []string{}},
        },
        Responses: []Response{
            {
                "Q1": {Val: "foo"},
                "Q2": {Val: []string{"a", "b"}},
                "Q3": {Val: 123},
            },
            {
                "Q1": {Val: nil},
                "Q2": {Val: []string{}},
                "Q3": {Val: nil},
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

func TestSchema_SearchForString(t *testing.T) {
    schema := Schema{
        &SchemaEntry{Key: "Q1", Text: "Favorite color", QType: SC, UsedOptions: []string{"red", "blue"}},
        &SchemaEntry{Key: "Q2", Text: "Programming languages", QType: MC, UsedOptions: []string{"Go", "Python"}},
        &SchemaEntry{Key: "Q3", Text: "Age", QType: TE},
    }

    tests := []struct {
        query    string
        wantKeys []string
    }{
        {"color", []string{"Q1"}},
        {"go", []string{"Q2"}},
        {"age", []string{"Q2", "Q3"}},
        {"blue", []string{"Q1"}},
        {"python", []string{"Q2"}},
        {"q1", []string{"Q1"}},
        {"q", []string{"Q1", "Q2", "Q3"}},
        {"notfound", []string{}},
    }

    for _, tt := range tests {
        got := schema.SearchForString(tt.query)
        if len(got) != len(tt.wantKeys) {
            t.Errorf("SearchForString(%q) got %d results, want %d", tt.query, len(got), len(tt.wantKeys))
            continue
        }
        gotKeys := make(map[string]bool)
        for _, entry := range got {
            gotKeys[entry.Key] = true
        }
        for _, wantKey := range tt.wantKeys {
            if !gotKeys[wantKey] {
                t.Errorf("SearchForString(%q) missing key %q", tt.query, wantKey)
            }
        }
    }
}

func TestSurveyData_CreateSubset(t *testing.T) {
    schema := Schema{
        &SchemaEntry{Key: "Q1", Text: "Favorite color", QType: SC, UsedOptions: []string{"red", "blue"}},
        &SchemaEntry{Key: "Q2", Text: "Languages", QType: MC, UsedOptions: []string{"Go", "Python"}},
        &SchemaEntry{Key: "Q3", Text: "Comment", QType: TE},
    }
    responses := []Response{
        {"Q1": ResponseValue{Val: "red"}, "Q2": ResponseValue{Val: []string{"Go", "Python"}}, "Q3": ResponseValue{Val: "Nice!"}},
        {"Q1": ResponseValue{Val: "blue"}, "Q2": ResponseValue{Val: []string{"Go"}}, "Q3": ResponseValue{Val: "Cool!"}},
        {"Q1": ResponseValue{Val: "green"}, "Q2": ResponseValue{Val: []string{"Java"}}, "Q3": ResponseValue{Val: "Okay!"}},
        {"Q1": ResponseValue{Val: "red"}, "Q2": ResponseValue{Val: []string{"Python"}}, "Q3": ResponseValue{Val: "Great!"}},
    }
    sd := &SurveyData{
        Schema:    schema,
        Responses: responses,
    }

    // SC: search for "red" in Q1
    subset := sd.CreateSubset("Q1", "red")
    if len(subset) != 2 {
        t.Errorf("CreateSubset SC: got %d, want 2", len(subset))
    }

    // MC: search for "python" in Q2 (case-insensitive, partial)
    subset = sd.CreateSubset("Q2", "pyth")
    if len(subset) != 2 {
        t.Errorf("CreateSubset MC: got %d, want 2", len(subset))
    }

    // MC: search for "go" in Q2
    subset = sd.CreateSubset("Q2", "go")
    if len(subset) != 2 {
        t.Errorf("CreateSubset MC: got %d, want 2", len(subset))
    }

    // MC: search for "java" in Q2
    subset = sd.CreateSubset("Q2", "java")
    if len(subset) != 1 {
        t.Errorf("CreateSubset MC: got %d, want 1", len(subset))
    }

    // Invalid key
    subset = sd.CreateSubset("QX", "red")
    if len(subset) != 0 {
        t.Errorf("CreateSubset invalid key: got %d, want 0", len(subset))
    }

    // TE: should not match anything
    subset = sd.CreateSubset("Q3", "Nice")
    if len(subset) != 0 {
        t.Errorf("CreateSubset TE: got %d, want 0", len(subset))
    }
}
