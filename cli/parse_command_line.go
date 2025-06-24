package cli

import (
    "fmt"
)

// ParseCommandLine splits a command line into command and arguments,
// supporting single and double quoted strings with escaping.
// Returns error if the command line is malformed (e.g. unclosed quotes).
func ParseCommandLine(line string) (cmd string, cmdArgs []string, err error) {
    var args []string
    var current []rune
    inQuote := rune(0)
    escaped := false

    for _, r := range line {
        switch {
        case escaped:
            current = append(current, r)
            escaped = false
        case r == '\\':
            escaped = true
        case inQuote != 0:
            if r == inQuote {
                inQuote = 0
            } else {
                current = append(current, r)
            }
        case r == '\'' || r == '"':
            inQuote = r
        case r == ' ' || r == '\t':
            if len(current) > 0 {
                args = append(args, string(current))
                current = current[:0]
            }
        default:
            current = append(current, r)
        }
    }
    if escaped {
        return "", nil, fmt.Errorf("unfinished escape at end of input")
    }
    if inQuote != 0 {
        return "", nil, fmt.Errorf("unclosed quote in input")
    }
    if len(current) > 0 {
        args = append(args, string(current))
    }
    if len(args) == 0 {
        return "", nil, nil
    }
    return args[0], args[1:], nil
}
