package main

import (
	"fmt"
	"os"

	"srg.de/jb/air_task3/reader"
)

func main() {
	jsonFile := "so_2024.json"
	xlsxFile := "so_2024_raw.xlsx"

	data, err := reader.ReadSurveyDataCached(jsonFile, xlsxFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load survey data: %v\n", err)
		os.Exit(1)
	}

	width := len(fmt.Sprintf("%d", len(data.Schema)))

	fmt.Println("Survey Questions:")
	var i int
	for _, v := range data.Schema {
		fmt.Printf("%0*d. [%s] (%s)\n    %s\n", width, i+1, v.Key, v.QType, v.Text)
		i++
	}
}
