package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/ax-lab/tango/common"
)

func ExportEntries(importDir, exportDir string) error {
	db := OpenDB(importDir, "entries.db")
	defer db.Close()

	out := NewDataWriter(exportDir)
	defer out.Close()

	entry := Entry{}
	entryReader := db.ScanTable(&entry)
	defer entryReader.Close()

	entryKanji := EntryKanji{}
	entryKanjiReader := db.ScanTable(&entryKanji)
	defer entryKanjiReader.Close()

	entryReading := EntryReading{}
	entryReadingReader := db.ScanTable(&entryReading)
	defer entryReadingReader.Close()

	entrySense := EntrySense{}
	entrySenseReader := db.ScanTable(&entrySense)
	defer entrySenseReader.Close()

	entrySenseGlossary := EntrySenseGlossary{}
	entrySenseGlossaryReader := db.ScanTable(&entrySenseGlossary)
	defer entrySenseGlossaryReader.Close()

	entrySenseSource := EntrySenseSource{}
	entrySenseSourceReader := db.ScanTable(&entrySenseSource)
	defer entrySenseSourceReader.Close()

	var (
		entries = make([]string, 0)
		indexes = make(map[string][]int)

		reverse      = make([]string, 0)
		reverseIndex = make(map[string][]int)
	)

	pushEntry := func(entry string, sequence int) {
		key := common.KanaToCommonHiragana(entry)
		index := indexes[key]
		if len(index) == 0 {
			entries = append(entries, key)
		}
		indexes[key] = append(index, sequence)
	}

	pushReverse := func(entry string, sequence int) {
		if entry == "" {
			return
		}
		key := strings.ToLower(entry)
		index := reverseIndex[key]
		if len(index) == 0 {
			reverse = append(reverse, key)
		}
		indexes[key] = append(index, sequence)
	}

	var (
		curKanji   []EntryKanji
		curReading []EntryReading
		curSense   []EntrySense
	)
	for entryReader.Next() {
		for entryKanjiReader.Next() {
			if entryKanji.Sequence != entry.Sequence {
				entryKanjiReader.Unget()
				break
			}
			curKanji = append(curKanji, entryKanji)
		}

		for entryReadingReader.Next() {
			if entryReading.Sequence != entry.Sequence {
				entryReadingReader.Unget()
				break
			}
			curReading = append(curReading, entryReading)
		}

		for entrySenseReader.Next() {
			if entrySense.Sequence != entry.Sequence {
				entrySenseReader.Unget()
				break
			}
			for entrySenseGlossaryReader.Next() {
				if entrySenseGlossary.Sequence != entry.Sequence {
					entrySenseGlossaryReader.Unget()
					break
				}
				entrySense.Glossary = append(entrySense.Glossary, entrySenseGlossary)
			}

			for entrySenseSourceReader.Next() {
				if entrySenseSource.Sequence != entry.Sequence {
					entrySenseSourceReader.Unget()
					break
				}
				entrySense.Source = append(entrySense.Source, entrySenseSource)
			}

			curSense = append(curSense, entrySense)
			entrySense.Glossary = nil
			entrySense.Source = nil
		}

		index := entry.Sequence % 1000
		fileName := fmt.Sprintf("entries/entry-%03d.json", index)

		for _, it := range curKanji {
			pushEntry(it.Text, entry.Sequence)
		}
		for _, it := range curReading {
			pushEntry(it.Text, entry.Sequence)
		}
		for _, sense := range curSense {
			for _, it := range sense.Glossary {
				pushReverse(it.Text, entry.Sequence)
			}
		}

		out.AppendToFile(fileName, []any{entry.Sequence, curKanji, curReading, curSense})

		curKanji = curKanji[:0]
		curReading = curReading[:0]
		curSense = curSense[:0]
	}

	sort.Strings(entries)
	for _, it := range entries {
		char := (([]rune)(it))[0] % 0x100
		file := fmt.Sprintf("entry-index/entry-%02X.json", char)
		out.AppendToFile(file, []any{it, indexes[it]})
	}

	sort.Strings(reverse)
	for _, it := range reverse {
		char := it[0]
		if char < 'a' || char > 'z' {
			char = '0'
		}
		file := fmt.Sprintf("entry-index/entry-reverse-%c.json", char)
		out.AppendToFile(file, []any{it, indexes[it]})
	}

	tag := Tag{}
	tagReader := db.ScanTable(&tag)
	for tagReader.Next() {
		out.AppendToFileIndent("entry-tags.json", tag)
	}

	out.Close()
	err := db.Error()
	if err == nil {
		err = out.Error()
	}

	return err
}

type Entry struct {
	Sequence int
}

func (tb *Entry) Query() string {
	return "SELECT sequence FROM entry"
}

func (tb *Entry) Read(row *sql.Rows) error {
	return row.Scan(&tb.Sequence)
}

type EntryKanji struct {
	Sequence int
	Text     string
	Info     []string
	Priority []string
}

func (item EntryKanji) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{item.Text, item.Info, item.Priority})
}

func (tb *EntryKanji) Query() string {
	return "SELECT sequence, text, info, priority FROM entry_kanji"
}

func (tb *EntryKanji) Read(row *sql.Rows) (err error) {
	var (
		info     string
		priority string
	)
	err = row.Scan(&tb.Sequence, &tb.Text, &info, &priority)
	tb.Info = TSV(info)
	tb.Priority = TSV(priority)
	return err
}

type EntryReading struct {
	Sequence    int
	Text        string
	Info        []string
	Priority    []string
	Restriction []string
	NoKanji     int
}

func (item EntryReading) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{item.Text, item.Info, item.Priority, item.Restriction, item.NoKanji})
}

func (tb *EntryReading) Query() string {
	return "SELECT sequence, text, info, priority, restriction, no_kanji FROM entry_reading"
}

func (tb *EntryReading) Read(row *sql.Rows) (err error) {
	var (
		info        string
		priority    string
		restriction string
	)
	err = row.Scan(&tb.Sequence, &tb.Text, &info, &priority, &restriction, &tb.NoKanji)
	tb.Info = TSV(info)
	tb.Priority = TSV(priority)
	tb.Restriction = TSV(restriction)
	return err
}

type EntrySense struct {
	Sequence int
	Position int
	Info     string
	Pos      []string
	StagK    []string
	StagR    []string
	Field    []string
	Misc     []string
	Dialect  []string
	Antonym  []string
	Xref     []string

	Glossary []EntrySenseGlossary
	Source   []EntrySenseSource
}

func (item EntrySense) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{item.Info, item.Pos, item.StagK, item.StagR, item.Field, item.Misc, item.Dialect, item.Antonym, item.Xref, item.Glossary, item.Source})
}

func (tb *EntrySense) Query() string {
	return "SELECT sequence, position, info, pos, stagk, stagr, field, misc, dialect, antonym, xref FROM entry_sense"
}

func (tb *EntrySense) Read(row *sql.Rows) (err error) {
	var (
		pos     string
		stagK   string
		stagR   string
		field   string
		misc    string
		dialect string
		antonym string
		xref    string
	)
	err = row.Scan(&tb.Sequence, &tb.Position, &tb.Info, &pos, &stagK, &stagR, &field, &misc, &dialect, &antonym, &xref)
	tb.Pos = TSV(pos)
	tb.StagK = TSV(stagK)
	tb.StagR = TSV(stagR)
	tb.Field = TSV(field)
	tb.Misc = TSV(misc)
	tb.Dialect = TSV(dialect)
	tb.Antonym = TSV(antonym)
	tb.Xref = TSV(xref)
	return err
}

type EntrySenseGlossary struct {
	Sequence int
	Sense    int
	Text     string
	Lang     string
	Type     string
}

func (item EntrySenseGlossary) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{item.Text, item.Lang, item.Type})
}

func (tb *EntrySenseGlossary) Query() string {
	return "SELECT sequence, sense, text, lang, type FROM entry_sense_glossary"
}

func (tb *EntrySenseGlossary) Read(row *sql.Rows) error {
	return row.Scan(&tb.Sequence, &tb.Sense, &tb.Text, &tb.Lang, &tb.Type)
}

type EntrySenseSource struct {
	Sequence int
	Sense    int
	Text     string
	Lang     string
	Type     string
	Wasei    bool
}

func (item EntrySenseSource) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{item.Text, item.Lang, item.Type, item.Wasei})
}

func (tb *EntrySenseSource) Query() string {
	return "SELECT sequence, sense, text, lang, type, wasei FROM entry_sense_source"
}

func (tb *EntrySenseSource) Read(row *sql.Rows) (err error) {
	var wasei string
	err = row.Scan(&tb.Sequence, &tb.Sense, &tb.Text, &tb.Lang, &tb.Type, &wasei)
	tb.Wasei = wasei == "y"
	return err
}
