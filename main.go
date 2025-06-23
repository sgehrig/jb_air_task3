package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "time"

    "srg.de/jb/air_task3/cli"
    "srg.de/jb/air_task3/reader"
)

func main() {
    xlsxFile := "so_2024_raw.xlsx"

    fmt.Printf("Loading survey data from %s...\n", xlsxFile)
    start := time.Now()
    data, err := reader.ReadSurveyDataCached(xlsxFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("Loaded survey data in %s\n", time.Since(start))

    commandSet, err := cli.InitCommands()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    }

    scanner := bufio.NewScanner(os.Stdin)
    fmt.Printf("Survey CLI. Type %s.\n", commandSet.Help())
    for {
        fmt.Print("(survey)> ")
        if !scanner.Scan() {
            break
        }
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        args := strings.Fields(line)
        cmd := args[0]
        cmdArgs := args[1:]

        cmdInstance, err := commandSet.Get(cmd)
        if err != nil {
            fmt.Printf("Unknown command: %s. Args: %v\n", cmd, cmdArgs)
            continue
        }

        cont, err := cmdInstance.Run(cmd, cmdArgs, data)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        }
        if !cont {
            break
        }
    }
}
