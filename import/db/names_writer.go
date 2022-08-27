package db

import (
	"database/sql"

	"github.com/ax-lab/tango/import/jmnedict"
)

type NamesWriter struct {
	db *sql.DB
}

func NewNamesWriter(outputFile string) (*NamesWriter, error) {
	db, err := Open(outputFile, `
		DROP TABLE IF EXISTS names_sense;
		DROP TABLE IF EXISTS names;

		CREATE TABLE names (
			sequence INTEGER NOT NULL PRIMARY KEY,
			kanji    TEXT,
			reading  TEXT
		);

		CREATE TABLE names_sense (
			sequence    INTEGER,
			position    INTEGER,
			info        TEXT,
			xref        TEXT,
			translation TEXT,
			PRIMARY KEY (sequence, position),
			FOREIGN KEY (sequence) REFERENCES names(sequence)
		);
	`)
	if err != nil {
		return nil, err
	}

	return &NamesWriter{
		db: db,
	}, nil
}

func (writer *NamesWriter) Close() {
	writer.db.Close()
}

func (writer *NamesWriter) WriteNames(names []*jmnedict.Name) error {
	tx := BeginTransaction(writer.db)

	insertName := tx.Prepare(`
		INSERT INTO names(sequence, kanji, reading) VALUES (?, ?, ?)
	`)

	insertNameSense := tx.Prepare(`
		INSERT INTO names_sense
		(sequence, position, info, xref, translation)
		VALUES (?, ?, ?, ?, ?)
	`)

	for _, name := range names {
		insertName.Exec(name.Sequence, csv(name.Kanji), csv(name.Reading))
		for pos, sense := range name.Sense {
			insertNameSense.Exec(
				name.Sequence, pos,
				csv(sense.Info), csv(sense.XRef), csv(sense.Translation))
		}
	}

	return tx.Finish()
}
