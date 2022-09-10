package main

import (
	"fmt"
	"os"
	"path"

	"github.com/ax-lab/tango/common"
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
	outputDir := common.GetOutputDir("output", "output directory")
	fmt.Printf("Importing dictionary data to `%s`...\n", outputDir)

	importIfNotExists(path.Join(outputDir, EntriesDB), importEntries)
	importIfNotExists(path.Join(outputDir, NamesDB), importNames)
	importIfNotExists(path.Join(outputDir, KanjiDB), importKanji)
	importIfNotExists(path.Join(outputDir, FreqDB), importFrequency)
}

func importIfNotExists(outputFile string, callback func(outputFile string)) {
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		callback(outputFile)
	} else if err != nil {
		common.ExitWithError("stat error: %v", err)
	} else {
		fmt.Printf("... `%s` already exists, skipping\n", outputFile)
	}
}

func importEntries(outputFile string) {
	var (
		entries []*jmdict.Entry
		tags    map[string]string
	)
	opImportEntries := common.Start("importing entries")
	if input, err := jmdict.Load(jmdict.DefaultFilePath); err != nil {
		common.ExitWithError("could not load JMdict data: %v", err)
	} else {
		decoder := jmdict.NewDecoder(input)
		for {
			entry, err := decoder.ReadEntry()
			if err != nil {
				common.ExitWithError("importing entries: %v", err)
			} else if entry == nil {
				break
			}
			entries = append(entries, entry)
		}
		tags = decoder.Tags
	}

	opImportEntries.Complete()
	fmt.Printf("... Loaded %d entries\n", len(entries))

	opWriteEntries := common.Start("writing " + EntriesDB)
	if db, dbErr := db.NewEntriesWriter(outputFile); dbErr != nil {
		common.ExitWithError("creating %s: %v", EntriesDB, dbErr)
	} else {
		defer db.Close()
		db.WriteTags(tags)
		if writeErr := db.WriteEntries(entries); writeErr != nil {
			common.ExitWithError("writing entries to %s: %v", EntriesDB, writeErr)
		}
	}
	opWriteEntries.Complete()
}

func importNames(outputFile string) {
	var (
		names []*jmnedict.Name
		tags  map[string]string
	)
	opImportNames := common.Start("importing names")
	if input, err := jmnedict.Load(jmnedict.DefaultFilePath); err != nil {
		common.ExitWithError("could not load JMnedict data: %v", err)
	} else {
		decoder := jmnedict.NewDecoder(input)
		for {
			entry, err := decoder.ReadEntry()
			if err != nil {
				common.ExitWithError("importing names: %v", err)
			} else if entry == nil {
				break
			}
			names = append(names, entry)
		}
		tags = decoder.Tags
	}

	opImportNames.Complete()
	fmt.Printf("... Loaded %d names\n", len(names))

	opWriteNames := common.Start("writing " + NamesDB)
	if db, dbErr := db.NewNamesWriter(outputFile); dbErr != nil {
		common.ExitWithError("creating %s: %v", NamesDB, dbErr)
	} else {
		defer db.Close()
		db.WriteTags(tags)
		if writeErr := db.WriteNames(names); writeErr != nil {
			common.ExitWithError("writing names to %s: %v", NamesDB, writeErr)
		}
	}
	opWriteNames.Complete()
}

func importKanji(outputFile string) {
	var (
		characters []*kanji.Character
		info       kanji.Info
	)
	opImportKanji := common.Start("importing kanji")
	if input, err := kanji.Load(kanji.DefaultFilePath); err != nil {
		common.ExitWithError("could not load Kanji data: %v", err)
	} else {
		decoder := kanji.NewDecoder(input)
		for {
			entry, err := decoder.ReadCharacter()
			if err != nil {
				common.ExitWithError("importing characters: %v", err)
			} else if entry == nil {
				break
			}
			characters = append(characters, entry)
		}
		info = decoder.Info
	}

	opImportKanji.Complete()
	fmt.Printf("... Loaded %d kanji\n", len(characters))

	opWriteKanji := common.Start("writing " + KanjiDB)
	if db, dbErr := db.NewKanjiWriter(outputFile); dbErr != nil {
		common.ExitWithError("creating %s: %v", KanjiDB, dbErr)
	} else {
		defer db.Close()
		db.WriteInfo(info)
		if writeErr := db.WriteCharacters(characters); writeErr != nil {
			common.ExitWithError("writing kanji to %s: %v", KanjiDB, writeErr)
		}
	}
	opWriteKanji.Complete()
}

func importFrequency(outputFile string) {
	opImportFrequency := common.Start("importing frequency")
	jparser, mecab, kanji, err := frequency.LoadPairs(frequency.DefaultPairsFile)
	if err != nil {
		common.ExitWithError("could not load frequency pairs: %v", err)
	}

	words, chars, err := frequency.LoadInfo(frequency.DefaultInfoFile)
	if err != nil {
		common.ExitWithError("could not load frequency info: %v", err)
	}
	opImportFrequency.Complete()
	fmt.Printf("... Loaded frequency information\n")

	opWriteFrequency := common.Start("writing " + FreqDB)
	if db, dbErr := db.NewFrequencyWriter(outputFile); dbErr != nil {
		common.ExitWithError("creating %s: %v", FreqDB, dbErr)
	} else {
		defer db.Close()
		db.AddCharInfo(chars)
		db.AddWordInfo(words)
		db.AddCharPairs(kanji)
		db.AddWordPairs(jparser, mecab)
		if writeErr := db.Write(); writeErr != nil {
			common.ExitWithError("writing frequency to %s: %v", FreqDB, writeErr)
		}
	}
	opWriteFrequency.Complete()
}
