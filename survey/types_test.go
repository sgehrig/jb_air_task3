package survey

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestSurveyData_WriteJSON(t *testing.T) {
	sd := &SurveyData{
		Schema: Schema{
			{Key: "Q1", Text: "Question 1", QType: SC},
			{Key: "Q2", Text: "Question 2", QType: MC},
			{Key: "Q3", Text: "Question 3", QType: TE},
		},
		Responses: []SurveyResponse{
			{
				"Q1": {Value: "foo"},
				"Q2": {Value: []string{"a", "b"}},
				"Q3": {Value: 42},
			},
			{
				"Q1": {Value: nil},
				"Q2": {Value: nil},
				"Q3": {Value: nil},
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
		{"string", ResponseValue{Value: "foo"}, "foo", true},
		{"int", ResponseValue{Value: 42}, "42", true},
		{"slice", ResponseValue{Value: []string{"a", "b"}}, "a;b", true},
		{"nil", ResponseValue{Value: nil}, "", false},
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
		{"slice", ResponseValue{Value: []string{"a", "b"}}, []string{"a", "b"}, true},
		{"string", ResponseValue{Value: "foo"}, []string{"foo"}, true},
		{"int", ResponseValue{Value: 42}, nil, false},
		{"nil", ResponseValue{Value: nil}, nil, false},
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
	if !((ResponseValue{Value: "foo"}).IsPresent()) {
		t.Error("expected IsPresent true for non-nil value")
	}
	if (ResponseValue{Value: nil}).IsPresent() {
		t.Error("expected IsPresent false for nil value")
	}
}

func TestSurveyData_WriteAndLoadJSON_Gzip(t *testing.T) {
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

func TestSchema_SearchForString(t *testing.T) {
	s := Schema{
		&SchemaEntry{Key: "Q1", Text: "Favorite color", QType: SC, Options: []string{"red", "blue", "green"}},
		&SchemaEntry{Key: "Q2", Text: "Hobbies", QType: MC, Options: []string{"reading", "sports", "music"}},
		&SchemaEntry{Key: "Q3", Text: "Age", QType: TE},
	}
	tests := []struct {
		query  string
		expect []string
	}{
		{"color", []string{"Q1"}},
		{"sports", []string{"Q2"}},
		{"q3", []string{"Q3"}},
		{"red", []string{"Q1"}},
		{"music", []string{"Q2"}},
		{"notfound", nil},
	}
	for _, tc := range tests {
		results := s.SearchForString(tc.query)
		var keys []string
		for _, entry := range results {
			keys = append(keys, entry.Key)
		}
		if !equalStringSlices(keys, tc.expect) {
			t.Errorf("query %q: got %v, want %v", tc.query, keys, tc.expect)
		}
	}
}

func TestSurveyData_CreateSubset(t *testing.T) {
	schema := Schema{
		&SchemaEntry{Key: "Q1", Text: "Favorite color", QType: SC, Options: []string{"red", "blue", "green"}},
		&SchemaEntry{Key: "Q2", Text: "Hobbies", QType: MC, Options: []string{"reading", "sports", "music"}},
		&SchemaEntry{Key: "Q3", Text: "Age", QType: TE},
	}
	responses := []SurveyResponse{
		{"Q1": ResponseValue{Value: "red"}, "Q2": ResponseValue{Value: []string{"reading", "music"}}, "Q3": ResponseValue{Value: "30"}},
		{"Q1": ResponseValue{Value: "blue"}, "Q2": ResponseValue{Value: []string{"sports"}}, "Q3": ResponseValue{Value: "25"}},
		{"Q1": ResponseValue{Value: "green"}, "Q2": ResponseValue{Value: []string{"reading"}}, "Q3": ResponseValue{Value: "40"}},
		{"Q1": ResponseValue{Value: "red"}, "Q2": ResponseValue{Value: []string{"music"}}, "Q3": ResponseValue{Value: "22"}},
	}
	sd := &SurveyData{Schema: schema, Responses: responses}
	tests := []struct {
		key     string
		query   string
		expectN int
	}{
		{"Q1", "red", 2},
		{"Q1", "blue", 1},
		{"Q2", "music", 2},
		{"Q2", "reading", 2},
		{"Q2", "sports", 1},
		{"Q3", "22", 1}, // TE, but will match as string
		{"Q1", "notfound", 0},
		{"QX", "red", 0}, // invalid key
	}
	for _, tc := range tests {
		result := sd.CreateSubset(tc.key, tc.query)
		if len(result) != tc.expectN {
			t.Errorf("CreateSubset(%q, %q): got %d, want %d", tc.key, tc.query, len(result), tc.expectN)
		}
	}
}
