package db_test

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ax-lab/tango/import/db"
	"github.com/stretchr/testify/require"
)

func TestOpensDatabase(t *testing.T) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		sql, err := db.Open(dbFile, `
			CREATE TABLE test(name TEXT)
		`)
		test.NoError(err)
		test.NotNil(sql)

		_, err = sql.Exec(`INSERT INTO test(name) VALUES ('abc')`)
		sql.Close()
		test.NoError(err)

		sql, err = db.Open(dbFile, `
			CREATE TABLE IF NOT EXISTS test(name TEXT)
		`)
		test.NoError(err)
		test.NotNil(sql)
		defer sql.Close()

		rows, err := sql.Query(`SELECT * FROM test`)
		test.NoError(err)
		test.True(rows.Next())

		var name string
		test.NoError(rows.Scan(&name))
		test.Equal("abc", name)
	})
}

func TestOpenFailsOnInvalidFile(t *testing.T) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		// use the file as a directory to force an error
		sql, err := db.Open(filepath.Join(dbFile, "some-db-file.db"), `
			CREATE TABLE test(name TEXT)
		`)
		test.Error(err)
		test.Nil(sql)
	})
}

func TestOpenFailsOnInvalidSchema(t *testing.T) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		// use the file as a directory to force an error
		sql, err := db.Open(dbFile, `
			CREATE TABLEx
		`)
		test.ErrorContains(err, "syntax error")
		test.Nil(sql)
	})
}

func TestValidTransaction(t *testing.T) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		func() {
			sql, err := db.Open(dbFile, "CREATE TABLE test(name TEXT)")
			test.NoError(err)
			defer sql.Close()

			tx := db.BeginTransaction(sql)
			insert := tx.Prepare("INSERT INTO test(name) VALUES (?)")
			insert.Exec("abc")
			insert.Exec("123")
			err = tx.Finish()
			test.NoError(err)
		}()

		sql, err := db.Open(dbFile, "")
		test.NoError(err)
		defer sql.Close()

		expected := []Data{
			{"name": "abc"},
			{"name": "123"},
		}
		rows := testQuery(test, sql, "SELECT * FROM test")
		test.EqualValues(expected, rows)
	})
}

func TestFailedTransaction(t *testing.T) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		func() {
			sql, err := db.Open(dbFile, "CREATE TABLE test(name TEXT NOT NULL)")
			test.NoError(err)
			defer sql.Close()

			tx := db.BeginTransaction(sql)
			insert := tx.Prepare("INSERT INTO test(name) VALUES (?)")
			insert.Exec("abc")
			insert.Exec(nil) // this will fail
			insert.Exec("123")
			err = tx.Finish()
			test.ErrorContains(err, "constraint")
		}()

		sql, err := db.Open(dbFile, "")
		test.NoError(err)
		defer sql.Close()

		expected := []Data{
			{"c": int64(0)}, // test that the transaction was rollback
		}
		rows := testQuery(test, sql, "SELECT COUNT(*) AS c FROM test")
		test.EqualValues(expected, rows)
	})
}

func testTempDB(t *testing.T, callback func(test *require.Assertions, dbFile string)) {
	test := require.New(t)

	file, err := ioutil.TempFile("", "tango-test-db-*.db")
	test.NoError(err)

	file.Close()
	defer os.Remove(file.Name())

	callback(test, file.Name())
}

type Data map[string]interface{}

func testQuery(test *require.Assertions, db *sql.DB, query string) (out []Data) {
	rows, err := db.Query(query)
	test.NoError(err, "error executing query")
	defer rows.Close()

	colNames, err := rows.Columns()
	test.NoError(err, "error retrieving column names")

	cols, args := make([]interface{}, len(colNames)), make([]any, len(colNames))
	for i := range cols {
		args[i] = &cols[i]
	}

	for rows.Next() {
		scanErr := rows.Scan(args...)
		test.NoError(scanErr, "error retrieving row values")

		row := make(Data)
		for i, name := range colNames {
			row[name] = cols[i]
			cols[i] = nil
		}
		out = append(out, row)
	}

	return out
}
