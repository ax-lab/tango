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

func TestEntriesWriterCanRewriteTheDatabase(t *testing.T) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		w, err := db.NewEntriesWriter(dbFile)
		test.NoError(err)
		w.Close()

		w, err = db.NewEntriesWriter(dbFile)
		test.NoError(err)
		w.Close()
	})
}

func TestEntriesWriterExportsEntries(t *testing.T) {
	testEntries(t,
		func(test *require.Assertions, db *db.EntriesWriter) {
			err := db.WriteEntries([]*jmdict.Entry{
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
			actual := testQuery(test, db, "SELECT sequence FROM entry")
			test.EqualValues(expected, actual)
		},
	)
}

func TestEntriesWriterExportsTags(t *testing.T) {
	testEntries(t,
		func(test *require.Assertions, db *db.EntriesWriter) {
			err := db.WriteTags(map[string]string{
				"a": "tag a",
				"b": "tag b",
				"c": "tag c",
			})
			test.NoError(err)
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := []Data{
				{"name": "a", "desc": "tag a"},
				{"name": "b", "desc": "tag b"},
				{"name": "c", "desc": "tag c"},
			}
			actual := testQuery(test, db, "SELECT name, desc FROM tag")
			test.EqualValues(expected, actual)
		},
	)
}

func TestEntriesWriterExportsKanji(t *testing.T) {
	testEntries(t,
		func(test *require.Assertions, db *db.EntriesWriter) {
			err := db.WriteEntries([]*jmdict.Entry{
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
							Info:     []string{"info 2a"},
							Priority: []string{"priority 2a"},
						},
						{
							Text:     "kanji 2b",
							Info:     []string{"info 2b", "x"},
							Priority: []string{"priority 2b", "y"},
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
					"info":     "info 2b\tx",
					"priority": "priority 2b\ty",
				},
			}
			actual := testQuery(test, db, "SELECT sequence, position, text, info, priority FROM entry_kanji")
			test.EqualValues(expected, actual)
		},
	)
}

func TestEntriesWriterExportsReading(t *testing.T) {
	testEntries(t,
		func(test *require.Assertions, db *db.EntriesWriter) {
			err := db.WriteEntries([]*jmdict.Entry{
				{
					Sequence: 1001,
					Reading: []jmdict.EntryReading{
						{Text: "reading 1"},
					},
				},
				{
					Sequence: 1002,
					Reading: []jmdict.EntryReading{
						{
							Text:        "reading 2a",
							Info:        []string{"info 2a"},
							Priority:    []string{"priority 2a"},
							Restriction: []string{"restriction 2a"},
						},
						{
							NoKanji:     true,
							Text:        "reading 2b",
							Info:        []string{"info 2b", "x"},
							Priority:    []string{"priority 2b", "y"},
							Restriction: []string{"restriction 2b", "z"},
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
					"text":        "reading 1",
					"info":        "",
					"priority":    "",
					"restriction": "",
					"no_kanji":    int64(0),
				},
				{
					"sequence":    int64(1002),
					"position":    int64(0),
					"text":        "reading 2a",
					"info":        "info 2a",
					"priority":    "priority 2a",
					"restriction": "restriction 2a",
					"no_kanji":    int64(0),
				},
				{
					"sequence":    int64(1002),
					"position":    int64(1),
					"text":        "reading 2b",
					"info":        "info 2b\tx",
					"priority":    "priority 2b\ty",
					"restriction": "restriction 2b\tz",
					"no_kanji":    int64(1),
				},
			}
			actual := testQuery(test, db,
				"SELECT sequence, position, text, info, priority, restriction, no_kanji FROM entry_reading")
			test.EqualValues(expected, actual)
		},
	)
}

func TestEntriesWriterExportsSense(t *testing.T) {
	testEntries(t,
		func(test *require.Assertions, db *db.EntriesWriter) {
			err := db.WriteEntries([]*jmdict.Entry{
				{
					Sequence: 1001,
					Sense: []jmdict.EntrySense{
						{},
					},
				},
				{
					Sequence: 1002,
					Sense: []jmdict.EntrySense{
						{
							Info:         []string{"info A"},
							PartOfSpeech: []string{"pos A"},
							StagKanji:    []string{"stagk A"},
							StagReading:  []string{"stagr A"},
							Field:        []string{"field A"},
							Misc:         []string{"misc A"},
							Dialect:      []string{"dialect A"},
							Antonym:      []string{"antonym A"},
							XRef:         []string{"xref A"},
						},
						{
							Info:         []string{"info A", "info B"},
							PartOfSpeech: []string{"pos A", "pos B"},
							StagKanji:    []string{"stagk A", "stagk B"},
							StagReading:  []string{"stagr A", "stagr B"},
							Field:        []string{"field A", "field B"},
							Misc:         []string{"misc A", "misc B"},
							Dialect:      []string{"dialect A", "dialect B"},
							Antonym:      []string{"antonym A", "antonym B"},
							XRef:         []string{"xref A", "xref B"},
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
					"sequence": int64(1001),
					"position": int64(0),
					"info":     "",
					"pos":      "",
					"stagk":    "",
					"stagr":    "",
					"field":    "",
					"misc":     "",
					"dialect":  "",
					"antonym":  "",
					"xref":     "",
				},
				{
					"sequence": int64(1002),
					"position": int64(0),
					"info":     "info A",
					"pos":      "pos A",
					"stagk":    "stagk A",
					"stagr":    "stagr A",
					"field":    "field A",
					"misc":     "misc A",
					"dialect":  "dialect A",
					"antonym":  "antonym A",
					"xref":     "xref A",
				},
				{
					"sequence": int64(1002),
					"position": int64(1),
					"info":     "info A\tinfo B",
					"pos":      "pos A\tpos B",
					"stagk":    "stagk A\tstagk B",
					"stagr":    "stagr A\tstagr B",
					"field":    "field A\tfield B",
					"misc":     "misc A\tmisc B",
					"dialect":  "dialect A\tdialect B",
					"antonym":  "antonym A\tantonym B",
					"xref":     "xref A\txref B",
				},
			}
			actual := testQuery(test, db, `
				SELECT
					sequence, position, info, pos, stagk, stagr, field, misc,
					dialect, antonym, xref
				FROM entry_sense`)
			test.EqualValues(expected, actual)
		},
	)
}

func TestEntriesWriterExportsSenseGlossary(t *testing.T) {
	testEntries(t,
		func(test *require.Assertions, db *db.EntriesWriter) {
			err := db.WriteEntries([]*jmdict.Entry{
				{
					Sequence: 1001,
					Sense: []jmdict.EntrySense{
						{
							Glossary: []jmdict.EntrySenseGlossary{
								{Text: "entry A"},
							},
						},
					},
				},
				{
					Sequence: 1002,
					Sense: []jmdict.EntrySense{
						{
							Glossary: []jmdict.EntrySenseGlossary{
								{Text: "entry B1 - A"},
								{Text: "entry B1 - B", Lang: "por", Type: "lit"},
							},
						},
						{
							Glossary: []jmdict.EntrySenseGlossary{
								{Text: "entry B2 - A"},
								{Text: "entry B2 - B", Lang: "ger"},
								{Text: "entry B2 - C", Lang: "fra", Type: "fig"},
							},
						},
					},
				},
				{
					Sequence: 9999,
					Sense: []jmdict.EntrySense{
						{},
					},
				},
			})
			test.NoError(err)
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := []Data{
				{
					"sequence": int64(1001),
					"sense":    int64(0),
					"position": int64(0),
					"text":     "entry A",
					"lang":     "",
					"type":     "",
				},
				{
					"sequence": int64(1002),
					"sense":    int64(0),
					"position": int64(0),
					"text":     "entry B1 - A",
					"lang":     "",
					"type":     "",
				},
				{
					"sequence": int64(1002),
					"sense":    int64(0),
					"position": int64(1),
					"text":     "entry B1 - B",
					"lang":     "por",
					"type":     "lit",
				},
				{
					"sequence": int64(1002),
					"sense":    int64(1),
					"position": int64(0),
					"text":     "entry B2 - A",
					"lang":     "",
					"type":     "",
				},
				{
					"sequence": int64(1002),
					"sense":    int64(1),
					"position": int64(1),
					"text":     "entry B2 - B",
					"lang":     "ger",
					"type":     "",
				},
				{
					"sequence": int64(1002),
					"sense":    int64(1),
					"position": int64(2),
					"text":     "entry B2 - C",
					"lang":     "fra",
					"type":     "fig",
				},
			}
			actual := testQuery(test, db, `
				SELECT
					sequence, sense, position, text, lang, type
				FROM entry_sense_glossary`)
			test.EqualValues(expected, actual)
		},
	)
}

func TestEntriesWriterExportsSenseSource(t *testing.T) {
	testEntries(t,
		func(test *require.Assertions, db *db.EntriesWriter) {
			err := db.WriteEntries([]*jmdict.Entry{
				{
					Sequence: 1001,
					Sense: []jmdict.EntrySense{
						{
							Source: []jmdict.EntrySenseSource{
								{Text: "source A"},
							},
						},
					},
				},
				{
					Sequence: 1002,
					Sense: []jmdict.EntrySense{
						{
							Source: []jmdict.EntrySenseSource{
								{Text: "source B1 - A"},
								{Text: "source B1 - B", Lang: "por", Type: "part"},
							},
						},
						{
							Source: []jmdict.EntrySenseSource{
								{Text: "source B2 - A"},
								{Text: "source B2 - B", Lang: "ger", Wasei: "n"},
								{Text: "source B2 - C", Lang: "fra", Type: "full", Wasei: "y"},
							},
						},
					},
				},
				{
					Sequence: 9999,
					Sense: []jmdict.EntrySense{
						{},
					},
				},
			})
			test.NoError(err)
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := []Data{
				{
					"sequence": int64(1001),
					"sense":    int64(0),
					"position": int64(0),
					"text":     "source A",
					"lang":     "",
					"type":     "",
					"wasei":    "",
				},
				{
					"sequence": int64(1002),
					"sense":    int64(0),
					"position": int64(0),
					"text":     "source B1 - A",
					"lang":     "",
					"type":     "",
					"wasei":    "",
				},
				{
					"sequence": int64(1002),
					"sense":    int64(0),
					"position": int64(1),
					"text":     "source B1 - B",
					"lang":     "por",
					"type":     "part",
					"wasei":    "",
				},
				{
					"sequence": int64(1002),
					"sense":    int64(1),
					"position": int64(0),
					"text":     "source B2 - A",
					"lang":     "",
					"type":     "",
					"wasei":    "",
				},
				{
					"sequence": int64(1002),
					"sense":    int64(1),
					"position": int64(1),
					"text":     "source B2 - B",
					"lang":     "ger",
					"type":     "",
					"wasei":    "n",
				},
				{
					"sequence": int64(1002),
					"sense":    int64(1),
					"position": int64(2),
					"text":     "source B2 - C",
					"lang":     "fra",
					"type":     "full",
					"wasei":    "y",
				},
			}
			actual := testQuery(test, db, `
				SELECT
					sequence, sense, position, text, lang, type, wasei
				FROM entry_sense_source`)
			test.EqualValues(expected, actual)
		},
	)
}

func testEntries(t *testing.T, prepare func(test *require.Assertions, db *db.EntriesWriter), eval func(test *require.Assertions, db *sql.DB)) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		func() {
			db, dbErr := db.NewEntriesWriter(dbFile)
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
