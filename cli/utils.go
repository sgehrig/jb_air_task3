package cli

import (
	"fmt"
	"slices"

	"srg.de/jb/air_task3/survey"
)

func outputSchemaEntry(entry *survey.SchemaEntry, i int, w int) {
	if i > 0 {
		if w < 1 {
			w = 1
		}
		fmt.Printf("%0*d. [%s] (%s)\n", w, i, entry.Key, entry.QType)
	} else {
		fmt.Printf("[%s] (%s)\n", entry.Key, entry.QType)
	}
	fmt.Printf("    %s\n", entry.Text)
	if ((entry.QType == survey.SC) || entry.QType == survey.MC) && (len(entry.Options) > 0) {
		fmt.Println("    Used options:")
		for _, opt := range entry.Options {
			fmt.Printf("        - %s\n", opt)
		}
	}
}

func outputSchemaEntries(entries []*survey.SchemaEntry) {
	i := 1
	width := len(fmt.Sprintf("%d", len(entries)))
	for _, entry := range entries {
		outputSchemaEntry(entry, i, width)
		i++
	}
}

func outputResponseValue(entry *survey.SchemaEntry, resp survey.SurveyResponse) {
	val, ok := resp[entry.Key]
	if !ok || !val.IsPresent() {
		fmt.Printf("    %s: n/a\n", entry.Key)
		return
	}
	switch entry.QType {
	case survey.SC:
		s, ok := val.AsString()
		if ok {
			fmt.Printf("    %s: %s\n", entry.Key, s)
		} else {
			fmt.Printf("    %s: (invalid))\n", entry.Key)
		}
	case survey.MC:
		ss, ok := val.AsStringSlice()
		if ok {
			fmt.Printf("    %s: ", entry.Key)
			for j, opt := range ss {
				if (j > 10) && (len(ss) > 10) {
					fmt.Printf("... (%d more)", len(ss)-j)
					break
				}
				if j > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%s", opt)
			}
			fmt.Println()
		} else {
			fmt.Printf("    %s: (invalid)\n", entry.Key)
		}
	case survey.TE:
		s, ok := val.AsString()
		if ok {
			fmt.Printf("    %s: %s\n", entry.Key, s)
		} else {
			fmt.Printf("    %s: (invalid)\n", entry.Key)
		}
	}
}

func outputResponse(schema survey.Schema, resp survey.SurveyResponse, keys []string) {
	for _, entry := range schema {
		if len(keys) > 0 && !slices.Contains(keys, entry.Key) {
			continue
		}
		outputResponseValue(entry, resp)
	}
}

func outputResponses(schema survey.Schema, responses []survey.SurveyResponse, keys []string) {
	for i, resp := range responses {
		fmt.Printf("Response %d:\n", i+1)
		outputResponse(schema, resp, keys)
	}
}
