package cli

import (
	"fmt"

	"srg.de/jb/air_task3/reader"
)

func outputSchemaEntry(entry *reader.SchemaEntry, i int, w int) {
	if i > 0 {
		if w < 1 {
			w = 1
		}
		fmt.Printf("%0*d. [%s] (%s)\n", w, i, entry.Key, entry.QType)
	} else {
		fmt.Printf("[%s] (%s)\n", entry.Key, entry.QType)
	}
	fmt.Printf("    %s\n", entry.Text)
	if ((entry.QType == reader.SC) || entry.QType == reader.MC) && (len(entry.Options) > 0) {
		fmt.Println("    Used options:")
		for _, opt := range entry.Options {
			fmt.Printf("        - %s\n", opt)
		}
	}
}

func outputSchemaEntries(entries []*reader.SchemaEntry) {
	i := 1
	width := len(fmt.Sprintf("%d", len(entries)))
	for _, entry := range entries {
		outputSchemaEntry(entry, i, width)
		i++
	}
}
