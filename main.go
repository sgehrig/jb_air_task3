package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"srg.de/jb/air_task3/cli"
	"srg.de/jb/air_task3/survey"
)

func main() {
	xlsxFile := "so_2024_raw.xlsx"

	fmt.Printf("Loading survey data from %s...\n", xlsxFile)
	start := time.Now()
	data, err := survey.ReadSurveyCached(xlsxFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load survey data: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Loaded survey data in %s\n", time.Since(start))

	commandSet, err := cli.InitCommands()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt: "(survey)> ",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize readline: %v\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	fmt.Printf("Survey CLI. Type %s.\n", commandSet.Help())
	for {
		line, err := rl.Readline()
		if errors.Is(err, readline.ErrInterrupt) {
			if len(line) == 0 {
				break
			}
			continue
		} else if errors.Is(err, io.EOF) {
			break
		}
		line = strings.TrimSpace(line)
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
