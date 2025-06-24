package cli

import (
    "reflect"
    "testing"
)

func TestParseCommandLine(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        wantCmd  string
        wantArgs []string
        wantErr  bool
    }{
        {
            name:     "simple",
            input:    "subset foo bar",
            wantCmd:  "subset",
            wantArgs: []string{"foo", "bar"},
        },
        {
            name:     "double quoted arg",
            input:    `subset "My Question"`,
            wantCmd:  "subset",
            wantArgs: []string{"My Question"},
        },
        {
            name:     "single quoted arg",
            input:    "subset 'My Question'",
            wantCmd:  "subset",
            wantArgs: []string{"My Question"},
        },
        {
            name:     "escaped single quote in single quotes",
            input:    `subset 'Someone\'s name'`,
            wantCmd:  "subset",
            wantArgs: []string{"Someone's name"},
        },
        {
            name:     "escaped double quote in double quotes",
            input:    `subset "A \"quote\" here"`,
            wantCmd:  "subset",
            wantArgs: []string{`A "quote" here`},
        },
        {
            name:     "mixed quotes",
            input:    `subset "Test 'ass'"`,
            wantCmd:  "subset",
            wantArgs: []string{`Test 'ass'`},
        },
        {
            name:     "mixed quotes 2",
            input:    `subset 'Test "ABC"'`,
            wantCmd:  "subset",
            wantArgs: []string{`Test "ABC"`},
        },
        {
            name:     "multiple quoted args",
            input:    `subset "A B" 'C D'`,
            wantCmd:  "subset",
            wantArgs: []string{"A B", "C D"},
        },
        {
            name:     "empty input",
            input:    "",
            wantCmd:  "",
            wantArgs: nil,
        },
        {
            name:     "only spaces",
            input:    "   ",
            wantCmd:  "",
            wantArgs: nil,
        },
        {
            name:    "unclosed single quote",
            input:   "subset 'foo",
            wantErr: true,
        },
        {
            name:    "unclosed double quote",
            input:   `subset "foo`,
            wantErr: true,
        },
        {
            name:     "illegal nesting",
            input:    `subset "foo 'bar"`,
            wantCmd:  "subset",
            wantArgs: []string{`foo 'bar`},
        },
        {
            name:     "escaped backslash",
            input:    `subset foo\\bar`,
            wantCmd:  "subset",
            wantArgs: []string{`foo\bar`},
        },
        {
            name:     "escaped quote outside quotes",
            input:    `subset foo\"bar`,
            wantCmd:  "subset",
            wantArgs: []string{`foo"bar`},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cmd, args, err := ParseCommandLine(tt.input)
            if tt.wantErr {
                if err == nil {
                    t.Errorf("expected error but got none, cmd=%q, args=%v", cmd, args)
                }
                return
            }
            if err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }
            if cmd != tt.wantCmd {
                t.Errorf("cmd = %q, want %q", cmd, tt.wantCmd)
            }
            if !reflect.DeepEqual(args, tt.wantArgs) {
                t.Errorf("args = %#v, want %#v", args, tt.wantArgs)
            }
        })
    }
}
