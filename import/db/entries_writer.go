package db

import (
	"database/sql"

	"github.com/ax-lab/tango/import/jmdict"
	_ "github.com/mattn/go-sqlite3"
)

type EntriesWriter struct {
	db *sql.DB
}

func NewEntriesWriter(outputFile string) (*EntriesWriter, error) {
	db, err := sql.Open("sqlite3", outputFile)
	if err != nil {
		return nil, err
	}

	_, execErr := db.Exec(`
		DROP TABLE IF EXISTS entries_kanji;
		DROP TABLE IF EXISTS entries;

		CREATE TABLE entries (
			sequence INTEGER NOT NULL PRIMARY KEY
		);

		CREATE TABLE entries_kanji (
			sequence INTEGER,
			position INTEGER,
			text     TEXT,
			info     TEXT,
			priority TEXT,
			PRIMARY KEY (sequence, position),
			FOREIGN KEY (sequence) REFERENCES entries(sequence)
		);
	`)
	if execErr != nil {
		db.Close()
		return nil, execErr
	}

	return &EntriesWriter{
		db: db,
	}, nil
}

func (writer *EntriesWriter) Close() {
	writer.db.Close()
}

func (writer *EntriesWriter) WriteEntries(entries []*jmdict.Entry) error {
	tx, err := writer.db.Begin()
	if err != nil {
		return err
	}

	insertEntry, err := tx.Prepare(`
		INSERT INTO entries(sequence) VALUES (?)
	`)
	if err != nil {
		return err
	}

	insertEntryKanji, err := tx.Prepare(`
		INSERT INTO entries_kanji
		(sequence, position, text, info, priority)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if _, cmdErr := insertEntry.Exec(entry.Sequence); cmdErr != nil {
			tx.Rollback()
			return cmdErr
		}

		for pos, kanji := range entry.Kanji {
			if _, cmdErr := insertEntryKanji.Exec(entry.Sequence, pos, kanji.Text, kanji.Info, kanji.Priority); cmdErr != nil {
				tx.Rollback()
				return cmdErr
			}
		}
	}

	return tx.Commit()
}
