package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/ax-lab/tango/import/db"
	"github.com/ax-lab/tango/import/jmdict"
	"github.com/ax-lab/tango/import/jmnedict"
)

const (
	EntriesDB = "entries.db"
	NamesDB   = "names.db"
)

func main() {
	outputDirPtr := flag.String("output", "", "output directory")
	flag.Parse()

	outputDir := *outputDirPtr
	if outputDir == "" {
		ExitWithError("invalid arguments: missing output directory")
	}

	fmt.Printf("Importing dictionary data to `%s`...\n", outputDir)

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		ExitWithError("failed to create output directory: %v", err)
	}

	importIfNotExists(path.Join(outputDir, EntriesDB), importEntries)
	importIfNotExists(path.Join(outputDir, NamesDB), importNames)
}

func importIfNotExists(outputFile string, callback func(outputFile string)) {
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		callback(outputFile)
	} else if err != nil {
		ExitWithError("stat error: %v", err)
	} else {
		fmt.Printf("... `%s` already exists, skipping\n", outputFile)
	}
}

func importEntries(outputFile string) {
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
	fmt.Printf("... Loaded %d entries\n", len(entries))

	opWriteEntries := Start("writing " + EntriesDB)
	if db, dbErr := db.NewEntriesWriter(outputFile); dbErr != nil {
		ExitWithError("creating %s: %v", EntriesDB, dbErr)
	} else {
		defer db.Close()
		if writeErr := db.WriteEntries(entries); writeErr != nil {
			ExitWithError("writing entries to %s: %v", EntriesDB, writeErr)
		}
	}
	opWriteEntries.Complete()
}

func importNames(outputFile string) {
	var names []*jmnedict.Name
	opImportNames := Start("importing names")
	if input, err := jmnedict.Load(jmnedict.DefaultFilePath); err != nil {
		ExitWithError("could not load JMnedict data: %v", err)
	} else {
		decoder := jmnedict.NewDecoder(input)
		for {
			entry, err := decoder.ReadEntry()
			if err != nil {
				ExitWithError("importing names: %v", err)
			} else if entry == nil {
				break
			}
			names = append(names, entry)
		}
	}

	opImportNames.Complete()
	fmt.Printf("... Loaded %d names\n", len(names))

	opWriteNames := Start("writing " + NamesDB)
	if db, dbErr := db.NewNamesWriter(outputFile); dbErr != nil {
		ExitWithError("creating %s: %v", NamesDB, dbErr)
	} else {
		defer db.Close()
		if writeErr := db.WriteNames(names); writeErr != nil {
			ExitWithError("writing names to %s: %v", NamesDB, writeErr)
		}
	}
	opWriteNames.Complete()
}

type Timer struct {
	name  string
	start time.Time
}

func Start(name string) Timer {
	fmt.Printf("\n--> Started %s...\n", name)
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
