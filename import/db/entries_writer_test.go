package db_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/ax-lab/tango/import/db"
	"github.com/ax-lab/tango/import/jmdict"
	"github.com/stretchr/testify/require"
)

func TestEntriesWriterFailsOnOpenError(t *testing.T) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		w, err := db.NewEntriesWriter(filepath.Join(dbFile, "force-error.db"))
		test.Error(err)
		test.Nil(w)
	})
}

func TestEntriesWriterFailsOnInvalidInsert(t *testing.T) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		w, err := db.NewEntriesWriter(dbFile)
		test.NoError(err)
		err = w.WriteEntries([]*jmdict.Entry{
			{Sequence: 100},
			{Sequence: 100},
		})
		test.ErrorContains(err, "constraint")
	})
}

func TestEntriesWriterExportsEntries(t *testing.T) {
	testData(t,
		func(db *db.EntriesWriter) {
			db.WriteEntries([]*jmdict.Entry{
				{Sequence: 1001},
				{Sequence: 1002},
				{Sequence: 1003},
			})
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := []Data{
				{"sequence": int64(1001)},
				{"sequence": int64(1002)},
				{"sequence": int64(1003)},
			}
			actual := testQuery(test, db, "SELECT sequence FROM entries")
			test.EqualValues(expected, actual)
		},
	)
}

func TestEntriesWriterExportsKanji(t *testing.T) {
	testData(t,
		func(db *db.EntriesWriter) {
			db.WriteEntries([]*jmdict.Entry{
				{
					Sequence: 1001,
					Kanji: []jmdict.EntryKanji{
						{Text: "kanji 1"},
					},
				},
				{
					Sequence: 1002,
					Kanji: []jmdict.EntryKanji{
						{
							Text:     "kanji 2a",
							Info:     "info 2a",
							Priority: "priority 2a",
						},
						{
							Text:     "kanji 2b",
							Info:     "info 2b",
							Priority: "priority 2b",
						},
					},
				},
				{
					Sequence: 1003,
				},
			})
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := []Data{
				{
					"sequence": int64(1001),
					"position": int64(0),
					"text":     "kanji 1",
					"info":     "",
					"priority": "",
				},
				{
					"sequence": int64(1002),
					"position": int64(0),
					"text":     "kanji 2a",
					"info":     "info 2a",
					"priority": "priority 2a",
				},
				{
					"sequence": int64(1002),
					"position": int64(1),
					"text":     "kanji 2b",
					"info":     "info 2b",
					"priority": "priority 2b",
				},
			}
			actual := testQuery(test, db, "SELECT sequence, position, text, info, priority FROM entries_kanji")
			test.EqualValues(expected, actual)
		},
	)
}

func testData(t *testing.T, prepare func(db *db.EntriesWriter), eval func(test *require.Assertions, db *sql.DB)) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		func() {
			db, dbErr := db.NewEntriesWriter(dbFile)
			if dbErr != nil {
				panic(dbErr)
			}

			defer db.Close()
			prepare(db)
		}()

		func() {
			db, err := sql.Open("sqlite3", dbFile)
			if err != nil {
				panic(err)
			}
			defer db.Close()

			test := require.New(t)
			eval(test, db)
		}()
	})
}
