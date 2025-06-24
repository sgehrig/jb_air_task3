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

func TestCommandSet_addCommand_DuplicateName(t *testing.T) {
	cs := CommandSet{}
	cmd1 := &dummyCmd{name: "foo"}
	cmd2 := &dummyCmd{name: "foo"}
	if err := cs.addCommand(cmd1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err := cs.addCommand(cmd2)
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected duplicate name error, got: %v", err)
	}
}

func TestCommandSet_addCommand_DuplicateAlias(t *testing.T) {
	cs := CommandSet{}
	cmd1 := &dummyCmd{name: "foo", aliases: []string{"bar"}}
	cmd2 := &dummyCmd{name: "baz", aliases: []string{"bar"}}
	if err := cs.addCommand(cmd1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err := cs.addCommand(cmd2)
	if err == nil || !strings.Contains(err.Error(), "command alias already exists") {
		t.Errorf("expected duplicate alias error, got: %v", err)
	}
}

func TestCommandSet_Help(t *testing.T) {
	cs := CommandSet{}
	cs.addCommand(&dummyCmd{name: "clear"})
	cs.addCommand(&dummyCmd{name: "foo", aliases: []string{"f"}})
	cs.addCommand(&dummyCmd{name: "bar"})
	cs.addCommand(&dummyCmd{name: "quit"})
	help := cs.Help()
	expected := "'clear', 'bar', 'foo' ('f') or 'quit'"
	if help != expected {
		t.Errorf("Help() output mismatch.\nGot:      %q\nExpected: %q", help, expected)
	}
}
