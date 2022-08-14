package db

import (
	"database/sql"
	"fmt"

	"github.com/ax-lab/tango/import/jmdict"
	_ "github.com/mattn/go-sqlite3"
)

type EntriesWriter struct {
	db *sql.DB
}

func NewEntriesWriter(outputFile string) (*EntriesWriter, error) {
	fmt.Println(outputFile)

	db, err := sql.Open("sqlite3", outputFile)
	if err != nil {
		return nil, err
	}

	_, execErr := db.Exec(`
		DROP TABLE IF EXISTS entries;

		CREATE TABLE IF NOT EXISTS entries (
			sequence INTEGER NOT NULL PRIMARY KEY
		)
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

	insertEntry, err := tx.Prepare("INSERT INTO entries(sequence) VALUES (?)")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		_, cmdErr := insertEntry.Exec(entry.Sequence)
		if cmdErr != nil {
			tx.Rollback()
			return cmdErr
		}
	}

	return tx.Commit()
}
