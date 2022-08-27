package db_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/ax-lab/tango/import/db"
	"github.com/ax-lab/tango/import/kanji"
	"github.com/stretchr/testify/require"
)

func TestKanjiWriterFailsOnOpenError(t *testing.T) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		w, err := db.NewKanjiWriter(filepath.Join(dbFile, "force-error.db"))
		test.Error(err)
		test.Nil(w)
	})
}

func TestKanjiWriterFailsOnInvalidInsert(t *testing.T) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		w, err := db.NewKanjiWriter(dbFile)
		test.NoError(err)
		err = w.WriteCharacters([]*kanji.Character{
			{Literal: "日"},
			{Literal: "日"},
		})
		test.ErrorContains(err, "constraint")
	})
}

func TestKanjiWriterCanRewriteTheDatabase(t *testing.T) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		w, err := db.NewKanjiWriter(dbFile)
		test.NoError(err)
		w.Close()

		w, err = db.NewKanjiWriter(dbFile)
		test.NoError(err)
		w.Close()
	})
}

func TestKanjiWriterExportsInfo(t *testing.T) {
	testKanji(t,
		func(test *require.Assertions, db *db.KanjiWriter) {
			err := db.WriteInfo(kanji.Info{
				Version: "2022-001",
				Created: "2022-01-01",
			})
			test.NoError(err)
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := []Data{
				{"version": "2022-001", "created": "2022-01-01"},
			}
			actual := testQuery(test, db, "SELECT version, created FROM info")
			test.EqualValues(expected, actual)
		},
	)
}

func TestKanjiWriterExportsCharacter(t *testing.T) {
	testKanji(t,
		func(test *require.Assertions, db *db.KanjiWriter) {
			err := db.WriteCharacters([]*kanji.Character{
				{Literal: "日"},
				{Literal: "本"},
			})
			test.NoError(err)
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := []Data{
				{"literal": "日"},
				{"literal": "本"},
			}
			actual := testQuery(test, db, "SELECT literal FROM character")
			test.EqualValues(expected, actual)
		},
	)
}

func TestKanjiWriterExportsCharacterData(t *testing.T) {
	testKanji(t,
		func(test *require.Assertions, db *db.KanjiWriter) {
			err := db.WriteCharacters([]*kanji.Character{
				{
					Literal:     "日",
					Grade:       8,
					Strokes:     []int{11, 12},
					Frequency:   123,
					JLPT:        4,
					RadicalName: []string{"name1", "name2"},
					Nanori:      []string{"nanori1", "nanori2"},
				},
			})
			test.NoError(err)
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := []Data{
				{
					"literal":      "日",
					"grade":        int64(8),
					"strokes":      "11\t12",
					"frequency":    int64(123),
					"jlpt":         int64(4),
					"radical_name": "name1\tname2",
					"nanori":       "nanori1\tnanori2",
				},
			}
			actual := testQuery(test, db, `
				SELECT
					literal, grade, strokes, frequency, jlpt, radical_name, nanori
				FROM character`)
			test.EqualValues(expected, actual)
		},
	)
}

func TestKanjiWriterExportsCharacterCodepoint(t *testing.T) {
	testKanji(t,
		func(test *require.Assertions, db *db.KanjiWriter) {
			err := db.WriteCharacters([]*kanji.Character{
				{
					Literal: "A",
					Codepoint: []kanji.CharacterCodepoint{
						{Type: "type A1", Text: "text A1"},
						{Type: "type A2", Text: "text A2"},
					},
				},
				{
					Literal: "B",
					Codepoint: []kanji.CharacterCodepoint{
						{Type: "type B1", Text: "text B1"},
					},
				},
				{
					Literal: "X",
				},
			})
			test.NoError(err)
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := []Data{
				{"literal": "A", "position": int64(0), "name": "type A1", "value": "text A1"},
				{"literal": "A", "position": int64(1), "name": "type A2", "value": "text A2"},
				{"literal": "B", "position": int64(0), "name": "type B1", "value": "text B1"},
			}
			actual := testQuery(test, db,
				`SELECT literal, position, name, value FROM character_codepoint`)
			test.EqualValues(expected, actual)
		},
	)
}

func TestKanjiWriterExportsCharacterRadical(t *testing.T) {
	testKanji(t,
		func(test *require.Assertions, db *db.KanjiWriter) {
			err := db.WriteCharacters([]*kanji.Character{
				{
					Literal: "A",
					Radical: []kanji.CharacterRadical{
						{Type: "type A1", Text: "text A1"},
						{Type: "type A2", Text: "text A2"},
					},
				},
				{
					Literal: "B",
					Radical: []kanji.CharacterRadical{
						{Type: "type B1", Text: "text B1"},
					},
				},
				{
					Literal: "X",
				},
			})
			test.NoError(err)
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := []Data{
				{"literal": "A", "position": int64(0), "name": "type A1", "value": "text A1"},
				{"literal": "A", "position": int64(1), "name": "type A2", "value": "text A2"},
				{"literal": "B", "position": int64(0), "name": "type B1", "value": "text B1"},
			}
			actual := testQuery(test, db,
				`SELECT literal, position, name, value FROM character_radical`)
			test.EqualValues(expected, actual)
		},
	)
}

func TestKanjiWriterExportsCharacterVariant(t *testing.T) {
	testKanji(t,
		func(test *require.Assertions, db *db.KanjiWriter) {
			err := db.WriteCharacters([]*kanji.Character{
				{
					Literal: "A",
					Variant: []kanji.CharacterVariant{
						{Type: "type A1", Text: "text A1"},
						{Type: "type A2", Text: "text A2"},
					},
				},
				{
					Literal: "B",
					Variant: []kanji.CharacterVariant{
						{Type: "type B1", Text: "text B1"},
					},
				},
				{
					Literal: "X",
				},
			})
			test.NoError(err)
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := []Data{
				{"literal": "A", "position": int64(0), "name": "type A1", "value": "text A1"},
				{"literal": "A", "position": int64(1), "name": "type A2", "value": "text A2"},
				{"literal": "B", "position": int64(0), "name": "type B1", "value": "text B1"},
			}
			actual := testQuery(test, db,
				`SELECT literal, position, name, value FROM character_variant`)
			test.EqualValues(expected, actual)
		},
	)
}

func TestKanjiWriterExportsCharacterReference(t *testing.T) {
	testKanji(t,
		func(test *require.Assertions, db *db.KanjiWriter) {
			err := db.WriteCharacters([]*kanji.Character{
				{
					Literal: "A",
					Reference: []kanji.CharacterReference{
						{Type: "type A1", Text: "text A1"},
						{Type: "type A2", Text: "text A2"},
					},
				},
				{
					Literal: "B",
					Reference: []kanji.CharacterReference{
						{Type: "type B1", Text: "text B1", Volume: "v1", Page: "p1"},
					},
				},
				{
					Literal: "X",
				},
			})
			test.NoError(err)
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := []Data{
				{"literal": "A", "position": int64(0), "name": "type A1", "value": "text A1", "volume": "", "page": ""},
				{"literal": "A", "position": int64(1), "name": "type A2", "value": "text A2", "volume": "", "page": ""},
				{
					"literal":  "B",
					"position": int64(0),
					"name":     "type B1",
					"value":    "text B1",
					"volume":   "v1",
					"page":     "p1",
				},
			}
			actual := testQuery(test, db,
				`SELECT literal, position, name, value, volume, page FROM character_reference`)
			test.EqualValues(expected, actual)
		},
	)
}

func TestKanjiWriterExportsCharacterQueryCode(t *testing.T) {
	testKanji(t,
		func(test *require.Assertions, db *db.KanjiWriter) {
			err := db.WriteCharacters([]*kanji.Character{
				{
					Literal: "A",
					QueryCode: []kanji.CharacterQueryCode{
						{Type: "type A1", Text: "text A1"},
						{Type: "type A2", Text: "text A2"},
					},
				},
				{
					Literal: "B",
					QueryCode: []kanji.CharacterQueryCode{
						{Type: "type B1", Text: "text B1", SkipMisclass: "miss"},
					},
				},
				{
					Literal: "X",
				},
			})
			test.NoError(err)
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := []Data{
				{"literal": "A", "position": int64(0), "name": "type A1", "value": "text A1", "skip_misclass": ""},
				{"literal": "A", "position": int64(1), "name": "type A2", "value": "text A2", "skip_misclass": ""},
				{
					"literal":       "B",
					"position":      int64(0),
					"name":          "type B1",
					"value":         "text B1",
					"skip_misclass": "miss",
				},
			}
			actual := testQuery(test, db,
				`SELECT literal, position, name, value, skip_misclass FROM character_query_code`)
			test.EqualValues(expected, actual)
		},
	)
}

func TestKanjiWriterExportsCharacterReadingMeaning(t *testing.T) {
	testKanji(t,
		func(test *require.Assertions, db *db.KanjiWriter) {
			err := db.WriteCharacters([]*kanji.Character{
				{
					Literal: "A",
					ReadingMeanings: []kanji.CharacterReadingMeaning{
						{
							Reading: []kanji.CharacterReading{
								{Type: "type A1x", Text: "read A1x"},
								{Type: "type A1y", Text: "read A1y"},
							},
							Meaning: []kanji.CharacterMeaning{
								{Lang: "lang A1x", Text: "mean A1x"},
								{Lang: "lang A1y", Text: "mean A1y"},
							},
						},
						{
							Reading: []kanji.CharacterReading{
								{Type: "type A2x", Text: "read A2x"},
							},
							Meaning: []kanji.CharacterMeaning{
								{Lang: "lang A2x", Text: "mean A2x"},
							},
						},
					},
				},
				{
					Literal: "B",
					ReadingMeanings: []kanji.CharacterReadingMeaning{
						{
							Reading: []kanji.CharacterReading{
								{Type: "type B1x", Text: "read B1x"},
								{Type: "type B1y", Text: "read B1y"},
							},
							Meaning: []kanji.CharacterMeaning{
								{Lang: "lang B1x", Text: "mean B1x"},
								{Lang: "lang B1y", Text: "mean B1y"},
							},
						},
					},
				},
				{
					Literal: "X",
				},
			})
			test.NoError(err)
		},
		func(test *require.Assertions, db *sql.DB) {
			expected_reading := []Data{
				{"literal": "A", "group_pos": int64(0), "position": int64(0), "name": "type A1x", "value": "read A1x"},
				{"literal": "A", "group_pos": int64(0), "position": int64(1), "name": "type A1y", "value": "read A1y"},
				{"literal": "A", "group_pos": int64(1), "position": int64(0), "name": "type A2x", "value": "read A2x"},
				{"literal": "B", "group_pos": int64(0), "position": int64(0), "name": "type B1x", "value": "read B1x"},
				{"literal": "B", "group_pos": int64(0), "position": int64(1), "name": "type B1y", "value": "read B1y"},
			}
			actual_reading := testQuery(test, db,
				`SELECT literal, position, group_pos, name, value FROM character_reading`)
			test.EqualValues(expected_reading, actual_reading)

			expected_meaning := []Data{
				{"literal": "A", "group_pos": int64(0), "position": int64(0), "lang": "lang A1x", "value": "mean A1x"},
				{"literal": "A", "group_pos": int64(0), "position": int64(1), "lang": "lang A1y", "value": "mean A1y"},
				{"literal": "A", "group_pos": int64(1), "position": int64(0), "lang": "lang A2x", "value": "mean A2x"},
				{"literal": "B", "group_pos": int64(0), "position": int64(0), "lang": "lang B1x", "value": "mean B1x"},
				{"literal": "B", "group_pos": int64(0), "position": int64(1), "lang": "lang B1y", "value": "mean B1y"},
			}
			actual_meaning := testQuery(test, db,
				`SELECT literal, position, group_pos, lang, value FROM character_meaning`)
			test.EqualValues(expected_meaning, actual_meaning)
		},
	)
}

func testKanji(t *testing.T, prepare func(test *require.Assertions, db *db.KanjiWriter), eval func(test *require.Assertions, db *sql.DB)) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		func() {
			db, dbErr := db.NewKanjiWriter(dbFile)
			if dbErr != nil {
				panic(dbErr)
			}

			defer db.Close()
			prepare(test, db)
		}()

		func() {
			db, err := sql.Open("sqlite3", dbFile)
			if err != nil {
				panic(err)
			}
			defer db.Close()
			eval(test, db)
		}()
	})
}
