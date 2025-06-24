package cli

import (
    "strings"
    "testing"

    "srg.de/jb/air_task3/survey"
)

type dummyCmd struct {
    name    string
    aliases []string
}

func (d *dummyCmd) Name() string      { return d.name }
func (d *dummyCmd) Aliases() []string { return d.aliases }
func (d *dummyCmd) Run(cmd string, args []string, data *survey.SurveyData) (bool, error) {
    return true, nil
}

func TestCommandSet_Get(t *testing.T) {
    cs := &CommandSet{}
    cmd1 := &dummyCmd{name: "foo", aliases: []string{"f"}}
    cmd2 := &dummyCmd{name: "bar", aliases: []string{"b", "bee"}}
    cs.addCommand(cmd1)
    cs.addCommand(cmd2)

    // By name
    got, err := cs.Get("foo")
    if err != nil || got != cmd1 {
        t.Errorf("Get(foo) = %v, %v; want %v, nil", got, err, cmd1)
    }
    // By alias
    got, err = cs.Get("b")
    if err != nil || got != cmd2 {
        t.Errorf("Get(b) = %v, %v; want %v, nil", got, err, cmd2)
    }
    got, err = cs.Get("bee")
    if err != nil || got != cmd2 {
        t.Errorf("Get(bee) = %v, %v; want %v, nil", got, err, cmd2)
    }
    // Unknown
    _, err = cs.Get("baz")
    if err == nil {
        t.Errorf("Get(baz) = nil error, want error")
    }
}

func TestCommandSet_DuplicateName(t *testing.T) {
    cs := &CommandSet{}
    cmd1 := &dummyCmd{name: "foo"}
    cmd2 := &dummyCmd{name: "foo"}
    cs.addCommand(cmd1)
    err := cs.addCommand(cmd2)
    if err == nil {
        t.Errorf("expected error for duplicate command name")
    }
}

func TestCommandSet_DuplicateAlias(t *testing.T) {
    cs := &CommandSet{}
    cmd1 := &dummyCmd{name: "foo", aliases: []string{"f"}}
    cmd2 := &dummyCmd{name: "bar", aliases: []string{"f"}}
    cs.addCommand(cmd1)
    err := cs.addCommand(cmd2)
    if err == nil {
        t.Errorf("expected error for duplicate alias")
    }
}

func TestCommandSet_Help(t *testing.T) {
    cs := &CommandSet{}
    cs.addCommand(&dummyCmd{name: "foo", aliases: []string{"f", "fo"}})
    cs.addCommand(&dummyCmd{name: "bar", aliases: []string{"b"}})
    help := cs.Help()
    expected := "'bar' ('b') or 'foo' ('f', 'fo')"
    if help != expected {
        t.Errorf("Help() = %q, want %q", help, expected)
    }
}

func TestCommandSet_HelpFormat(t *testing.T) {
    cs := &CommandSet{}
    cs.addCommand(&dummyCmd{name: "foo", aliases: []string{"f"}})
    help := cs.Help()
    if !strings.Contains(help, "foo") || !strings.Contains(help, "f") {
        t.Errorf("Help() format incorrect: %s", help)
    }
}

func TestCommandSet_Empty(t *testing.T) {
    cs := &CommandSet{}
    help := cs.Help()
    if help == "" {
        t.Errorf("Help() should not be empty for empty CommandSet")
    }
    _, err := cs.Get("foo")
    if err == nil {
        t.Errorf("expected error for Get on empty CommandSet")
    }
}
