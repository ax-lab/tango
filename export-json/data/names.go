package data

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"io"
	"os"
	"path"
)

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

func (sense NameSense) MarshalJSON() ([]byte, error) {
	info := TSV(sense.Info)
	xref := TSV(sense.XRef)
	translation := TSV(sense.Translation)
	return json.Marshal([]any{info, xref, translation})
}

func EncodeNameAndSenses(output io.Writer, name Name, senses []NameSense) (err error) {
	out := func(v interface{}) {
		if err == nil {
			var bytes []byte
			bytes, err = json.Marshal(v)
			if err == nil {
				_, err = output.Write(bytes)
			}
		}
	}

	var (
		kanji   = TSV(name.Kanji)
		reading = TSV(name.Reading)
	)

	out([]any{name.Sequence, kanji, reading, senses})

	return err
}

func (tb *NameSense) Query() string {
	return "SELECT sequence, info, xref, translation FROM name_sense"
}

func (tb *NameSense) Read(row *sql.Rows) error {
	return row.Scan(&tb.Sequence, &tb.Info, &tb.XRef, &tb.Translation)
}

func ExportNames(importDir, exportDir string) error {
	db := OpenDB(importDir, "names.db")
	defer db.Close()

	name := Name{}
	nameReader := db.ScanTable(&name)
	defer nameReader.Close()

	sense := NameSense{}
	senseReader := db.ScanTable(&sense)
	defer senseReader.Close()

	outputPath := path.Join(exportDir, "names-2.json")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	output := bufio.NewWriter(outputFile)

	if _, err = output.WriteString("[\n"); err != nil {
		return err
	}

	var first = true
	var curSenses []NameSense
	for nameReader.Next() {
		for senseReader.Next() {
			if sense.Sequence != name.Sequence {
				senseReader.Unget()
				break
			}
			curSenses = append(curSenses, sense)
		}

		if !first {
			if _, err = output.WriteString(",\n"); err != nil {
				return err
			}
		}
		first = false

		if err = EncodeNameAndSenses(output, name, curSenses); err != nil {
			return err
		}
		curSenses = curSenses[:0]
	}

	if _, err = output.WriteString("\n]"); err != nil {
		return err
	}

	output.Flush()

	return db.Done()
}
