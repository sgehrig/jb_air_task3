package reader

import (
	"bytes"
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
	for key, entry := range data.Schema {
		if entry.QType != SC && entry.QType != MC && entry.QType != TE {
			t.Errorf("unexpected question type for key %s: %s", key, entry.QType)
		}
	}

	// Check values for first response (if available)
	resp := data.Responses[0]
	for key, entry := range data.Schema {
		val := resp[key]
		switch entry.QType {
		case SC:
			if val.IsPresent() {
				if _, ok := val.AsString(); !ok {
					t.Errorf("expected string for SC key %s", key)
				}
			}
		case MC:
			if val.IsPresent() {
				if _, ok := val.AsStringSlice(); !ok {
					t.Errorf("expected []string for MC key %s", key)
				}
			}
		case TE:
			if val.IsPresent() {
				if _, ok := val.AsInt(); !ok {
					t.Errorf("expected int for TE key %s", key)
				}
			}
		}
	}

	// Check NA and empty handling
	for _, resp := range data.Responses {
		for key := range data.Schema {
			val := resp[key]
			if !val.IsPresent() {
				// Should be nil for NA or empty
				if _, ok := val.AsString(); ok {
					t.Errorf("expected nil for NA/empty, got string for key %s", key)
				}
				if _, ok := val.AsInt(); ok {
					t.Errorf("expected nil for NA/empty, got int for key %s", key)
				}
				if _, ok := val.AsStringSlice(); ok {
					t.Errorf("expected nil for NA/empty, got []string for key %s", key)
				}
			}
		}
	}
	// Add more checks as needed for your test data
}

func TestSurveyData_LoadJSON(t *testing.T) {
	sd := &SurveyData{
		Schema: Schema{
			"Q1": {Key: "Q1", Text: "Question 1", QType: SC},
			"Q2": {Key: "Q2", Text: "Question 2", QType: MC},
			"Q3": {Key: "Q3", Text: "Question 3", QType: TE},
		},
		Responses: []SurveyResponse{
			{
				"Q1": ResponseValue{value: "foo"},
				"Q2": ResponseValue{value: []string{"a", "b"}},
				"Q3": ResponseValue{value: 42},
			},
			{
				"Q1": ResponseValue{value: nil},
				"Q2": ResponseValue{value: nil},
				"Q3": ResponseValue{value: nil},
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
