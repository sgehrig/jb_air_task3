package survey

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ResponseQuery struct {
	Keys  []string
	Range RangeSelector
}

func (rq *ResponseQuery) Limit(responses []SurveyResponse) []SurveyResponse {
	n := len(responses)
	startIndex := 0
	endIndex := n - 1

	switch rq.Range.Start.Type {
	case "first":
		startIndex = 0 + rq.Range.Start.Offset
	case "last":
		startIndex = (n - 1) + rq.Range.Start.Offset
	case "index":
		startIndex = rq.Range.Start.Offset
	}
	switch rq.Range.End.Type {
	case "first":
		endIndex = 0 + rq.Range.End.Offset
	case "last":
		endIndex = (n - 1) + rq.Range.End.Offset
	case "index":
		endIndex = rq.Range.End.Offset
	}
	// Clamp indices
	if startIndex < 0 {
		startIndex = 0
	}
	if endIndex >= n {
		endIndex = n - 1
	}
	if endIndex < startIndex || startIndex >= n {
		return []SurveyResponse{}
	}
	return responses[startIndex : endIndex+1]
}

func AllResponseQuery() *ResponseQuery {
	return &ResponseQuery{
		Keys: nil,
		Range: RangeSelector{
			Start: RangeEndpoint{Type: "first", Offset: 0, RawString: "first"},
			End:   RangeEndpoint{Type: "last", Offset: 0, RawString: "last"},
		},
	}
}

type RangeSelector struct {
	Start RangeEndpoint
	End   RangeEndpoint
}

type RangeEndpoint struct {
	Type      string // "first", "last", "index"
	Offset    int    // applies to first/last, e.g., +2 or -3
	RawString string // original representation for debugging
}

func ParseResponseQuery(input string) (*ResponseQuery, error) {
	input = strings.TrimSpace(input)
	sections := splitSections(input)
	q := &ResponseQuery{Keys: nil, Range: RangeSelector{
		Start: RangeEndpoint{Type: "first", Offset: 0, RawString: "first"},
		End:   RangeEndpoint{Type: "last", Offset: 0, RawString: "last"},
	}}
	for _, sec := range sections {
		sec = strings.TrimSpace(sec)
		if strings.HasPrefix(sec, "keys:") || strings.HasPrefix(sec, "keys=") {
			keysStr := strings.TrimSpace(sec[strings.Index(sec, ":")+1:])
			if eq := strings.Index(keysStr, "="); eq != -1 {
				keysStr = keysStr[eq+1:]
			}
			for _, k := range splitKeys(keysStr) {
				k = strings.TrimSpace(k)
				if len(k) >= 2 && ((k[0] == '\'' && k[len(k)-1] == '\'') || (k[0] == '"' && k[len(k)-1] == '"')) {
					quoteChar := k[0]
					k = k[1 : len(k)-1]
					k = strings.ReplaceAll(k, `\\`, `\`)
					switch quoteChar {
					case '\'':
						k = strings.ReplaceAll(k, `''`, `'`)
					case '"':
						k = strings.ReplaceAll(k, `\"`, `"`)
						k = strings.ReplaceAll(k, `""`, `"`)
					}
				}
				if k != "" {
					q.Keys = append(q.Keys, k)
				}
			}
			// No error if keys is empty
		} else if strings.HasPrefix(sec, "range:") || strings.HasPrefix(sec, "range=") {
			rangeStr := strings.TrimSpace(sec[strings.Index(sec, ":")+1:])
			if eq := strings.Index(rangeStr, "="); eq != -1 {
				rangeStr = rangeStr[eq+1:]
			}
			if !strings.HasPrefix(rangeStr, "[") || !strings.HasSuffix(rangeStr, "]") {
				return nil, fmt.Errorf("range must be enclosed in [..]: got %q", rangeStr)
			}
			inner := strings.TrimSpace(rangeStr[1 : len(rangeStr)-1])
			parts := strings.Split(inner, "..")
			if len(parts) != 2 {
				return nil, fmt.Errorf("range must be of form [start..end]")
			}
			start, err := parseRangeEndpoint(parts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid range start: %w", err)
			}
			end, err := parseRangeEndpoint(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid range end: %w", err)
			}
			q.Range = RangeSelector{Start: start, End: end}
		}
	}
	return q, nil
}

// splitKeys splits a keys string on commas not inside single or double quotes.
func splitKeys(s string) []string {
	var res []string
	var cur strings.Builder
	inSingle, inDouble := false, false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\'' && !inDouble {
			inSingle = !inSingle
			cur.WriteByte(c)
		} else if c == '"' && !inSingle {
			inDouble = !inDouble
			cur.WriteByte(c)
		} else if c == ',' && !inSingle && !inDouble {
			res = append(res, cur.String())
			cur.Reset()
		} else {
			cur.WriteByte(c)
		}
	}
	if cur.Len() > 0 {
		res = append(res, cur.String())
	}
	return res
}

var reEndpoint = regexp.MustCompile(`^(first|last|[0-9]+)([+-][0-9]+)?$`)

func parseRangeEndpoint(s string) (RangeEndpoint, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return RangeEndpoint{}, fmt.Errorf("empty endpoint")
	}
	m := reEndpoint.FindStringSubmatch(s)
	if m == nil {
		return RangeEndpoint{RawString: s}, fmt.Errorf("invalid endpoint: %q", s)
	}
	var typ string
	var offset int
	if m[1] == "first" || m[1] == "last" {
		typ = m[1]
		if m[2] != "" {
			offset, _ = strconv.Atoi(m[2])
		}
	} else {
		typ = "index"
		parsed, _ := strconv.Atoi(m[1])
		offset = parsed
		if m[2] != "" {
			// For index+N or index-N, offset is not used, but we could error or ignore
		}
	}
	return RangeEndpoint{Type: typ, Offset: offset, RawString: s}, nil
}

func splitSections(input string) []string {
	// Split on both semicolons and newlines
	var out []string
	for _, part := range strings.Split(input, "\n") {
		for _, sub := range strings.Split(part, ";") {
			trimmed := strings.TrimSpace(sub)
			if trimmed != "" {
				out = append(out, trimmed)
			}
		}
	}
	return out
}
