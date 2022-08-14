package db_test

import (
	"database/sql"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ax-lab/tango/import/db"
	"github.com/ax-lab/tango/import/jmdict"
	"github.com/stretchr/testify/require"
)

func TestEntriesWriterExportsEntries(t *testing.T) {
	testData(t,
		func(db *db.EntriesWriter) {
			db.WriteEntries([]*jmdict.Entry{
				{Sequence: 1001},
				{Sequence: 1002},
				{Sequence: 1003},
			})
		},
		func(test *require.Assertions, db TestData) {
			expected := []Data{
				{"sequence": int64(1001)},
				{"sequence": int64(1002)},
				{"sequence": int64(1003)},
			}
			actual := db.Query("SELECT sequence FROM entries")
			test.EqualValues(expected, actual)
		},
	)
}

type Data map[string]interface{}

type TestData struct {
	test *require.Assertions
	db   *sql.DB
}

func (data TestData) Query(sql string) (out []Data) {
	rows, err := data.db.Query(sql)
	data.test.NoError(err, "error executing query")
	defer rows.Close()

	cols, err := rows.Columns()
	data.test.NoError(err, "error retrieving column names")

	row, args := make([]interface{}, len(cols)), make([]any, len(cols))
	for i := range row {
		args[i] = &row[i]
	}

	for rows.Next() {
		scanErr := rows.Scan(args...)
		data.test.NoError(scanErr, "error retrieving row values")

		data := make(Data)
		for i, name := range cols {
			data[name] = row[i]
			row[i] = nil
		}
		out = append(out, data)
	}

	return out
}

func testData(t *testing.T, prepare func(db *db.EntriesWriter), eval func(test *require.Assertions, db TestData)) {
	file, err := ioutil.TempFile("", "tango-test-db-*.db")
	if err != nil {
		panic(err)
	}

	file.Close()
	defer os.Remove(file.Name())

	func() {
		db, dbErr := db.NewEntriesWriter(file.Name())
		if dbErr != nil {
			panic(dbErr)
		}

		defer db.Close()
		prepare(db)
	}()

	func() {
		db, err := sql.Open("sqlite3", file.Name())
		if err != nil {
			panic(err)
		}
		defer db.Close()

		test := require.New(t)
		eval(test, TestData{test, db})
	}()
}
