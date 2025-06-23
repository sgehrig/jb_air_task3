package cli

import (
	"fmt"

	"srg.de/jb/air_task3/reader"
)

type ListCommand struct {
}

func (c *ListCommand) Name() string { return "list" }

func (c *ListCommand) Run(cmd string, args []string, data *reader.SurveyData) (bool, error) {
	// Output in the order as in the file (insertion order of Schema)
	fmt.Println("Survey Questions:")
	i := 1
	width := len(fmt.Sprintf("%d", len(data.Schema)))
	for _, entry := range data.Schema {
		fmt.Printf("%0*d. [%s] (%s)\n", width, i, entry.Key, entry.QType)
		fmt.Printf("    %s\n", entry.Text)
		i++
	}
	return true, nil
}
