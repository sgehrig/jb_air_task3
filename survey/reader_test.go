package reader

import (
    "bytes"
    "compress/gzip"
    "os"
    "path/filepath"
    "testing"
)

func TestSurveyData_ReadSurveyData(t *testing.T) {
    testFile := filepath.Join(".", "so_test.xlsx")
    data, err := ReadSurveyData(testFile)
    if err != nil {
        t.Fatalf("failed to read test data: %v", err)
    }

    // Basic checks
    if len(data.Schema) == 0 {
        t.Error("schema should not be empty")
    }
    if len(data.Responses) == 0 {
        t.Error("responses should not be empty")
    }

    // Check schema types
    for _, entry := range data.Schema {
        if entry.QType != SC && entry.QType != MC && entry.QType != TE {
            t.Errorf("unexpected question type for key %s: %s", entry.Key, entry.QType)
        }
    }

    // Check values for first response (if available)
    resp := data.Responses[0]
    for _, entry := range data.Schema {
        val := resp[entry.Key]
        switch entry.QType {
        case SC:
            if val.Present() {
                if _, ok := val.AsString(); !ok {
                    t.Errorf("expected string for SC key %s", entry.Key)
                }
            }
        case MC:
            if val.Present() {
                if _, ok := val.AsStringSlice(); !ok {
                    t.Errorf("expected []string for MC key %s", entry.Key)
                }
            }
        case TE:
            if val.Present() {
                if _, ok := val.AsString(); !ok {
                    t.Errorf("expected int for TE key %s", entry.Key)
                }
            }
        }
    }

    // Check NA and empty handling
    for _, resp := range data.Responses {
        for _, entry := range data.Schema {
            val := resp[entry.Key]
            if !val.Present() {
                // Should be nil for NA or empty
                if _, ok := val.AsString(); ok {
                    t.Errorf("expected nil for NA/empty, got string for key %s", entry.Key)
                }
                if _, ok := val.AsStringSlice(); ok {
                    t.Errorf("expected nil for NA/empty, got []string for key %s", entry.Key)
                }
            }
        }
    }
}

func TestSurveyData_LoadJSON(t *testing.T) {
    sd := &SurveyData{
        Schema: Schema{
            {Key: "Q1", Text: "Question 1", QType: SC},
            {Key: "Q2", Text: "Question 2", QType: MC},
            {Key: "Q3", Text: "Question 3", QType: TE},
        },
        Responses: []Response{
            {
                "Q1": {Val: "foo"},
                "Q2": {Val: []string{"a", "b"}},
                "Q3": {Val: 42},
            },
            {
                "Q1": {Val: nil},
                "Q2": {Val: nil},
                "Q3": {Val: nil},
            },
        },
    }

    var buf bytes.Buffer
    err := sd.WriteJSON(&buf)
    if err != nil {
        t.Fatalf("WriteJSON failed: %v", err)
    }

    loaded, err := LoadSurveyData(&buf)
    if err != nil {
        t.Fatalf("LoadJSON failed: %v", err)
    }

    if len(loaded.Schema) != len(sd.Schema) {
        t.Errorf("Loaded schema length mismatch: got %d, want %d", len(loaded.Schema), len(sd.Schema))
    }
    if len(loaded.Responses) != len(sd.Responses) {
        t.Errorf("Loaded responses length mismatch: got %d, want %d", len(loaded.Responses), len(sd.Responses))
    }

    // --- Gzipped JSON roundtrip ---
    var gzBuf bytes.Buffer
    gzWriter := gzip.NewWriter(&gzBuf)
    err = sd.WriteJSON(gzWriter)
    if err != nil {
        t.Fatalf("WriteJSON (gzip) failed: %v", err)
    }
    gzWriter.Close()

    gzReader, err := gzip.NewReader(&gzBuf)
    if err != nil {
        t.Fatalf("gzip.NewReader failed: %v", err)
    }
    loadedGz, err := LoadSurveyData(gzReader)
    if err != nil {
        t.Fatalf("LoadJSON (gzip) failed: %v", err)
    }
    if len(loadedGz.Schema) != len(sd.Schema) {
        t.Errorf("Loaded (gzip) schema length mismatch: got %d, want %d", len(loadedGz.Schema), len(sd.Schema))
    }
    if len(loadedGz.Responses) != len(sd.Responses) {
        t.Errorf("Loaded (gzip) responses length mismatch: got %d, want %d", len(loadedGz.Responses), len(sd.Responses))
    }
}

func TestSurveyData_WriteJSONToFile_Gzipped(t *testing.T) {
    sd := &SurveyData{
        Schema: Schema{
            {Key: "Q1", Text: "Question 1", QType: SC},
        },
        Responses: []Response{
            {"Q1": {Val: "foo"}},
        },
    }
    gzFile := "test.cache.json.gz"
    defer os.Remove(gzFile)

    // Write gzipped
    err := sd.WriteJSONToFile(gzFile)
    if err != nil {
        t.Fatalf("WriteJSONToFile (gz) failed: %v", err)
    }

    // Read back
    loaded, err := LoadSurveyDataFromFile(gzFile)
    if err != nil {
        t.Fatalf("LoadSurveyDataFromFile (gz) failed: %v", err)
    }
    if len(loaded.Schema) != len(sd.Schema) {
        t.Errorf("Loaded (gz) schema length mismatch: got %d, want %d", len(loaded.Schema), len(sd.Schema))
    }
    if len(loaded.Responses) != len(sd.Responses) {
        t.Errorf("Loaded (gz) responses length mismatch: got %d, want %d", len(loaded.Responses), len(sd.Responses))
    }
}
