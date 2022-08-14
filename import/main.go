package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/ax-lab/tango/import/db"
	"github.com/ax-lab/tango/import/jmdict"
)

const (
	EntriesDB = "entries.db"
)

func main() {
	outputDirPtr := flag.String("output", "", "output directory")
	flag.Parse()

	outputDir := *outputDirPtr
	if outputDir == "" {
		ExitWithError("invalid arguments: missing output directory")
	}

	fmt.Printf("Importing dictionary data to `%s`...\n\n", outputDir)

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		ExitWithError("failed to create output directory: %v", err)
	}

	var entries []*jmdict.Entry
	opImportEntries := Start("importing entries")
	if input, err := jmdict.Load(jmdict.DefaultFilePath); err != nil {
		ExitWithError("could not load JMdict data: %v", err)
	} else {
		decoder := jmdict.NewDecoder(input)
		for {
			entry, err := decoder.ReadEntry()
			if err != nil {
				ExitWithError("importing entries: %v", err)
			} else if entry == nil {
				break
			}
			entries = append(entries, entry)
		}
	}

	opImportEntries.Complete()
	fmt.Printf("=== Loaded %d entries", len(entries))

	opWriteEntries := Start("writing " + EntriesDB)
	if db, dbErr := db.NewEntriesWriter(path.Join(outputDir, EntriesDB)); dbErr != nil {
		ExitWithError("creating %s: %v", EntriesDB, dbErr)
	} else {
		defer db.Close()
		if writeErr := db.WriteEntries(entries); writeErr != nil {
			ExitWithError("writing entries to %s: %v", EntriesDB, writeErr)
		}
	}
	opWriteEntries.Complete()
}

type Timer struct {
	name  string
	start time.Time
}

func Start(name string) Timer {
	fmt.Printf("--> Started %s...\n", name)
	return Timer{name, time.Now()}
}

func (t Timer) Complete() {
	fmt.Printf("... %s took %v\n", t.name, time.Since(t.start))
}

func ExitWithError(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg, args...)
	fmt.Fprintln(os.Stderr)
	os.Exit(2)
}
