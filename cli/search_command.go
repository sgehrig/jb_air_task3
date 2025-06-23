package cli

import (
    "fmt"

    "srg.de/jb/air_task3/survey"
)

type SearchCommand struct {
}

func (c *SearchCommand) Name() string { return "search" }

func (c *SearchCommand) Aliases() []string { return []string{"find", "query"} }

func (c *SearchCommand) Run(cmd string, args []string, data *survey.SurveyData) (bool, error) {
    if len(args) == 0 {
        return true, fmt.Errorf("missing search term")
    }

    // Output in the order as in the file (insertion order of Schema)
    fmt.Println("Survey Questions matching: '" + args[0] + "':")
    outputSchemaEntries(data.Schema.SearchForString(args[0]))
    return true, nil
}
