package cli

import (
    "fmt"
    "slices"
    "strings"

    "srg.de/jb/air_task3/survey"
)

type SubsetCommand struct{}

func (c *SubsetCommand) Name() string { return "subsets" }

func (c *SubsetCommand) Aliases() []string { return []string{"subset", "sub"} }

func (c *SubsetCommand) Run(cmd string, args []string, data *survey.SurveyData) (bool, error) {
    if len(args) < 2 {
        return true, fmt.Errorf("missing question and/or option")
    }
    questionKey := args[0]
    option := strings.ToLower(args[1])

    entry, ok := data.Schema.Get(questionKey)
    if !ok {
        return true, fmt.Errorf("question %q not found", questionKey)
    }

    found := []survey.Response{}
    for _, resp := range data.Responses {
        val, ok := resp[entry.Key]
        if !ok || !val.Present() {
            continue
        }
        if entry.QType == survey.SC {
            strVal, ok := val.AsString()
            if !ok {
                continue
            }
            if strings.ToLower(strVal) == option {
                found = append(found, resp)
            }
        } else if entry.QType == survey.MC {
            strVal, ok := val.AsStringSlice()
            if !ok {
                continue
            }
            for _, opt := range strVal {
                if strings.ToLower(opt) == option {
                    found = append(found, resp)
                    break
                }
            }
        }
    }

    showKeys := []string{questionKey}
    if len(args) > 2 {
        queryString := args[2]
        query, err := survey.ParseResponseQuery(queryString)
        if err != nil {
            return true, err
        }
        found = query.Limit(found)
        if slices.Contains(query.Keys, "*") {
            showKeys = nil
        } else {
            showKeys = append(showKeys, query.Keys...)
        }
    }

    outputResponses(data.Schema, found, showKeys)
    return true, nil
}
