package data

import (
	"database/sql"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

type Data map[string]interface{}

type Database struct {
	inner *sql.DB
	err   error
}

func OpenDB(dir, file string) *Database {
	filePath := path.Join(dir, file)
	db, err := sql.Open("sqlite3", filePath)

	out := &Database{inner: db, err: nil}
	out.FlagError(err)
	return out
}

func (db *Database) Close() {
	if db.inner != nil {
		db.inner.Close()
		db.inner = nil
	}
}

func (db *Database) FlagError(err error) bool {
	if err != nil {
		db.err = err
	}
	return db.err != nil
}

func (db *Database) Error() error {
	return db.err
}

func (db *Database) HasError() bool {
	return db.err != nil
}

type Scanner interface {
	Query() string
	Read(rows *sql.Rows) error
}

type ScannerResult struct {
	database *Database
	scanner  Scanner
	rows     *sql.Rows
	hasNext  bool
}

func (res *ScannerResult) Close() {
	if res.rows != nil {
		res.rows.Close()
		res.rows = nil
	}
}

func (res *ScannerResult) Next() bool {
	if res.hasNext {
		res.hasNext = false
		return true
	}

	if res.rows != nil && res.rows.Next() {
		if err := res.scanner.Read(res.rows); res.database.FlagError(err) {
			return false
		}
		return true
	}

	return false
}

func (res *ScannerResult) Unget() {
	res.hasNext = true
}

func (db *Database) ScanTable(scanner Scanner) *ScannerResult {
	out := &ScannerResult{
		database: db,
		scanner:  scanner,
	}

	if db.HasError() {
		return out
	}

	rows, err := db.inner.Query(scanner.Query())
	if db.FlagError(err) {
		return out
	}

	out.rows = rows
	return out
}
