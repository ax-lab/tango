package data

import "database/sql"

type Tag struct {
	Name string
	Desc string
}

func (tb *Tag) Query() string {
	return "SELECT name, desc FROM tag"
}

func (tb *Tag) Read(row *sql.Rows) error {
	return row.Scan(&tb.Name, &tb.Desc)
}
