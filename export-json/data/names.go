package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
)

func ExportNames(importDir, exportDir string) error {
	db := OpenDB(importDir, "names.db")
	defer db.Close()

	out := NewDataWriter(exportDir)
	defer out.Close()

	name := Name{}
	nameReader := db.ScanTable(&name)
	defer nameReader.Close()

	sense := NameSense{}
	senseReader := db.ScanTable(&sense)
	defer senseReader.Close()

	var (
		entries = make([]string, 0)
		indexes = make(map[string][]int)
	)

	pushEntry := func(entry string, sequence int) {
		index := indexes[entry]
		if len(index) == 0 {
			entries = append(entries, entry)
		}
		indexes[entry] = append(index, sequence)
	}

	var curSenses []NameSense
	for nameReader.Next() {
		for senseReader.Next() {
			if sense.Sequence != name.Sequence {
				senseReader.Unget()
				break
			}
			curSenses = append(curSenses, sense)
		}

		index := name.Sequence % 1000
		fileName := fmt.Sprintf("names/name-%03d.json", index)

		var (
			sequence = name.Sequence
			kanji    = TSV(name.Kanji)
			reading  = TSV(name.Reading)
		)

		for _, it := range kanji {
			pushEntry(it, sequence)
		}
		for _, it := range reading {
			pushEntry(it, sequence)
		}

		out.AppendToFile(fileName, []any{sequence, kanji, reading, curSenses})

		curSenses = curSenses[:0]
	}

	sort.Strings(entries)
	for _, it := range entries {
		char := (([]rune)(it))[0] % 0x100
		file := fmt.Sprintf("names-index/name-%02X.json", char)
		out.AppendToFile(file, []any{it, indexes[it]})
	}

	tag := Tag{}
	tagReader := db.ScanTable(&tag)
	for tagReader.Next() {
		out.AppendToFileIndent("name-tags.json", tag)
	}

	out.Close()
	err := db.Error()
	if err == nil {
		err = out.Error()
	}

	return err
}

type Name struct {
	Sequence int
	Kanji    string
	Reading  string
}

func (tb *Name) Query() string {
	return "SELECT sequence, kanji, reading FROM name"
}

func (tb *Name) Read(row *sql.Rows) error {
	return row.Scan(&tb.Sequence, &tb.Kanji, &tb.Reading)
}

type NameSense struct {
	Sequence    int
	Info        string
	XRef        string
	Translation string
}

func (tb *NameSense) Query() string {
	return "SELECT sequence, info, xref, translation FROM name_sense"
}

func (tb *NameSense) Read(row *sql.Rows) error {
	return row.Scan(&tb.Sequence, &tb.Info, &tb.XRef, &tb.Translation)
}

func (sense NameSense) MarshalJSON() ([]byte, error) {
	info := TSV(sense.Info)
	xref := TSV(sense.XRef)
	translation := TSV(sense.Translation)
	return json.Marshal([]any{info, xref, translation})
}

type Tag struct {
	Name string
	Desc string
}

func (tb *Tag) Query() string {
	return "SELECT name, desc FROM tag"
}

func (tb *Tag) Read(row *sql.Rows) error {
	return row.Scan(&tb.Name, &tb.Desc)
}
