package cli

import (
    "fmt"
    "sort"
    "strings"

    "srg.de/jb/air_task3/reader"
)

type Command interface {
    Name() string
    Run(cmd string, args []string, data *reader.SurveyData) (bool, error)
}

var (
    commands = []Command{
        &ClearCommand{},
        &QuitCommand{},
        &ListCommand{},
    }
)

type CommandSet map[string]Command

func (cs CommandSet) addCommand(cmd Command) error {
    if _, ok := cs[cmd.Name()]; ok {
        return fmt.Errorf("command already exists: %s", cmd.Name())
    }
    cs[cmd.Name()] = cmd
    return nil
}

func (cs CommandSet) Names() []string {
    names := make([]string, 0, len(cs))
    for name := range cs {
        names = append(names, name)
    }
    // Custom sort: "clear" first, "quit" last, rest alphabetical
    sort.Slice(names, func(i, j int) bool {
        if names[i] == "clear" {
            return true
        }
        if names[j] == "clear" {
            return false
        }
        if names[i] == "quit" {
            return false
        }
        if names[j] == "quit" {
            return true
        }
        return names[i] < names[j]
    })
    return names
}

func (cs CommandSet) Help() string {
    names := cs.Names()
    for i, n := range names {
        names[i] = "'" + n + "'"
    }
    var namesStr string
    if len(names) == 1 {
        namesStr = names[0]
    } else if len(names) > 1 {
        namesStr = strings.Join(names[:len(names)-1], ", ") + ", or " + names[len(names)-1]
    }
    return namesStr
}

func (cs CommandSet) Get(name string) (Command, error) {
    if cmd, ok := cs[name]; ok {
        return cmd, nil
    }
    return nil, fmt.Errorf("command not found: %s", name)
}

func InitCommands() (CommandSet, error) {
    cs := CommandSet{}

    for _, cmd := range commands {
        err := cs.addCommand(cmd)
        if err != nil {
            return nil, fmt.Errorf("error adding command %s: %w", cmd.Name(), err)
        }
    }
    return cs, nil
}
