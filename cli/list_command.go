package cli

import (
	"fmt"

	"srg.de/jb/air_task3/reader"
)

type ListCommand struct {
}

func (c *ListCommand) Name() string { return "list" }

func (c *ListCommand) Aliases() []string { return []string{"ls"} }

func (c *ListCommand) Run(cmd string, args []string, data *reader.SurveyData) (bool, error) {
	// Output in the order as in the file (insertion order of Schema)
	fmt.Println("Survey Questions:")
	outputSchemaEntries(data.Schema)
	return true, nil
}
