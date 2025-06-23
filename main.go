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
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    width := len(fmt.Sprintf("%d", len(data.Schema)))

    // Output in the order as in the file (insertion order of Schema)
    fmt.Println("Survey Questions:")
    i := 1
    for _, entry := range data.Schema {
        fmt.Printf("%0*d. [%s] (%s)\n", width, i, entry.Key, entry.QType)
        fmt.Printf("    %s\n", entry.Text)
        i++
    }
}
