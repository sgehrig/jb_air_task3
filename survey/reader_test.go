package survey

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestReadSurvey(t *testing.T) {
	testFile := filepath.Join(".", "so_test.xlsx")
	data, err := ReadSurvey(testFile)
	if err != nil {
		t.Fatalf("failed to read survey: %v", err)
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
			if val.IsPresent() {
				if _, ok := val.AsString(); !ok {
					t.Errorf("expected string for SC key %s", entry.Key)
				}
			}
		case MC:
			if val.IsPresent() {
				if _, ok := val.AsStringSlice(); !ok {
					t.Errorf("expected []string for MC key %s", entry.Key)
				}
			}
		case TE:
			if val.IsPresent() {
				if s, ok := val.AsString(); !ok || s == "" {
					t.Errorf("expected string for TE key %s", entry.Key)
				}
			}
		}
	}

	// Check NA and empty handling
	for _, resp := range data.Responses {
		for _, entry := range data.Schema {
			val := resp[entry.Key]
			if !val.IsPresent() {
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
	// Add more checks as needed for your test data
}

func TestSurveyData_LoadJSON(t *testing.T) {
	sd := &SurveyData{
		Schema: Schema{
			{Key: "Q1", Text: "Question 1", QType: SC},
			{Key: "Q2", Text: "Question 2", QType: MC},
			{Key: "Q3", Text: "Question 3", QType: TE},
		},
		Responses: []SurveyResponse{
			{
				"Q1": ResponseValue{Value: "foo"},
				"Q2": ResponseValue{Value: []string{"a", "b"}},
				"Q3": ResponseValue{Value: 42},
			},
			{
				"Q1": ResponseValue{Value: nil},
				"Q2": ResponseValue{Value: nil},
				"Q3": ResponseValue{Value: nil},
			},
		},
	}

	var buf bytes.Buffer
	err := sd.WriteJSON(&buf)
	if err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	loaded, err := LoadJSON(&buf)
	if err != nil {
		t.Fatalf("LoadJSON failed: %v", err)
	}

	if len(loaded.Schema) != len(sd.Schema) {
		t.Errorf("Loaded schema length mismatch: got %d, want %d", len(loaded.Schema), len(sd.Schema))
	}
	if len(loaded.Responses) != len(sd.Responses) {
		t.Errorf("Loaded responses length mismatch: got %d, want %d", len(loaded.Responses), len(sd.Responses))
	}
}

func TestReadSurveyCached(t *testing.T) {
	xlsxFile := "so_test.xlsx"
	cacheFile := createCacheFilename(xlsxFile)

	// Ensure cleanup before and after
	_ = os.Remove(cacheFile)
	defer os.Remove(cacheFile)

	// 1. Cache does not exist, excel is valid
	data1, err := ReadSurveyCached(xlsxFile)
	if err != nil {
		t.Fatalf("expected to read from xlsx, got error: %v", err)
	}
	if data1 == nil || len(data1.Schema) == 0 {
		t.Fatalf("expected valid data from xlsx")
	}
	if _, err := os.Stat(cacheFile); err != nil {
		t.Errorf("expected cache file to be created")
	}

	// 2. Cache exists and is valid
	data2, err := ReadSurveyCached(xlsxFile)
	if err != nil {
		t.Fatalf("expected to read from cache, got error: %v", err)
	}
	if data2 == nil || len(data2.Schema) == 0 {
		t.Fatalf("expected valid data from cache")
	}

	// 3. Cache exists but is invalid
	// Overwrite cache file with invalid JSON
	if err := os.WriteFile(cacheFile, []byte("{invalid json"), 0644); err != nil {
		t.Fatalf("failed to write invalid cache: %v", err)
	}
	data3, err := ReadSurveyCached(xlsxFile)
	if err != nil {
		t.Fatalf("expected fallback to xlsx, got error: %v", err)
	}
	if data3 == nil || len(data3.Schema) == 0 {
		t.Fatalf("expected valid data after fallback from invalid cache")
	}
}
