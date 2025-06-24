package survey

import (
	"reflect"
	"testing"
)

func TestParseQuery_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  ResponseQuery
	}{
		{
			"single-line full",
			"keys=name,email;range=[first+1..last-2]",
			ResponseQuery{
				Keys: []string{"name", "email"},
				Range: RangeSelector{
					Start: RangeEndpoint{Type: "first", Offset: 1, RawString: "first+1"},
					End:   RangeEndpoint{Type: "last", Offset: -2, RawString: "last-2"},
				},
			},
		},
		{
			"multi-line full",
			"keys: id \n range: [0..5]",
			ResponseQuery{
				Keys: []string{"id"},
				Range: RangeSelector{
					Start: RangeEndpoint{Type: "index", Offset: 0, RawString: "0"},
					End:   RangeEndpoint{Type: "index", Offset: 5, RawString: "5"},
				},
			},
		},
		{
			"keys only",
			"range:[first..last]",
			ResponseQuery{
				Keys: nil,
				Range: RangeSelector{
					Start: RangeEndpoint{Type: "first", Offset: 0, RawString: "first"},
					End:   RangeEndpoint{Type: "last", Offset: 0, RawString: "last"},
				},
			},
		},
		{
			"empty query",
			"",
			ResponseQuery{
				Keys: nil,
				Range: RangeSelector{
					Start: RangeEndpoint{Type: "first", Offset: 0, RawString: "first"},
					End:   RangeEndpoint{Type: "last", Offset: 0, RawString: "last"},
				},
			},
		},
		{
			"keys only",
			"keys:x,y,z",
			ResponseQuery{
				Keys: []string{"x", "y", "z"},
				Range: RangeSelector{
					Start: RangeEndpoint{Type: "first", Offset: 0, RawString: "first"},
					End:   RangeEndpoint{Type: "last", Offset: 0, RawString: "last"},
				},
			},
		},
		{
			"keys with spaces",
			"keys:foo bar, baz qux",
			ResponseQuery{
				Keys: []string{"foo bar", "baz qux"},
				Range: RangeSelector{
					Start: RangeEndpoint{Type: "first", Offset: 0, RawString: "first"},
					End:   RangeEndpoint{Type: "last", Offset: 0, RawString: "last"},
				},
			},
		},
		{
			"single and double quoted keys",
			`keys:'foo,bar', "baz qux", plain, 'with ''quote''', "with \"quote\""`,
			ResponseQuery{
				Keys: []string{"foo,bar", "baz qux", "plain", "with 'quote'", "with \"quote\""},
				Range: RangeSelector{
					Start: RangeEndpoint{Type: "first", Offset: 0, RawString: "first"},
					End:   RangeEndpoint{Type: "last", Offset: 0, RawString: "last"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseResponseQuery(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got.Keys, tt.want.Keys) {
				t.Errorf("ParseQuery() Keys = %+v, want %+v", got.Keys, tt.want.Keys)
			}
			if !reflect.DeepEqual(got.Range, tt.want.Range) {
				t.Errorf("ParseQuery() Range = %+v, want %+v", got.Range, tt.want.Range)
			}
		})
	}
}

func TestParseQuery_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"invalid range endpoints", "keys:foo;range:[foo..bar]"},
		{"empty range", "keys:x;range:[..]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseResponseQuery(tt.input)
			if err == nil {
				t.Errorf("expected error for input %q", tt.input)
			}
		})
	}
}

func TestResponseQuery_Limit(t *testing.T) {
	// Helper to make SurveyResponse with just an "id" field for easy comparison
	makeResp := func(id int) SurveyResponse {
		return SurveyResponse{"id": ResponseValue{Value: id}}
	}
	responses := []SurveyResponse{
		makeResp(0), makeResp(1), makeResp(2), makeResp(3), makeResp(4), makeResp(5),
	}
	tests := []struct {
		name  string
		query ResponseQuery
		want  []SurveyResponse
	}{
		{
			"full range",
			ResponseQuery{
				Range: RangeSelector{
					Start: RangeEndpoint{Type: "first", Offset: 0},
					End:   RangeEndpoint{Type: "last", Offset: 0},
				},
			},
			[]SurveyResponse{makeResp(0), makeResp(1), makeResp(2), makeResp(3), makeResp(4), makeResp(5)},
		},
		{
			"first 3",
			ResponseQuery{
				Range: RangeSelector{
					Start: RangeEndpoint{Type: "first", Offset: 0},
					End:   RangeEndpoint{Type: "index", Offset: 2},
				},
			},
			[]SurveyResponse{makeResp(0), makeResp(1), makeResp(2)},
		},
		{
			"last 2",
			ResponseQuery{
				Range: RangeSelector{
					Start: RangeEndpoint{Type: "index", Offset: 4},
					End:   RangeEndpoint{Type: "last", Offset: 0},
				},
			},
			[]SurveyResponse{makeResp(4), makeResp(5)},
		},
		{
			"middle slice",
			ResponseQuery{
				Range: RangeSelector{
					Start: RangeEndpoint{Type: "index", Offset: 2},
					End:   RangeEndpoint{Type: "index", Offset: 4},
				},
			},
			[]SurveyResponse{makeResp(2), makeResp(3), makeResp(4)},
		},
		{
			"last-2 to last",
			ResponseQuery{
				Range: RangeSelector{
					Start: RangeEndpoint{Type: "last", Offset: -2},
					End:   RangeEndpoint{Type: "last", Offset: 0},
				},
			},
			[]SurveyResponse{makeResp(3), makeResp(4), makeResp(5)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.query.Limit(responses)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Limit() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
