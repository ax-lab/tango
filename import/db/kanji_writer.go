package db

import (
	"database/sql"
	"fmt"

	"github.com/ax-lab/tango/import/kanji"
)

type KanjiWriter struct {
	db *sql.DB
}

func NewKanjiWriter(outputFile string) (*KanjiWriter, error) {
	db, err := Open(outputFile, `
		DROP TABLE IF EXISTS info;

		DROP TABLE IF EXISTS character_codepoint;
		DROP TABLE IF EXISTS character_radical;
		DROP TABLE IF EXISTS character_variant;
		DROP TABLE IF EXISTS character_reference;
		DROP TABLE IF EXISTS character_query_code;
		DROP TABLE IF EXISTS character_reading;
		DROP TABLE IF EXISTS character_meaning;

		DROP TABLE IF EXISTS character;

		CREATE TABLE info (
			created TEXT,
			version TEXT
		);

		CREATE TABLE character (
			literal      TEXT PRIMARY KEY,
			grade        INTEGER,
			strokes      TEXT,
			frequency    INTEGER,
			jlpt         INTEGER,
			radical_name TEXT,
			nanori       TEXT
		);

		CREATE TABLE character_codepoint (
			literal      TEXT,
			position     INTEGER,
			name         TEXT,
			value        TEXT,
			PRIMARY KEY (literal, position),
			FOREIGN KEY (literal) REFERENCES character(literal)
		);

		CREATE TABLE character_radical (
			literal      TEXT,
			position     INTEGER,
			name         TEXT,
			value        TEXT,
			PRIMARY KEY (literal, position),
			FOREIGN KEY (literal) REFERENCES character(literal)
		);

		CREATE TABLE character_variant (
			literal      TEXT,
			position     INTEGER,
			name         TEXT,
			value        TEXT,
			PRIMARY KEY (literal, position),
			FOREIGN KEY (literal) REFERENCES character(literal)
		);

		CREATE TABLE character_reference (
			literal      TEXT,
			position     INTEGER,
			name         TEXT,
			value        TEXT,
			volume       TEXT,
			page         TEXT,
			PRIMARY KEY (literal, position),
			FOREIGN KEY (literal) REFERENCES character(literal)
		);

		CREATE TABLE character_query_code (
			literal       TEXT,
			position      INTEGER,
			name          TEXT,
			value         TEXT,
			skip_misclass TEXT,
			PRIMARY KEY (literal, position),
			FOREIGN KEY (literal) REFERENCES character(literal)
		);

		CREATE TABLE character_reading (
			literal       TEXT,
			group_pos     INTEGER,
			position      INTEGER,
			name          TEXT,
			value         TEXT,
			PRIMARY KEY (literal, group_pos, position),
			FOREIGN KEY (literal) REFERENCES character(literal)
		);

		CREATE TABLE character_meaning (
			literal       TEXT,
			group_pos     INTEGER,
			position      INTEGER,
			lang          TEXT,
			value         TEXT,
			PRIMARY KEY (literal, group_pos, position),
			FOREIGN KEY (literal) REFERENCES character(literal)
		);
	`)
	if err != nil {
		return nil, err
	}

	return &KanjiWriter{
		db: db,
	}, nil
}

func (writer *KanjiWriter) Close() {
	writer.db.Close()
}

func (writer *KanjiWriter) WriteInfo(info kanji.Info) error {
	tx := BeginTransaction(writer.db)

	tx.Exec("DELETE FROM info")
	tx.Exec("INSERT INTO info(created, version) VALUES (?, ?)", info.Created, info.Version)
	return tx.Finish()
}

func (writer *KanjiWriter) WriteCharacters(characters []*kanji.Character) error {
	tx := BeginTransaction(writer.db)

	insertCharacter := tx.Prepare(`
		INSERT INTO character
		(literal, grade, strokes, frequency, jlpt, radical_name, nanori)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)

	insertCharacterCodepoint := tx.Prepare(`
		INSERT INTO character_codepoint
		(literal, position, name, value)
		VALUES (?, ?, ?, ?)
	`)

	insertCharacterRadical := tx.Prepare(`
		INSERT INTO character_radical
		(literal, position, name, value)
		VALUES (?, ?, ?, ?)
	`)

	insertCharacterVariant := tx.Prepare(`
		INSERT INTO character_variant
		(literal, position, name, value)
		VALUES (?, ?, ?, ?)
	`)

	insertCharacterReference := tx.Prepare(`
		INSERT INTO character_reference
		(literal, position, name, value, volume, page)
		VALUES (?, ?, ?, ?, ?, ?)
	`)

	insertCharacterQueryCode := tx.Prepare(`
		INSERT INTO character_query_code
		(literal, position, name, value, skip_misclass)
		VALUES (?, ?, ?, ?, ?)
	`)

	insertCharacterReading := tx.Prepare(`
		INSERT INTO character_reading
		(literal, group_pos, position, name, value)
		VALUES (?, ?, ?, ?, ?)
	`)

	insertCharacterMeaning := tx.Prepare(`
		INSERT INTO character_meaning
		(literal, group_pos, position, lang, value)
		VALUES (?, ?, ?, ?, ?)
	`)

	for _, character := range characters {
		var (
			grade     interface{}
			frequency interface{}
			jlpt      interface{}
			strokes   []string
		)
		if character.Grade > 0 {
			grade = character.Grade
		}
		if character.Frequency > 0 {
			frequency = character.Frequency
		}
		if character.JLPT > 0 {
			jlpt = character.JLPT
		}
		for _, it := range character.Strokes {
			strokes = append(strokes, fmt.Sprint(it))
		}
		insertCharacter.Exec(character.Literal, grade, csv(strokes), frequency, jlpt,
			csv(character.RadicalName), csv(character.Nanori))

		for i, it := range character.Codepoint {
			insertCharacterCodepoint.Exec(character.Literal, i, it.Type, it.Text)
		}

		for i, it := range character.Radical {
			insertCharacterRadical.Exec(character.Literal, i, it.Type, it.Text)
		}

		for i, it := range character.Variant {
			insertCharacterVariant.Exec(character.Literal, i, it.Type, it.Text)
		}

		for i, it := range character.Reference {
			insertCharacterReference.Exec(character.Literal, i, it.Type, it.Text, it.Volume, it.Page)
		}

		for i, it := range character.QueryCode {
			insertCharacterQueryCode.Exec(character.Literal, i, it.Type, it.Text, it.SkipMisclass)
		}

		for pos, group := range character.ReadingMeanings {
			for i, it := range group.Reading {
				insertCharacterReading.Exec(character.Literal, pos, i, it.Type, it.Text)
			}
			for i, it := range group.Meaning {
				insertCharacterMeaning.Exec(character.Literal, pos, i, it.Lang, it.Text)
			}
		}
	}

	return tx.Finish()
}
