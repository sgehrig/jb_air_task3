package cli

import (
	"fmt"
	"sort"
	"strings"

	"srg.de/jb/air_task3/survey"
)

type AnalyzeCommand struct{}

func (c *AnalyzeCommand) Name() string { return "analyze" }

func (c *AnalyzeCommand) Aliases() []string { return []string{"distribution", "dist"} }

func (c *AnalyzeCommand) Run(cmd string, args []string, data *survey.SurveyData) (bool, error) {
	if len(args) < 1 {
		return true, fmt.Errorf("usage: analyze <question-key>")
	}
	key := args[0]
	entry, ok := data.Schema.Get(key)
	if !ok {
		return true, fmt.Errorf("question key not found: %s", key)
	}
	if entry.QType != survey.SC && entry.QType != survey.MC {
		return true, fmt.Errorf("analyze only supports single or multiple choice questions")
	}
	// Count answers
	counts := make(map[string]int)
	missing := 0
	for _, resp := range data.Responses {
		val, ok := resp[key]
		if !ok || !val.IsPresent() {
			missing++
			continue
		}
		switch entry.QType {
		case survey.SC:
			ans, _ := val.AsString()
			counts[ans]++
		case survey.MC:
			anss, _ := val.AsStringSlice()
			for _, ans := range anss {
				counts[ans]++
			}
		}
	}
	// Print distribution
	total := 0
	for _, v := range counts {
		total += v
	}
	totalAll := total + missing
	fmt.Printf("Distribution for %s: %s\n", key, entry.Text)
	options := append([]string{}, entry.Options...)
	sort.Strings(options)
	// Determine max width for option column (up to 25)
	maxOptLen := 0
	for _, opt := range options {
		if l := len(opt); l > maxOptLen {
			maxOptLen = l
		}
	}
	if maxOptLen > 25 {
		maxOptLen = 25
	}
	format := fmt.Sprintf("  %%-%ds : %%7d  (%%6.2f%%%%) %%s\n", maxOptLen)
	ellipsis := "…"
	trim := func(s string) string {
		if len(s) > maxOptLen {
			if maxOptLen > len(ellipsis) {
				return s[:maxOptLen-len(ellipsis)] + ellipsis
			}
			return s[:maxOptLen]
		}
		return s
	}
	bar := func(count int) string {
		barLen := 0
		if totalAll > 0 {
			barLen = int(float64(count)*15.0/float64(totalAll) + 0.5)
		}
		if barLen > 15 {
			barLen = 15
		}
		return strings.Repeat("█", barLen) + strings.Repeat(" ", 15-barLen)
	}
	for _, opt := range options {
		cnt := counts[opt]
		share := 0.0
		if totalAll > 0 {
			share = float64(cnt) * 100.0 / float64(totalAll)
		}
		fmt.Printf(format, trim(opt), cnt, share, bar(cnt))
	}
	shareMissing := 0.0
	if totalAll > 0 {
		shareMissing = float64(missing) * 100.0 / float64(totalAll)
	}
	fmt.Printf(format, "(missing)", missing, shareMissing, "")
	fmt.Printf(format, "(total)", totalAll, 100.0, "")
	return true, nil
}
