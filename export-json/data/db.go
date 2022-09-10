package data

import (
	"database/sql"
	"path"
	"strings"

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
	return &Database{db, err}
}

func (db *Database) Done() error {
	if db.inner != nil {
		db.inner.Close()
		db.inner = nil
	}
	return db.err
}

func (db *Database) LoadTable(sql string, args ...any) (out []Data) {
	if db.err != nil {
		return
	}

	rows, err := db.inner.Query(sql, args...)
	if err != nil {
		db.err = err
		return
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		db.err = err
		return
	}

	rowBuffer := make([]any, len(cols))
	rowBufferPtr := make([]any, len(cols))
	for i := range rowBuffer {
		rowBufferPtr[i] = &rowBuffer[i]
	}

	for rows.Next() {
		if db.err = rows.Scan(rowBufferPtr...); db.err != nil {
			return
		}
		row := make(Data)
		for i, name := range cols {
			row[name] = rowBuffer[i]
		}
		out = append(out, row)
	}

	return out
}

func splitTabs(value interface{}) []string {
	input := value.(string)
	if input == "" {
		return nil
	}
	return strings.Split(input, "\t")
}
