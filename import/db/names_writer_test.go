package db_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/ax-lab/tango/import/db"
	"github.com/ax-lab/tango/import/jmnedict"
	"github.com/stretchr/testify/require"
)

func TestNamesWriterFailsOnOpenError(t *testing.T) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		w, err := db.NewNamesWriter(filepath.Join(dbFile, "force-error.db"))
		test.Error(err)
		test.Nil(w)
	})
}

func TestNamesWriterFailsOnInvalidInsert(t *testing.T) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		w, err := db.NewNamesWriter(dbFile)
		test.NoError(err)
		err = w.WriteNames([]*jmnedict.Name{
			{Sequence: 100},
			{Sequence: 100},
		})
		test.ErrorContains(err, "constraint")
	})
}

func TestNamesWriterCanRewriteTheDatabase(t *testing.T) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		w, err := db.NewNamesWriter(dbFile)
		test.NoError(err)
		w.Close()

		w, err = db.NewNamesWriter(dbFile)
		test.NoError(err)
		w.Close()
	})
}

func TestNamesWriterExportsNames(t *testing.T) {
	testNames(t,
		func(test *require.Assertions, db *db.NamesWriter) {
			err := db.WriteNames([]*jmnedict.Name{
				{Sequence: 1001},
				{Sequence: 1002},
				{Sequence: 1003},
			})
			test.NoError(err)
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := []Data{
				{"sequence": int64(1001)},
				{"sequence": int64(1002)},
				{"sequence": int64(1003)},
			}
			actual := testQuery(test, db, "SELECT sequence FROM names")
			test.EqualValues(expected, actual)
		},
	)
}

func TestNamesWriterExportsKanji(t *testing.T) {
	testNames(t,
		func(test *require.Assertions, db *db.NamesWriter) {
			err := db.WriteNames([]*jmnedict.Name{
				{
					Sequence: 1001,
					Kanji:    []string{"kanji 1"},
				},
				{
					Sequence: 1002,
					Kanji:    []string{"kanji 2a", "kanji 2b"},
				},
				{
					Sequence: 9999,
				},
			})
			test.NoError(err)
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := []Data{
				{
					"sequence": int64(1001),
					"kanji":    "kanji 1",
				},
				{
					"sequence": int64(1002),
					"kanji":    "kanji 2a\tkanji 2b",
				},
				{
					"sequence": int64(9999),
					"kanji":    "",
				},
			}
			actual := testQuery(test, db, "SELECT sequence, kanji FROM names")
			test.EqualValues(expected, actual)
		},
	)
}

func TestNamesWriterExportsReading(t *testing.T) {
	testNames(t,
		func(test *require.Assertions, db *db.NamesWriter) {
			err := db.WriteNames([]*jmnedict.Name{
				{
					Sequence: 1001,
					Reading:  []string{"reading 1"},
				},
				{
					Sequence: 1002,
					Reading:  []string{"reading 2a", "reading 2b"},
				},
				{
					Sequence: 9999,
				},
			})
			test.NoError(err)
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := []Data{
				{
					"sequence": int64(1001),
					"reading":  "reading 1",
				},
				{
					"sequence": int64(1002),
					"reading":  "reading 2a\treading 2b",
				},
				{
					"sequence": int64(9999),
					"reading":  "",
				},
			}
			actual := testQuery(test, db,
				"SELECT sequence, reading FROM names")
			test.EqualValues(expected, actual)
		},
	)
}

func TestNamesWriterExportsSense(t *testing.T) {
	testNames(t,
		func(test *require.Assertions, db *db.NamesWriter) {
			err := db.WriteNames([]*jmnedict.Name{
				{
					Sequence: 1001,
					Sense: []jmnedict.NameSense{
						{},
					},
				},
				{
					Sequence: 1002,
					Sense: []jmnedict.NameSense{
						{
							Info:        []string{"info A1"},
							XRef:        []string{"xref A1"},
							Translation: []string{"translation A1"},
						},
						{
							Info:        []string{"info B1", "info B2"},
							XRef:        []string{"xref B1", "xref B2"},
							Translation: []string{"translation B1", "translation B2"},
						},
					},
				},
				{
					Sequence: 9999,
				},
			})
			test.NoError(err)
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := []Data{
				{
					"sequence":    int64(1001),
					"position":    int64(0),
					"info":        "",
					"xref":        "",
					"translation": "",
				},
				{
					"sequence":    int64(1002),
					"position":    int64(0),
					"info":        "info A1",
					"xref":        "xref A1",
					"translation": "translation A1",
				},
				{
					"sequence":    int64(1002),
					"position":    int64(1),
					"info":        "info B1\tinfo B2",
					"xref":        "xref B1\txref B2",
					"translation": "translation B1\ttranslation B2",
				},
			}
			actual := testQuery(test, db, `
				SELECT
					sequence, position, info, xref, translation
				FROM names_sense`)
			test.EqualValues(expected, actual)
		},
	)
}

func testNames(t *testing.T, prepare func(test *require.Assertions, db *db.NamesWriter), eval func(test *require.Assertions, db *sql.DB)) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		func() {
			db, dbErr := db.NewNamesWriter(dbFile)
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
