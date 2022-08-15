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
	db, err := Open(outputFile, `
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
	if err != nil {
		return nil, err
	}

	return &EntriesWriter{
		db: db,
	}, nil
}

func (writer *EntriesWriter) Close() {
	writer.db.Close()
}

func (writer *EntriesWriter) WriteEntries(entries []*jmdict.Entry) error {
	tx := BeginTransaction(writer.db)

	insertEntry := tx.Prepare(`
		INSERT INTO entries(sequence) VALUES (?)
	`)

	insertEntryKanji := tx.Prepare(`
		INSERT INTO entries_kanji
		(sequence, position, text, info, priority)
		VALUES (?, ?, ?, ?, ?)
	`)

	for _, entry := range entries {
		insertEntry.Exec(entry.Sequence)
		for pos, kanji := range entry.Kanji {
			insertEntryKanji.Exec(entry.Sequence, pos, kanji.Text, kanji.Info, kanji.Priority)
		}
	}

	return tx.Finish()
}
