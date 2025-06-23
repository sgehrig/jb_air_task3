package cli

import (
	"fmt"

	"srg.de/jb/air_task3/reader"
)

type QuitCommand struct{}

func (c *QuitCommand) Name() string { return "quit" }

func (c *QuitCommand) Aliases() []string { return []string{"exit"} }

func (c *QuitCommand) Run(cmd string, args []string, data *reader.SurveyData) (bool, error) {
	fmt.Println("Bye!")
	return false, nil
}
