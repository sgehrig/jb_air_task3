package cli

// ParseCommandLine splits a command line into arguments, supporting single and double quotes, and backslash escapes.
// Example: subset "My Question" 'Someone\'s name' -> ["subset", "My Question", "Someone's name"]
func ParseCommandLine(line string) ([]string, error) {
	var args []string
	var cur []rune
	inSingle, inDouble := false, false
	escaped := false
	for _, r := range line {
		switch {
		case escaped:
			cur = append(cur, r)
			escaped = false
		case r == '\\':
			escaped = true
		case inSingle:
			if r == '\'' {
				inSingle = false
			} else {
				cur = append(cur, r)
			}
		case inDouble:
			if r == '"' {
				inDouble = false
			} else {
				cur = append(cur, r)
			}
		case r == '\'' && !inDouble:
			inSingle = true
		case r == '"' && !inSingle:
			inDouble = true
		case r == ' ' || r == '\t':
			if len(cur) > 0 {
				args = append(args, string(cur))
				cur = nil
			}
			// skip whitespace
		default:
			cur = append(cur, r)
		}
	}
	if escaped {
		return nil, ErrParseCommandLine("unfinished escape at end of input")
	}
	if inSingle || inDouble {
		return nil, ErrParseCommandLine("unclosed quote in input")
	}
	if len(cur) > 0 {
		args = append(args, string(cur))
	}
	return args, nil
}

type ErrParseCommandLine string

func (e ErrParseCommandLine) Error() string { return string(e) }
