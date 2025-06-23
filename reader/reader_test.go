package reader

import (
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
				if _, err := val.AsString(); err != nil {
					t.Errorf("expected string for SC key %s", key)
				}
			}
		case MC:
			if val.IsPresent() {
				if _, err := val.AsStringSlice(); err != nil {
					t.Errorf("expected []string for MC key %s", key)
				}
			}
		case TE:
			if val.IsPresent() {
				if _, err := val.AsInt(); err != nil {
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
				if _, err := val.AsString(); err == nil {
					t.Errorf("expected nil for NA/empty, got string for key %s", key)
				}
				if _, err := val.AsInt(); err == nil {
					t.Errorf("expected nil for NA/empty, got int for key %s", key)
				}
				if _, err := val.AsStringSlice(); err == nil {
					t.Errorf("expected nil for NA/empty, got []string for key %s", key)
				}
			}
		}
	}
	// Add more checks as needed for your test data
}
