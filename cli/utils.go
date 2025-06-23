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
    if ((entry.QType == reader.SC) || entry.QType == reader.MC) && (len(entry.UsedOptions) > 0) {
        fmt.Println("    Used options:")
        for _, opt := range entry.UsedOptions {
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

func outputResponseValue(entry *reader.SchemaEntry, resp reader.Response) {
    val, ok := resp[entry.Key]
    if !ok || !val.Present() {
        fmt.Printf("    %s: n/a\n", entry.Key)
        return
    }
    switch entry.QType {
    case reader.SC:
        s, ok := val.AsString()
        if ok {
            fmt.Printf("    %s: %s\n", entry.Key, s)
        } else {
            fmt.Printf("    %s: (invalid))\n", entry.Key)
        }
    case reader.MC:
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
    case reader.TE:
        s, ok := val.AsString()
        if ok {
            fmt.Printf("    %s: %s\n", entry.Key, s)
            break
        } else {
            fmt.Printf("    %s: (invalid)\n", entry.Key)
        }
    }
}

func outputResponse(schema reader.Schema, resp reader.Response) {

    for _, entry := range schema {
        outputResponseValue(entry, resp)
    }
}

func outputResponses(schema reader.Schema, responses []reader.Response) {
    for i, resp := range responses {
        fmt.Printf("Response %d:\n", i+1)
        outputResponse(schema, resp)
    }
}
