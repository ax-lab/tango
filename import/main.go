package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/ax-lab/tango/import/db"
	"github.com/ax-lab/tango/import/frequency"
	"github.com/ax-lab/tango/import/jmdict"
	"github.com/ax-lab/tango/import/jmnedict"
	"github.com/ax-lab/tango/import/kanji"
)

const (
	EntriesDB = "entries.db"
	NamesDB   = "names.db"
	KanjiDB   = "kanji.db"
	FreqDB    = "freq.db"
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
	importIfNotExists(path.Join(outputDir, KanjiDB), importKanji)
	importIfNotExists(path.Join(outputDir, FreqDB), importFrequency)
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
	var (
		entries []*jmdict.Entry
		tags    map[string]string
	)
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
		tags = decoder.Tags
	}

	opImportEntries.Complete()
	fmt.Printf("... Loaded %d entries\n", len(entries))

	opWriteEntries := Start("writing " + EntriesDB)
	if db, dbErr := db.NewEntriesWriter(outputFile); dbErr != nil {
		ExitWithError("creating %s: %v", EntriesDB, dbErr)
	} else {
		defer db.Close()
		db.WriteTags(tags)
		if writeErr := db.WriteEntries(entries); writeErr != nil {
			ExitWithError("writing entries to %s: %v", EntriesDB, writeErr)
		}
	}
	opWriteEntries.Complete()
}

func importNames(outputFile string) {
	var (
		names []*jmnedict.Name
		tags  map[string]string
	)
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
		tags = decoder.Tags
	}

	opImportNames.Complete()
	fmt.Printf("... Loaded %d names\n", len(names))

	opWriteNames := Start("writing " + NamesDB)
	if db, dbErr := db.NewNamesWriter(outputFile); dbErr != nil {
		ExitWithError("creating %s: %v", NamesDB, dbErr)
	} else {
		defer db.Close()
		db.WriteTags(tags)
		if writeErr := db.WriteNames(names); writeErr != nil {
			ExitWithError("writing names to %s: %v", NamesDB, writeErr)
		}
	}
	opWriteNames.Complete()
}

func importKanji(outputFile string) {
	var (
		characters []*kanji.Character
		info       kanji.Info
	)
	opImportKanji := Start("importing kanji")
	if input, err := kanji.Load(kanji.DefaultFilePath); err != nil {
		ExitWithError("could not load Kanji data: %v", err)
	} else {
		decoder := kanji.NewDecoder(input)
		for {
			entry, err := decoder.ReadCharacter()
			if err != nil {
				ExitWithError("importing characters: %v", err)
			} else if entry == nil {
				break
			}
			characters = append(characters, entry)
		}
		info = decoder.Info
	}

	opImportKanji.Complete()
	fmt.Printf("... Loaded %d kanji\n", len(characters))

	opWriteKanji := Start("writing " + KanjiDB)
	if db, dbErr := db.NewKanjiWriter(outputFile); dbErr != nil {
		ExitWithError("creating %s: %v", KanjiDB, dbErr)
	} else {
		defer db.Close()
		db.WriteInfo(info)
		if writeErr := db.WriteCharacters(characters); writeErr != nil {
			ExitWithError("writing kanji to %s: %v", KanjiDB, writeErr)
		}
	}
	opWriteKanji.Complete()
}

func importFrequency(outputFile string) {
	opImportFrequency := Start("importing frequency")
	jparser, mecab, kanji, err := frequency.LoadPairs(frequency.DefaultPairsFile)
	if err != nil {
		ExitWithError("could not load frequency pairs: %v", err)
	}

	words, chars, err := frequency.LoadInfo(frequency.DefaultInfoFile)
	if err != nil {
		ExitWithError("could not load frequency info: %v", err)
	}
	opImportFrequency.Complete()
	fmt.Printf("... Loaded frequency information\n")

	opWriteFrequency := Start("writing " + FreqDB)
	if db, dbErr := db.NewFrequencyWriter(outputFile); dbErr != nil {
		ExitWithError("creating %s: %v", FreqDB, dbErr)
	} else {
		defer db.Close()
		db.AddCharInfo(chars)
		db.AddWordInfo(words)
		db.AddCharPairs(kanji)
		db.AddWordPairs(jparser, mecab)
		if writeErr := db.Write(); writeErr != nil {
			ExitWithError("writing frequency to %s: %v", FreqDB, writeErr)
		}
	}
	opWriteFrequency.Complete()
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
