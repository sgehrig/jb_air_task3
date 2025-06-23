package cli

import (
    "fmt"
    "sort"
    "strings"

    "srg.de/jb/air_task3/survey"
)

type Command interface {
    Name() string
    Aliases() []string
    Run(cmd string, args []string, data *survey.SurveyData) (bool, error)
}

var (
    commands = []Command{
        &ClearCommand{},
        &QuitCommand{},
        &ListCommand{},
        &SearchCommand{},
        &ResponsesCommand{},
        &SubsetCommand{},
        &AnalyzeCommand{},
    }
)

type CommandSet map[string]Command

func (cs CommandSet) addCommand(cmd Command) error {
    if _, ok := cs[cmd.Name()]; ok {
        return fmt.Errorf("command already exists: %s", cmd.Name())
    }
    for _, alias := range cmd.Aliases() {
        if _, ok := cs[alias]; ok {
            return fmt.Errorf("command alias already exists: %s", alias)
        }
    }
    cs[cmd.Name()] = cmd
    return nil
}

func (cs CommandSet) Help() string {
    names := make([]struct {
        key  string
        name string
    }, 0, len(cs))
    for n := range cs {
        name := "'" + n + "'"
        aliases := cs[n].Aliases()
        if len(aliases) > 0 {
            name += " ('" + strings.Join(aliases, "', '") + "')"
        }
        names = append(names, struct {
            key  string
            name string
        }{
            key:  n,
            name: name,
        })
    }

    if len(names) == 0 {
        return "No commands available"
    }

    // Custom sort: "clear" first, "quit" last, rest alphabetical
    sort.Slice(names, func(i, j int) bool {
        if names[i].key == "clear" {
            return true
        }
        if names[j].key == "clear" {
            return false
        }
        if names[i].key == "quit" {
            return false
        }
        if names[j].key == "quit" {
            return true
        }
        return names[i].key < names[j].key
    })

    namesStr := names[0].name
    if len(names) > 1 {
        for _, name := range names[1 : len(names)-1] {
            namesStr += ", " + name.name
        }
        namesStr += " or " + names[len(names)-1].name
    }
    return namesStr
}

func (cs CommandSet) Get(name string) (Command, error) {
    if cmd, ok := cs[name]; ok {
        return cmd, nil
    }
    for _, cmd := range cs {
        for _, alias := range cmd.Aliases() {
            if alias == name {
                return cmd, nil
            }
        }
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
