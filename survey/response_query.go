package survey

import (
    "errors"
    "fmt"
    "regexp"
    "strconv"
    "strings"
)

type ResponseQuery struct {
    Keys  []string
    Range RangeSelector
}

func (rq *ResponseQuery) Limit(responses []Response) []Response {
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
        return []Response{}
    }
    return responses[startIndex : endIndex+1]
}

func AllResponseQuery() *ResponseQuery {
    return &ResponseQuery{
        Keys: []string{},
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
    if input == "" {
        // Use default first..last range selector for empty query, and empty keys slice
        return AllResponseQuery(), nil
    }

    // Support both single-line and multi-line variants
    sections := splitSections(input)

    var keysSection, rangeSection string
    for _, sec := range sections {
        sec = strings.TrimSpace(sec)
        if strings.HasPrefix(sec, "keys:") || strings.HasPrefix(sec, "keys=") {
            keysSection = sec
        } else if strings.HasPrefix(sec, "range:") || strings.HasPrefix(sec, "range=") {
            rangeSection = sec
        }
    }

    var keys []string
    var err error
    if keysSection != "" {
        keys, err = parseKeys(keysSection)
        if err != nil {
            return nil, fmt.Errorf("keys: %w", err)
        }
    } else {
        keys = []string{}
    }

    var rng *RangeSelector
    if rangeSection != "" {
        rng, err = parseRange(rangeSection)
        if err != nil {
            return nil, fmt.Errorf("range: %w", err)
        }
    }

    // If no range specified, use default first..last
    if rng == nil {
        rng = &RangeSelector{
            Start: RangeEndpoint{Type: "first", Offset: 0, RawString: "first"},
            End:   RangeEndpoint{Type: "last", Offset: 0, RawString: "last"},
        }
    }

    return &ResponseQuery{Keys: keys, Range: *rng}, nil
}

func splitSections(input string) []string {
    // Accept both ; and newlines as section separators
    parts := strings.Split(input, ";")
    if len(parts) == 1 {
        parts = strings.Split(input, "\n")
    }
    return parts
}

func parseKeys(sec string) ([]string, error) {
    sep := ":"
    if strings.Contains(sec, "=") {
        sep = "="
    }
    parts := strings.SplitN(sec, sep, 2)
    if len(parts) != 2 {
        return nil, errors.New("invalid keys section")
    }
    raw := strings.TrimSpace(parts[1])
    if raw == "" {
        return nil, errors.New("no keys specified")
    }

    // Custom parser for quoted and unquoted keys
    var out []string
    i := 0
    for i < len(raw) {
        // Skip whitespace and commas
        for i < len(raw) && (raw[i] == ' ' || raw[i] == ',') {
            i++
        }
        if i >= len(raw) {
            break
        }
        if raw[i] == '\'' || raw[i] == '"' {
            quote := raw[i]
            i++
            var sb strings.Builder
            for i < len(raw) {
                if raw[i] == quote {
                    // Check for escaped quote ('' or "")
                    if i+1 < len(raw) && raw[i+1] == quote {
                        sb.WriteByte(quote)
                        i += 2
                        continue
                    }
                    i++
                    break
                }
                // For double quotes, allow \" as escape
                if quote == '"' && raw[i] == '\\' && i+1 < len(raw) && raw[i+1] == '"' {
                    sb.WriteByte('"')
                    i += 2
                    continue
                }
                sb.WriteByte(raw[i])
                i++
            }
            out = append(out, sb.String())
        } else {
            // Unquoted key: read until comma
            start := i
            for i < len(raw) && raw[i] != ',' {
                i++
            }
            key := strings.TrimSpace(raw[start:i])
            if key != "" {
                out = append(out, key)
            }
        }
        // Skip trailing comma if present
        if i < len(raw) && raw[i] == ',' {
            i++
        }
    }

    if len(out) == 0 {
        return nil, errors.New("no keys specified")
    }
    return out, nil
}

func parseRange(sec string) (*RangeSelector, error) {
    sep := ":"
    if strings.Contains(sec, "=") {
        sep = "="
    }
    parts := strings.SplitN(sec, sep, 2)
    if len(parts) != 2 {
        return nil, errors.New("invalid range section")
    }
    raw := strings.TrimSpace(parts[1])
    if !strings.HasPrefix(raw, "[") || !strings.HasSuffix(raw, "]") {
        return nil, errors.New("range must be enclosed in [..]")
    }
    raw = strings.TrimPrefix(raw, "[")
    raw = strings.TrimSuffix(raw, "]")
    rangeParts := strings.SplitN(raw, "..", 2)
    if len(rangeParts) != 2 {
        return nil, errors.New("range must use .. to separate start and end")
    }
    start, err := parseRangeEndpoint(strings.TrimSpace(rangeParts[0]))
    if err != nil {
        return nil, fmt.Errorf("invalid start: %w", err)
    }
    end, err := parseRangeEndpoint(strings.TrimSpace(rangeParts[1]))
    if err != nil {
        return nil, fmt.Errorf("invalid end: %w", err)
    }
    return &RangeSelector{Start: start, End: end}, nil
}

var endpointRe = regexp.MustCompile(`^(first|last)([+-]\d+)?$|^(\d+)$`)

func parseRangeEndpoint(s string) (RangeEndpoint, error) {
    s = strings.TrimSpace(s)
    if s == "" {
        return RangeEndpoint{}, errors.New("empty endpoint")
    }
    m := endpointRe.FindStringSubmatch(s)
    if m == nil {
        return RangeEndpoint{RawString: s}, errors.New("must be first, last, or integer (with optional +N/-N for first/last)")
    }
    if m[1] != "" { // first or last
        offset := 0
        if m[2] != "" {
            var err error
            offset, err = strconv.Atoi(m[2])
            if err != nil {
                return RangeEndpoint{RawString: s}, fmt.Errorf("invalid offset: %w", err)
            }
        }
        return RangeEndpoint{
            Type:      m[1],
            Offset:    offset,
            RawString: s,
        }, nil
    }
    if m[3] != "" { // integer index
        idx, err := strconv.Atoi(m[3])
        if err != nil {
            return RangeEndpoint{RawString: s}, fmt.Errorf("invalid index: %w", err)
        }
        return RangeEndpoint{
            Type:      "index",
            Offset:    idx,
            RawString: s,
        }, nil
    }
    return RangeEndpoint{RawString: s}, errors.New("unrecognized endpoint")
}
