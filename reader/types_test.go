package reader

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestSurveyData_WriteJSON(t *testing.T) {
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

	// Check that the output is valid JSON and can be unmarshaled
	var out map[string]any
	dec := json.NewDecoder(bytes.NewReader(buf.Bytes()))
	if err := dec.Decode(&out); err != nil {
		t.Errorf("output is not valid JSON: %v", err)
	}

	// Optionally, check for expected keys
	if _, ok := out["Schema"]; !ok {
		t.Error("JSON output missing 'Schema' key")
	}
	if _, ok := out["Responses"]; !ok {
		t.Error("JSON output missing 'Responses' key")
	}
}

func TestResponseValue_AsString(t *testing.T) {
	cases := []struct {
		name   string
		val    ResponseValue
		expect string
		ok     bool
	}{
		{"string", ResponseValue{value: "foo"}, "foo", true},
		{"int", ResponseValue{value: 42}, "42", true},
		{"slice", ResponseValue{value: []string{"a", "b"}}, "a;b", true},
		{"nil", ResponseValue{value: nil}, "", false},
	}
	for _, c := range cases {
		got, ok := c.val.AsString()
		if ok != c.ok || got != c.expect {
			t.Errorf("%s: got (%q, %v), want (%q, %v)", c.name, got, ok, c.expect, c.ok)
		}
	}
}

func TestResponseValue_AsStringSlice(t *testing.T) {
	cases := []struct {
		name   string
		val    ResponseValue
		expect []string
		ok     bool
	}{
		{"slice", ResponseValue{value: []string{"a", "b"}}, []string{"a", "b"}, true},
		{"string", ResponseValue{value: "foo"}, []string{"foo"}, true},
		{"int", ResponseValue{value: 42}, nil, false},
		{"nil", ResponseValue{value: nil}, nil, false},
	}
	for _, c := range cases {
		got, ok := c.val.AsStringSlice()
		if ok != c.ok || (ok && !equalStringSlices(got, c.expect)) {
			t.Errorf("%s: got (%v, %v), want (%v, %v)", c.name, got, ok, c.expect, c.ok)
		}
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestResponseValue_IsPresent(t *testing.T) {
	if !((ResponseValue{value: "foo"}).IsPresent()) {
		t.Error("expected IsPresent true for non-nil value")
	}
	if (ResponseValue{value: nil}).IsPresent() {
		t.Error("expected IsPresent false for nil value")
	}
}

func TestSurveyData_WriteAndLoadJSON_Gzip(t *testing.T) {
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

	f, err := os.CreateTemp("", "test-survey-*.json.gz")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	f.Close()

	err = sd.WriteJSONToFile(f.Name())
	if err != nil {
		t.Fatalf("WriteJSONToFile (gzip) failed: %v", err)
	}

	loaded, err := LoadJSONFromFile(f.Name())
	if err != nil {
		t.Fatalf("LoadJSONFromFile (gzip) failed: %v", err)
	}

	if len(loaded.Schema) != len(sd.Schema) {
		t.Errorf("Loaded schema length mismatch: got %d, want %d", len(loaded.Schema), len(sd.Schema))
	}
	if len(loaded.Responses) != len(sd.Responses) {
		t.Errorf("Loaded responses length mismatch: got %d, want %d", len(loaded.Responses), len(sd.Responses))
	}
}
