package reader

import (
	"bytes"
	"encoding/json"
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
