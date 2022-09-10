package db

import (
	"database/sql"
	"strings"

	"github.com/ax-lab/tango/import/jmdict"
)

type EntriesWriter struct {
	db *sql.DB
}

func NewEntriesWriter(outputFile string) (*EntriesWriter, error) {
	db, err := Open(outputFile, `
		DROP TABLE IF EXISTS tag;

		DROP TABLE IF EXISTS entry_kanji;
		DROP TABLE IF EXISTS entry_reading;
		DROP TABLE IF EXISTS entry_sense_source;
		DROP TABLE IF EXISTS entry_sense_glossary;
		DROP TABLE IF EXISTS entry_sense;

		DROP TABLE IF EXISTS entry;

		CREATE TABLE tag (
			name TEXT,
			desc TEXT
		);

		CREATE TABLE entry (
			sequence INTEGER NOT NULL PRIMARY KEY
		);

		CREATE TABLE entry_kanji (
			sequence INTEGER,
			position INTEGER,
			text     TEXT,
			info     TEXT,
			priority TEXT,
			PRIMARY KEY (sequence, position),
			FOREIGN KEY (sequence) REFERENCES entry(sequence)
		);

		CREATE TABLE entry_reading (
			sequence    INTEGER,
			position    INTEGER,
			text        TEXT,
			info        TEXT,
			priority    TEXT,
			restriction TEXT,
			no_kanji    INTEGER,
			PRIMARY KEY (sequence, position),
			FOREIGN KEY (sequence) REFERENCES entry(sequence)
		);

		CREATE TABLE entry_sense (
			sequence    INTEGER,
			position    INTEGER,
			info        TEXT,
			pos         TEXT,
			stagk       TEXT,
			stagr       TEXT,
			field       TEXT,
			misc        TEXT,
			dialect     TEXT,
			antonym     TEXT,
			xref        TEXT,
			PRIMARY KEY (sequence, position),
			FOREIGN KEY (sequence) REFERENCES entry(sequence)
		);

		CREATE TABLE entry_sense_glossary (
			sequence    INTEGER,
			sense       INTEGER,
			position    INTEGER,
			text        TEXT,
			lang        TEXT,
			type        TEXT,
			PRIMARY KEY (sequence, sense, position),
			FOREIGN KEY (sequence) REFERENCES entry(sequence),
			FOREIGN KEY (sequence, sense) REFERENCES entry_sense(sequence, position)
		);

		CREATE TABLE entry_sense_source (
			sequence    INTEGER,
			sense       INTEGER,
			position    INTEGER,
			text        TEXT,
			lang        TEXT,
			type        TEXT,
			wasei       TEXT,
			PRIMARY KEY (sequence, sense, position),
			FOREIGN KEY (sequence) REFERENCES entry(sequence),
			FOREIGN KEY (sequence, sense) REFERENCES entry_sense(sequence, position)
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

func (writer *EntriesWriter) WriteTags(tags map[string]string) error {
	return writeTags(writer.db, tags)
}

func (writer *EntriesWriter) WriteEntries(entries []*jmdict.Entry) error {
	tx := BeginTransaction(writer.db)

	insertEntry := tx.Prepare(`
		INSERT INTO entry(sequence) VALUES (?)
	`)

	insertEntryKanji := tx.Prepare(`
		INSERT INTO entry_kanji
		(sequence, position, text, info, priority)
		VALUES (?, ?, ?, ?, ?)
	`)

	insertEntryReading := tx.Prepare(`
		INSERT INTO entry_reading
		(sequence, position, text, info, priority, restriction, no_kanji)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)

	insertEntrySense := tx.Prepare(`
		INSERT INTO entry_sense
		(sequence, position, info, pos, stagk, stagr, field, misc, dialect, antonym, xref)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)

	insertEntrySenseGlossary := tx.Prepare(`
		INSERT INTO entry_sense_glossary
		(sequence, sense, position, text, lang, type)
		VALUES (?, ?, ?, ?, ?, ?)
	`)

	insertEntrySenseSource := tx.Prepare(`
		INSERT INTO entry_sense_source
		(sequence, sense, position, text, lang, type, wasei)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)

	for _, entry := range entries {
		insertEntry.Exec(entry.Sequence)
		for pos, kanji := range entry.Kanji {
			insertEntryKanji.Exec(entry.Sequence, pos, kanji.Text, csv(kanji.Info), csv(kanji.Priority))
		}
		for pos, reading := range entry.Reading {
			insertEntryReading.Exec(
				entry.Sequence, pos, reading.Text,
				csv(reading.Info), csv(reading.Priority), csv(reading.Restriction), reading.NoKanji)
		}
		for pos, sense := range entry.Sense {
			insertEntrySense.Exec(
				entry.Sequence, pos,
				csv(sense.Info), csv(sense.PartOfSpeech), csv(sense.StagKanji),
				csv(sense.StagReading), csv(sense.Field), csv(sense.Misc),
				csv(sense.Dialect), csv(sense.Antonym), csv(sense.XRef))

			for pos_glossary, glossary := range sense.Glossary {
				insertEntrySenseGlossary.Exec(
					entry.Sequence, pos, pos_glossary,
					glossary.Text, glossary.Lang, glossary.Type)
			}

			for pos_source, source := range sense.Source {
				insertEntrySenseSource.Exec(
					entry.Sequence, pos, pos_source,
					source.Text, source.Lang, source.Type, source.Wasei)
			}
		}
	}

	return tx.Finish()
}

func csv(values []string) string {
	return strings.Join(values, "\t")
}
