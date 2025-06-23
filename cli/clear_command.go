package cli

import (
    "fmt"

    "srg.de/jb/air_task3/survey"
)

type ClearCommand struct{}

func (c *ClearCommand) Name() string { return "clear" }

func (c *ClearCommand) Aliases() []string { return []string{"cls"} }

func (c *ClearCommand) Run(cmd string, args []string, data *reader.SurveyData) (bool, error) {
    fmt.Print("\033[H\033[2J") // ANSI clear screen
    return true, nil
}
