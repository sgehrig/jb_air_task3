package cli

import (
    "fmt"
    "strings"

    "srg.de/jb/air_task3/survey"
)

type AnalyzeCommand struct{}

func (c *AnalyzeCommand) Name() string { return "analyze" }

func (c *AnalyzeCommand) Aliases() []string { return []string{"distribution", "dist"} }

func (c *AnalyzeCommand) Run(cmd string, args []string, data *survey.SurveyData) (bool, error) {
    if len(args) < 1 {
        return true, fmt.Errorf("missing question key")
    }
    questionKey := args[0]
    entry, ok := data.Schema.Get(questionKey)
    if !ok {
        return true, fmt.Errorf("question %q not found", questionKey)
    }
    if entry.QType != survey.SC && entry.QType != survey.MC {
        return true, fmt.Errorf("question %q is not single or multi choice", questionKey)
    }

    // Prepare option counts
    counts := map[string]int{}
    for _, opt := range entry.UsedOptions {
        counts[opt] = 0
    }
    counts["(n/a)"] = 0

    for _, resp := range data.Responses {
        val, ok := resp[entry.Key]
        if !ok || !val.Present() {
            counts["(n/a)"]++
            continue
        }
        if entry.QType == survey.SC {
            strVal, ok := val.AsString()
            if !ok || strVal == "" {
                counts["(n/a)"]++
                continue
            }
            if _, exists := counts[strVal]; exists {
                counts[strVal]++
            } else {
                counts[strVal] = 1
            }
        } else if entry.QType == survey.MC {
            strVals, ok := val.AsStringSlice()
            if !ok || len(strVals) == 0 {
                counts["(n/a)"]++
                continue
            }
            seen := map[string]bool{}
            for _, opt := range strVals {
                if _, exists := counts[opt]; exists {
                    counts[opt]++
                } else {
                    counts[opt] = 1
                }
                seen[opt] = true
            }
            // If no valid options found, count as n/a
            if len(seen) == 0 {
                counts["(n/a)"]++
            }
        }
    }

    // Output
    // Find max width for option column (capped at 25)
    maxOptLen := 0
    for opt := range counts {
        l := len(opt)
        if l > 25 {
            l = 25
        }
        if l > maxOptLen {
            maxOptLen = l
        }
    }
    if maxOptLen < 5 {
        maxOptLen = 5
    }
    optFmt := fmt.Sprintf("  %%-%ds", maxOptLen)
    numFmt := "%6d"
    pctFmt := "%7.1f%%%%"
    graphLen := 15

    fmt.Printf("Distribution for [%s] (%s):\n", entry.Key, entry.QType)
    total := 0
    for _, cnt := range counts {
        total += cnt
    }
    for opt, cnt := range counts {
        percent := 0.0
        if total > 0 {
            percent = float64(cnt) * 100.0 / float64(total)
        }
        displayOpt := opt
        if len(displayOpt) > 25 {
            displayOpt = displayOpt[:22] + "..."
        }
        // ASCII bar graph
        barCount := 0
        if total > 0 {
            barCount = int((float64(cnt)/float64(total))*float64(graphLen) + 0.5)
        }
        if barCount > graphLen {
            barCount = graphLen
        }
        bar := strings.Repeat("█", barCount) + strings.Repeat(" ", graphLen-barCount)

        fmt.Printf(optFmt, displayOpt)
        fmt.Printf(" "+numFmt+" "+pctFmt+" |%s|\n", cnt, percent, bar)
    }
    // Align "Total" with the count column
    fmt.Printf(optFmt, "Total")
    fmt.Printf(" "+numFmt+" "+pctFmt+"\n", total, 100.0)
    return true, nil
}
