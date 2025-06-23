package cli

import (
	"srg.de/jb/air_task3/survey"
)

type ResponsesCommand struct{}

func (c *ResponsesCommand) Name() string { return "responses" }

func (c *ResponsesCommand) Aliases() []string { return []string{"resp", "r"} }

func (c *ResponsesCommand) Run(cmd string, args []string, data *survey.SurveyData) (bool, error) {
	outputResponses(data.Schema, data.Responses)
	return true, nil
}
