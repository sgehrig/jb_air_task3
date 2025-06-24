package cli

import (
	"srg.de/jb/air_task3/survey"
)

type ResponsesCommand struct{}

func (c *ResponsesCommand) Name() string { return "responses" }

func (c *ResponsesCommand) Aliases() []string { return []string{"response", "resp"} }

func (c *ResponsesCommand) Run(cmd string, args []string, data *survey.SurveyData) (bool, error) {
	responses := data.Responses
	var showKeys []string
	if len(args) > 0 {
		queryString := args[0]
		query, err := survey.ParseResponseQuery(queryString)
		if err != nil {
			return true, err
		}
		responses = query.Limit(responses)
		showKeys = query.Keys
	}

	outputResponses(data.Schema, responses, showKeys)
	return true, nil
}
